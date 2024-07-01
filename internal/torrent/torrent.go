// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package torrent provides facilities for organizing torrents.
package torrent

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/matthewdargan/epify/internal/media"
)

// A File represents a torrent file.
type File struct {
	Dir, Name, DstDir string
}

// Rename renames a torrent to an episode in a season directory.
func Rename(t *File) error {
	info, err := os.Stat(t.Dir)
	if err != nil {
		return fmt.Errorf("invalid directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", t.Dir)
	}
	ents, err := os.ReadDir(t.Dir)
	if err != nil {
		return err
	}
	if len(ents) == 0 {
		return errors.New("no shows found")
	}
	var showDir string
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		show, _, ok := strings.Cut(name, media.YearSep)
		if !ok {
			continue
		}
		if strings.Contains(t.Name, show) {
			showDir = filepath.Join(t.Dir, name)
			break
		}
	}
	if showDir == "" {
		return fmt.Errorf("no show directory for %q", t.Name)
	}
	ents, err = os.ReadDir(showDir)
	if err != nil {
		return err
	}
	if len(ents) == 0 {
		return errors.New("no seasons found")
	}
	var seasonDir string
	var largest int
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		season := strings.TrimPrefix(name, "Season ")
		if name == season {
			continue
		}
		n, err := strconv.Atoi(season)
		if err != nil {
			return fmt.Errorf("invalid season: %w", err)
		}
		if n > largest {
			largest = n
			seasonDir = filepath.Join(showDir, name)
		}
	}
	if seasonDir == "" {
		return fmt.Errorf("no season directory in %q", showDir)
	}
	a := media.Addition{
		SeasonDir: seasonDir,
		Episodes:  []string{filepath.Join(t.Dir, t.Name)},
	}
	if err := media.AddEpisodes(&a); err != nil {
		return err
	}
	return nil
}
