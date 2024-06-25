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
//	$ epify season 3 '/media/shows/The Office (2005) [tvdbid-73244]' /downloads/the_office_s3_p1/ep*.mkv
//
// Add episodes to `/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:
//
//	$ epify add '/media/shows/The Office (2005) [tvdbid-73244]/Season 03' /downloads/the_office_s3_p2/ep*.mkv
//
// [Jellyfin naming scheme]: https://jellyfin.org/docs/general/server/media/shows/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/matthewdargan/epify/internal/epify"
)

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
		s := epify.Show{
			Name:   args[1],
			Year:   args[2],
			TVDBID: args[3],
			Dir:    args[4],
		}
		if err := epify.MkShow(&s); err != nil {
			log.Fatal(err)
		}
	case "season":
		if flag.NArg() < 4 {
			usage()
		}
		s := epify.Season{
			N:        args[1],
			ShowDir:  args[2],
			Episodes: args[3:],
		}
		if err := epify.MkSeason(&s); err != nil {
			log.Fatal(err)
		}
	case "add":
		if flag.NArg() < 3 {
			usage()
		}
		s := epify.SeasonAddition{
			SeasonDir: args[1],
			Episodes:  args[2:],
		}
		if err := epify.AddEpisodes(&s); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}
