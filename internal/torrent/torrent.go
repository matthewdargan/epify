// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package torrent provides facilities for organizing torrents.
package torrent

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/matthewdargan/epify/internal/media"
)

// A File represents a torrent file.
type File struct {
	Dir    string
	Name   string
	DstDir string
}

// Rename renames a torrent to an episode in a season directory.
func Rename(t *File) error {
	info, err := os.Stat(t.Dir)
	if err != nil {
		log.Fatal(err)
	}
	if !info.IsDir() {
		log.Fatalf("%q is not a directory", t.Dir)
	}
	ents, err := os.ReadDir(t.Dir)
	if err != nil {
		log.Fatal(err)
	}
	if len(ents) == 0 {
		log.Fatalf("%q is empty", t.Dir)
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
		log.Fatalf("no corresponding show directory for %q", t.Name)
	}
	ents, err = os.ReadDir(showDir)
	if err != nil {
		log.Fatal(err)
	}
	if len(ents) == 0 {
		log.Fatalf("%q is empty", showDir)
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
			log.Fatalf("invalid season: %v", err)
		}
		if n > largest {
			largest = n
			seasonDir = filepath.Join(showDir, name)
		}
	}
	if seasonDir == "" {
		log.Fatalf("no season directory in %q", showDir)
	}
	a := media.Addition{
		SeasonDir: seasonDir,
		Episodes:  []string{filepath.Join(t.Dir, t.Name)},
	}
	if err := media.AddEpisodes(&a); err != nil {
		log.Fatal(err)
	}
	return nil
}
