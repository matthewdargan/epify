// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package epify provides facilities for categorizing [shows] and [movies]
// using the Jellyfin naming scheme.
//
// [shows]: https://jellyfin.org/docs/general/server/media/shows/
// [movies]: https://jellyfin.org/docs/general/server/media/movies/
package epify

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Media represents metadata for a show or movie.
type Media struct {
	Name string // The media name.
	Year string // The year the media premiered.
	ID   string // The media ID.
	Dir  string // The directory to create the media in.
}

// MkShow creates a show directory like "Series Name (2018) [tvdbid-65567]".
func MkShow(s *Media) error {
	if len(s.Name) == 0 {
		return errors.New("empty show name")
	}
	year, err := strconv.Atoi(s.Year)
	if err != nil {
		return fmt.Errorf("invalid year: %w", err)
	}
	tvdbid, err := strconv.Atoi(s.ID)
	if err != nil {
		return fmt.Errorf("invalid TVDBID: %w", err)
	}
	path := fmt.Sprintf("%s (%d) [tvdbid-%d]", s.Name, year, tvdbid)
	if err := os.MkdirAll(filepath.Join(s.Dir, path), 0o755); err != nil {
		return err
	}
	return nil
}

// A Movie represents a movie.
type Movie struct {
	Media
	File string // The movie file to add.
}

// AddMovie adds a movie to a directory. Movies are labeled like
// "Film (2018) [tmdbid-65567]".
func AddMovie(m *Movie) error {
	if len(m.Name) == 0 {
		return errors.New("empty movie name")
	}
	year, err := strconv.Atoi(m.Year)
	if err != nil {
		return fmt.Errorf("invalid year: %w", err)
	}
	tmdbid, err := strconv.Atoi(m.ID)
	if err != nil {
		return fmt.Errorf("invalid TMDBID: %w", err)
	}
	info, err := os.Stat(m.Dir)
	if err != nil {
		return fmt.Errorf("invalid directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", m.Dir)
	}
	info, err = os.Stat(m.File)
	if err != nil {
		return fmt.Errorf("invalid movie: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("%q is a directory", m.File)
	}
	path := fmt.Sprintf("%s (%d) [tmdbid-%d]%s", m.Name, year, tmdbid, filepath.Ext(m.File))
	if err := os.Rename(m.File, filepath.Join(m.Dir, path)); err != nil {
		return err
	}
	return nil
}

// A Season represents a season of a TV show.
type Season struct {
	N          string   // The season number.
	ShowDir    string   // The show directory.
	Episodes   []string // The episodes to populate the season.
	MatchIndex int      // The index of the episode number in filenames.
}

var errNoEpisodes = errors.New("no episodes found")

// YearPrefix is the prefix for a show directory.
const YearPrefix = " ("

// MkSeason populates a season directory with episodes. Episodes are labeled
// like "Series Name S01E01.mkv".
func MkSeason(s *Season) error {
	n, err := strconv.Atoi(s.N)
	if err != nil {
		return fmt.Errorf("invalid season: %w", err)
	}
	info, err := os.Stat(s.ShowDir)
	if err != nil {
		return fmt.Errorf("invalid directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", s.ShowDir)
	}
	show, _, ok := strings.Cut(filepath.Base(s.ShowDir), YearPrefix)
	if !ok {
		return fmt.Errorf("invalid directory %q", s.ShowDir)
	}
	if len(s.Episodes) == 0 {
		return errNoEpisodes
	}
	for _, e := range s.Episodes {
		info, err = os.Stat(e)
		if err != nil {
			return fmt.Errorf("invalid episode: %w", err)
		}
		if info.IsDir() {
			return fmt.Errorf("%q is a directory", e)
		}
	}
	if err = sortEpisodes(s.Episodes, s.MatchIndex); err != nil {
		return err
	}
	path := fmt.Sprintf("Season %02d", n)
	seasonDir := filepath.Join(s.ShowDir, path)
	if err = os.Mkdir(seasonDir, 0o755); err != nil {
		return err
	}
	for i, e := range s.Episodes {
		ep := fmt.Sprintf("%s S%02dE%02d%s", show, n, i+1, filepath.Ext(e))
		if err := os.Rename(e, filepath.Join(seasonDir, ep)); err != nil {
			return err
		}
	}
	return nil
}

// A SeasonAddition represents episodes to add to a season.
type SeasonAddition struct {
	SeasonDir  string   // The season directory.
	Episodes   []string // The episodes to add.
	MatchIndex int      // The index of the episode number in filenames.
}

var episodeRe = regexp.MustCompile(`E(\d+)\.`)

// AddEpisodes adds episodes to a season directory, continuing at the previous
// episode increment.
func AddEpisodes(s *SeasonAddition) error {
	info, err := os.Stat(s.SeasonDir)
	if err != nil {
		return fmt.Errorf("invalid season directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", s.SeasonDir)
	}
	base := filepath.Base(s.SeasonDir)
	season := strings.TrimPrefix(base, "Season ")
	if base == season {
		return fmt.Errorf("invalid season directory %q", s.SeasonDir)
	}
	n, err := strconv.Atoi(season)
	if err != nil {
		return fmt.Errorf("invalid season: %w", err)
	}
	showDir := filepath.Dir(s.SeasonDir)
	show, _, ok := strings.Cut(filepath.Base(showDir), YearPrefix)
	if !ok {
		return fmt.Errorf("invalid show directory %q", showDir)
	}
	if len(s.Episodes) == 0 {
		return errNoEpisodes
	}
	for _, e := range s.Episodes {
		info, err = os.Stat(e)
		if err != nil {
			return fmt.Errorf("invalid episode: %w", err)
		}
		if info.IsDir() {
			return fmt.Errorf("%q is a directory", e)
		}
	}
	if err = sortEpisodes(s.Episodes, s.MatchIndex); err != nil {
		return err
	}
	ents, err := os.ReadDir(s.SeasonDir)
	if err != nil {
		return err
	}
	var epn int
	if len(ents) > 0 {
		prevEp := ents[len(ents)-1].Name()
		m := episodeRe.FindStringSubmatch(prevEp)
		if len(m) != 2 {
			return fmt.Errorf("invalid episode %q", prevEp)
		}
		epn, _ = strconv.Atoi(m[1])
	}
	for _, e := range s.Episodes {
		epn++
		ep := fmt.Sprintf("%s S%02dE%02d%s", show, n, epn, filepath.Ext(e))
		if err := os.Rename(e, filepath.Join(s.SeasonDir, ep)); err != nil {
			return err
		}
	}
	return nil
}

var re = regexp.MustCompile(`\d+`)

func sortEpisodes(eps []string, i int) error {
	for _, e := range eps {
		base := filepath.Base(e)
		m := re.FindAllString(base, -1)
		if len(m) == 0 {
			return fmt.Errorf("episode %q must contain number", e)
		}
		if i < 0 || i >= len(m) {
			return fmt.Errorf("invalid match index %d", i)
		}
	}
	slices.SortFunc(eps, func(a, b string) int {
		e1, _ := strconv.Atoi(re.FindAllString(filepath.Base(a), -1)[i])
		e2, _ := strconv.Atoi(re.FindAllString(filepath.Base(b), -1)[i])
		return cmp.Compare(e1, e2)
	})
	return nil
}
