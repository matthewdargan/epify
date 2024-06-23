// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epify

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
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
		return errors.New("no episodes found")
	}
	fis := make([]fs.FileInfo, len(s.Episodes))
	for i, e := range s.Episodes {
		info, err = os.Stat(e)
		if err != nil {
			return fmt.Errorf("invalid episode: %w", err)
		}
		fis[i] = info
	}
	eps, err := episodes(fis)
	if err != nil {
		return err
	}
	if err = sortEpisodes(eps); err != nil {
		return err
	}
	path := fmt.Sprintf("Season %02d", n)
	seasonDir := filepath.Join(s.ShowDir, path)
	if err = os.Mkdir(seasonDir, 0o755); err != nil {
		return err
	}
	for i, e := range eps {
		ep := fmt.Sprintf("S%02dE%02d%s", n, i+1, filepath.Ext(e.Name()))
		if err := os.Rename(e.Name(), filepath.Join(seasonDir, ep)); err != nil {
			return err
		}
	}
	return nil
}

func episodes(fis []fs.FileInfo) ([]fs.FileInfo, error) {
	eps := make([]fs.FileInfo, 0, len(fis))
	for _, e := range fis {
		if !e.IsDir() {
			eps = append(eps, e)
			continue
		}
		ents, err := os.ReadDir(e.Name())
		if err != nil {
			return nil, fmt.Errorf("invalid episode directory: %w", err)
		}
		for _, ep := range ents {
			info, err := ep.Info()
			if err != nil {
				return nil, fmt.Errorf("invalid episode: %w", err)
			}
			eps = append(eps, info)
		}
	}
	return eps, nil
}

var re = regexp.MustCompile(`\d+`)

func sortEpisodes(eps []fs.FileInfo) error {
	for _, e := range eps {
		name := e.Name()
		if _, err := strconv.Atoi(re.FindString(name)); err != nil {
			return fmt.Errorf("episode %q must contain number: %w", name, err)
		}
	}
	slices.SortFunc(eps, func(a, b fs.FileInfo) int {
		e1, _ := strconv.Atoi(re.FindString(a.Name()))
		e2, _ := strconv.Atoi(re.FindString(b.Name()))
		return cmp.Compare(e1, e2)
	})
	return nil
}
