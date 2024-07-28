// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package media provides facilities for categorizing [shows] and [movies]
// using the Jellyfin naming scheme.
//
// [shows]: https://jellyfin.org/docs/general/server/media/shows/
// [movies]: https://jellyfin.org/docs/general/server/media/movies/
package media

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

	"golang.org/x/sync/errgroup"
)

// A Show represents a TV show.
type Show struct {
	Name, Year, ID, Dir string
}

// MkShow creates a show directory. The directory will be labeled like
// "Series Name (2018) [tvdbid-65567]".
func MkShow(s Show) error {
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
	Show
	File string
}

// AddMovie adds a movie to a directory. Movies are labeled like
// "Film (2018) [tmdbid-65567]".
func AddMovie(m Movie) error {
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

// A Season represents a TV show season.
type Season struct {
	N          string // season number
	ShowDir    string
	Episodes   []string
	MatchIndex int // index of the episode number in filenames
}

var errNoEpisodes = errors.New("no episodes found")

const YearSep = " (" // YearSep separates the show name from the year.

// MkSeason creates a season directory and moves episodes into it. Episodes are
// labeled like "Series Name S01E01.mkv".
func MkSeason(s Season) error {
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
	show, _, ok := strings.Cut(filepath.Base(s.ShowDir), YearSep)
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
	var g errgroup.Group
	for i, e := range s.Episodes {
		g.Go(func() error {
			ep := fmt.Sprintf("%s S%02dE%02d%s", show, n, i+1, filepath.Ext(e))
			return os.Rename(e, filepath.Join(seasonDir, ep))
		})
	}
	return g.Wait()
}

// An Addition represents episodes to add to a season.
type Addition struct {
	SeasonDir  string
	Episodes   []string
	MatchIndex int // index of the episode number in filenames
}

var episodeRe = regexp.MustCompile(`E(\d+)\.`)

// AddEpisodes adds episodes to a season directory. Episode numbers continue at
// the previous episode increment.
func AddEpisodes(a Addition) error {
	info, err := os.Stat(a.SeasonDir)
	if err != nil {
		return fmt.Errorf("invalid season directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", a.SeasonDir)
	}
	base := filepath.Base(a.SeasonDir)
	season := strings.TrimPrefix(base, "Season ")
	if base == season {
		return fmt.Errorf("invalid season directory %q", a.SeasonDir)
	}
	n, err := strconv.Atoi(season)
	if err != nil {
		return fmt.Errorf("invalid season: %w", err)
	}
	showDir := filepath.Dir(a.SeasonDir)
	show, _, ok := strings.Cut(filepath.Base(showDir), YearSep)
	if !ok {
		return fmt.Errorf("invalid show directory %q", showDir)
	}
	if len(a.Episodes) == 0 {
		return errNoEpisodes
	}
	for _, e := range a.Episodes {
		info, err = os.Stat(e)
		if err != nil {
			return fmt.Errorf("invalid episode: %w", err)
		}
		if info.IsDir() {
			return fmt.Errorf("%q is a directory", e)
		}
	}
	if err = sortEpisodes(a.Episodes, a.MatchIndex); err != nil {
		return err
	}
	ents, err := os.ReadDir(a.SeasonDir)
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
	var g errgroup.Group
	for i, e := range a.Episodes {
		g.Go(func() error {
			ep := fmt.Sprintf("%s S%02dE%02d%s", show, n, epn+i+1, filepath.Ext(e))
			return os.Rename(e, filepath.Join(a.SeasonDir, ep))
		})
	}
	return g.Wait()
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
