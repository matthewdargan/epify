// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// A Show represents a TV show.
type Show struct {
	Name   string // The name of the show.
	Year   string // The year the show premiered.
	TVDBID string // The TVDB ID of the show.
	Dir    string // The directory to create the show in.
}

// MkShow creates a show directory like "Series Name (2018) [tvdbid-65567]".
func MkShow(s *Show) error {
	if len(s.Name) == 0 {
		return errors.New("empty show name")
	}
	year, err := strconv.Atoi(s.Year)
	if err != nil {
		return fmt.Errorf("invalid year: %w", err)
	}
	tvdbid, err := strconv.Atoi(s.TVDBID)
	if err != nil {
		return fmt.Errorf("invalid TVDBID: %w", err)
	}
	path := fmt.Sprintf("%s (%d) [tvdbid-%d]", s.Name, year, tvdbid)
	if err := os.MkdirAll(filepath.Join(s.Dir, path), 0o755); err != nil {
		return err
	}
	return nil
}

// A Season represents a season of a TV show.
type Season struct {
	N        string   // The season number.
	ShowDir  string   // The show directory.
	Episodes []string // The episodes to populate the season.
}

var errNoEpisodes = errors.New("no episodes found")

// MkSeason populates a season directory with episodes. Episodes are labeled
// like "Series Name S01E01.mkv".
func MkSeason(s *Season) error {
	n, err := strconv.Atoi(s.N)
	if err != nil {
		return fmt.Errorf("invalid season: %w", err)
	}
	info, err := os.Stat(s.ShowDir)
	if err != nil {
		return fmt.Errorf("invalid show directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", s.ShowDir)
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
	if err = sortEpisodes(s.Episodes); err != nil {
		return err
	}
	path := fmt.Sprintf("Season %02d", n)
	seasonDir := filepath.Join(s.ShowDir, path)
	if err = os.Mkdir(seasonDir, 0o755); err != nil {
		return err
	}
	for i, e := range s.Episodes {
		ep := fmt.Sprintf("S%02dE%02d%s", n, i+1, filepath.Ext(e))
		if err := os.Rename(e, filepath.Join(seasonDir, ep)); err != nil {
			return err
		}
	}
	return nil
}

// A SeasonAddition represents episodes to add to a season.
type SeasonAddition struct {
	SeasonDir string
	Episodes  []string
}

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
	if err = sortEpisodes(s.Episodes); err != nil {
		return err
	}
	ents, err := os.ReadDir(s.SeasonDir)
	if err != nil {
		return err
	}
	var epn int
	if len(ents) > 0 {
		prevEp := ents[len(ents)-1].Name()
		i := strings.Index(prevEp, "E")
		j := strings.Index(prevEp, ".")
		if i == -1 || j == -1 || i >= j {
			return fmt.Errorf("invalid episode %q", prevEp)
		}
		epn, err = strconv.Atoi(prevEp[i+1 : j])
		if err != nil {
			return fmt.Errorf("invalid episode number: %w", err)
		}
	}
	for _, e := range s.Episodes {
		epn++
		ep := fmt.Sprintf("S%02dE%02d%s", n, epn, filepath.Ext(e))
		if err := os.Rename(e, filepath.Join(s.SeasonDir, ep)); err != nil {
			return err
		}
	}
	return nil
}

var re = regexp.MustCompile(`\d+`)

func sortEpisodes(eps []string) error {
	for _, e := range eps {
		base := filepath.Base(e)
		if _, err := strconv.Atoi(re.FindString(base)); err != nil {
			return fmt.Errorf("episode %q must contain number: %w", e, err)
		}
	}
	slices.SortFunc(eps, func(a, b string) int {
		e1, _ := strconv.Atoi(re.FindString(filepath.Base(a)))
		e2, _ := strconv.Atoi(re.FindString(filepath.Base(b)))
		return cmp.Compare(e1, e2)
	})
	return nil
}
