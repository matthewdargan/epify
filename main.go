// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Epify categorizes shows using the [Jellyfin naming scheme].
//
// Usage:
//
//	epify show name year tvdbid dir
//	epify season seasonnum showdir episode...
//	epify add seasondir episode...
//
// `epify show` creates a show directory like
// "Series Name (2018) [tvdbid-65567]".
//
// `epify season` populates a season directory with episodes. Episodes are
// labeled like "Series Name S01E01.mkv".
//
// `epify add` adds episodes to a season directory, continuing at the previous
// episode increment.
//
// Examples:
//
// Create show directory `/media/shows/The Office (2005) [tvdbid-73244]`:
//
//	$ epify show 'The Office' 2005 73244 '/media/shows'
//
// Populate season directory
// `/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:
//
//	$ epify season 3 '/media/shows/The Office (2005) [tvdbid-73244]' '/downloads/the_office_s3_part_1' '/downloads/the_office_s3_ep13.mkv'
//
// Add episodes to `/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:
//
//	$ epify add '/media/shows/The Office (2005) [tvdbid-73244]/Season 03' '/downloads/the_office_s3_ep23.mkv' '/downloads/the_office_s3_p2'
//
// [Jellyfin naming scheme]: https://jellyfin.org/docs/general/server/media/shows/
package main

import (
	"cmp"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var re = regexp.MustCompile(`\d+`)

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "\tepify show name year tvdbid dir\n")
	fmt.Fprintf(os.Stderr, "\tepify season seasonnum showdir episode...\n")
	fmt.Fprintf(os.Stderr, "\tepify add seasondir episode...\n")
	os.Exit(2)
}

func main() {
	log.SetPrefix("epify: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}
	args := flag.Args()
	switch args[0] {
	case "show":
		if flag.NArg() != 5 {
			usage()
		}
		name := args[1]
		year, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatal(err)
		}
		tvdbid, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatal(err)
		}
		dir := args[4]
		parent := filepath.Dir(dir)
		if _, err := os.Stat(parent); err != nil {
			log.Fatal(err)
		}
		if err = mkShow(name, year, tvdbid, dir); err != nil {
			log.Fatal(err)
		}
	case "season":
		if flag.NArg() < 4 {
			usage()
		}
		season, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal(err)
		}
		showDir := args[2]
		info, err := os.Stat(showDir)
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() {
			log.Fatalf("%q is not a directory", showDir)
		}
		episodes := args[3:]
		efis := make([]fs.FileInfo, len(episodes))
		for i, e := range episodes {
			info, err = os.Stat(e)
			if err != nil {
				log.Fatal(err)
			}
			efis[i] = info
		}
		if err = mkSeason(season, showDir, efis); err != nil {
			log.Fatal(err)
		}
	case "add":
		if flag.NArg() < 3 {
			usage()
		}
		seasonDir := args[1]
		info, err := os.Stat(seasonDir)
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() {
			log.Fatalf("%q is not a directory", seasonDir)
		}
		episodes := args[2:]
		efis := make([]fs.FileInfo, len(episodes))
		for i, e := range episodes {
			info, err = os.Stat(e)
			if err != nil {
				log.Fatal(err)
			}
			efis[i] = info
		}
		if err = addEpisodes(seasonDir, efis); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}

func mkShow(name string, year int, tvdbid int, dir string) error {
	path := fmt.Sprintf("%s (%d) [tvdbid-%d]", name, year, tvdbid)
	if err := os.Mkdir(filepath.Join(dir, path), 0o750); err != nil {
		return err
	}
	return nil
}

func mkSeason(season int, showDir string, episodes []fs.FileInfo) error {
	if len(episodes) == 0 {
		return fmt.Errorf("no episodes found")
	}
	eps := make([]fs.FileInfo, 0, len(episodes))
	for _, e := range episodes {
		if e.IsDir() {
			ents, err := os.ReadDir(e.Name())
			if err != nil {
				return err
			}
			for _, ep := range ents {
				info, err := ep.Info()
				if err != nil {
					return err
				}
				eps = append(eps, info)
			}
		} else {
			eps = append(eps, e)
		}
	}
	if err := validateEps(eps); err != nil {
		return err
	}
	if err := sortEps(eps); err != nil {
		return err
	}
	path := fmt.Sprintf("Season %02d", season)
	seasonDir := filepath.Join(showDir, path)
	if err := os.Mkdir(seasonDir, 0o750); err != nil {
		log.Fatal(err)
	}
	for i, e := range eps {
		ep := fmt.Sprintf("S%02dE%02d%s", season, i+1, filepath.Ext(e.Name()))
		if err := os.Rename(e.Name(), filepath.Join(seasonDir, ep)); err != nil {
			return err
		}
	}
	return nil
}

func addEpisodes(seasonDir string, episodes []fs.FileInfo) error {
	if len(episodes) == 0 {
		return fmt.Errorf("no episodes found")
	}
	base := filepath.Base(seasonDir)
	season := strings.TrimLeft(base, "Season ")
	if base == season {
		return fmt.Errorf("invalid season directory %q", seasonDir)
	}
	eps := make([]fs.FileInfo, 0, len(episodes))
	for _, e := range episodes {
		if e.IsDir() {
			ents, err := os.ReadDir(e.Name())
			if err != nil {
				return err
			}
			for _, ep := range ents {
				info, err := ep.Info()
				if err != nil {
					return err
				}
				eps = append(eps, info)
			}
		} else {
			eps = append(eps, e)
		}
	}
	if err := validateEps(eps); err != nil {
		return err
	}
	if err := sortEps(eps); err != nil {
		return err
	}
	ents, err := os.ReadDir(seasonDir)
	if err != nil {
		return err
	}
	prevEp := ents[len(ents)-1].Name()
	i := strings.Index(prevEp, "E")
	j := strings.Index(prevEp, ".")
	if i == -1 || j == -1 || j >= i {
		return fmt.Errorf("invalid episode %q", prevEp)
	}
	epn, err := strconv.Atoi(prevEp[i+1 : j])
	if err != nil {
		return err
	}
	for _, e := range eps {
		epn++
		name := e.Name()
		ep := fmt.Sprintf("S%sE%02d%s", season, epn, filepath.Ext(name))
		if err := os.Rename(name, filepath.Join(seasonDir, ep)); err != nil {
			return err
		}
	}
	return nil
}

func validateEps(eps []fs.FileInfo) error {
	for _, e := range eps {
		name := e.Name()
		if _, err := strconv.Atoi(re.FindString(name)); err != nil {
			return fmt.Errorf("episode %q must contain number: %w", name, err)
		}
	}
	return nil
}

func sortEps(eps []fs.FileInfo) error {
	var err2 error
	slices.SortFunc(eps, func(a, b fs.FileInfo) int {
		e1, err := strconv.Atoi(re.FindString(a.Name()))
		if err != nil {
			err2 = err
			return 0
		}
		e2, err := strconv.Atoi(re.FindString(b.Name()))
		if err != nil {
			err2 = err
			return 0
		}
		return cmp.Compare(e1, e2)
	})
	return err2
}
