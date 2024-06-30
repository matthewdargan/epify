// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Trdone organizes completed torrent downloads.
//
// Usage:
//
//	trdone dir
//
// Trdone should be used with the `script-torrent-done-enabled` and
// `script-torrent-done-filename` [Transmission settings].
//
// The `TR_TORRENT_DIR` and `TR_TORRENT_NAME` [environment variables] must be
// defined.
//
// Example:
//
// Move completed downloads into respective show directories in `/media/shows`:
//
//	$ trdone '/media/shows'
//
// [Transmission settings]: https://github.com/transmission/transmission/blob/main/docs/Editing-Configuration-Files.md#misc
// [environment variables]: https://github.com/transmission/transmission/blob/main/docs/Scripts.md#on-torrent-completion
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/matthewdargan/epify/internal/epify"
)

var (
	trDir  = os.Getenv("TR_TORRENT_DIR")
	trName = os.Getenv("TR_TORRENT_NAME")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: trdone dir\n")
	os.Exit(2)
}

func main() {
	log.SetPrefix("trdone: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	switch {
	case flag.NArg() != 1:
		usage()
	case trDir == "":
		log.Fatal("$TR_TORRENT_DIR is not defined")
	case trName == "":
		log.Fatal("$TR_TORRENT_NAME is not defined")
	}
	dir := flag.Arg(0)
	info, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err)
	}
	if !info.IsDir() {
		log.Fatalf("%q is not a directory", dir)
	}
	ents, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	if len(ents) == 0 {
		log.Fatalf("%q is empty", dir)
	}
	var showDir string
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		show, _, ok := strings.Cut(name, epify.YearPrefix)
		if !ok {
			continue
		}
		if strings.Contains(trName, show) {
			showDir = filepath.Join(dir, name)
			break
		}
	}
	if showDir == "" {
		log.Fatalf("no corresponding show directory for %q", trName)
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
	s := epify.SeasonAddition{
		SeasonDir: seasonDir,
		Episodes:  []string{filepath.Join(trDir, trName)},
	}
	if err := epify.AddEpisodes(&s); err != nil {
		log.Fatal(err)
	}
}
