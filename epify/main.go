// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Epify categorizes [shows] and [movies] using the Jellyfin naming scheme.
//
// Usage:
//
//	epify show name year tvdbid dir
//	epify movie name year tmdbid dir movie
//	epify season [-m index] seasonnum showdir episode...
//	epify add [-m index] seasondir episode...
//
// `epify show` creates a show directory like
// "Series Name (2018) [tvdbid-65567]".
//
// `epify movie` adds a movie to a directory. Movies are labeled like
// "Film (2018) [tmdbid-65567]".
//
// `epify season` populates a season directory with episodes. Episodes are
// labeled like "Series Name S01E01.mkv".
//
// `epify add` adds episodes to a season directory, continuing at the previous
// episode increment.
//
// The `-m` flag specifies the index of the episode number in filenames for
// the `epify season` and `epify add` commands.
//
// Examples:
//
// Create show directory `/media/shows/The Office (2005) [tvdbid-73244]`:
//
//	$ epify show 'The Office' 2005 73244 '/media/shows'
//
// Add movie to `/media/movies`:
//
//	$ epify movie 'Braveheart' 1995 197 '/media/movies' '/downloads/braveheart.mkv'
//
// Populate season directory
// `/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:
//
//	$ epify season 3 '/media/shows/The Office (2005) [tvdbid-73244]' /downloads/the_office_s3_p1/ep*.mkv
//
// Populate season directory
// `/media/shows/Breaking Bad (2008) [tvdbid-81189]/Season 04`:
//
//	$ epify season -m 1 4 '/media/shows/Breaking Bad (2008) [tvdbid-81189]' /downloads/breaking_bad_s4_p1/s4ep*.mkv
//
// Add episodes to `/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:
//
//	$ epify add '/media/shows/The Office (2005) [tvdbid-73244]/Season 03' /downloads/the_office_s3_p2/ep*.mkv
//
// Add episodes to `/media/shows/Breaking Bad (2008) [tvdbid-81189]/Season 04`:
//
//	$ epify add -m 1 '/media/shows/Breaking Bad (2008) [tvdbid-81189]/Season 04' /downloads/breaking_bad_s4_p2/s4ep*.mkv
//
// [shows]: https://jellyfin.org/docs/general/server/media/shows/
// [movies]: https://jellyfin.org/docs/general/server/media/movies/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/matthewdargan/epify/internal/epify"
)

var (
	seasonCmd   = flag.NewFlagSet("season", flag.ExitOnError)
	seasonMatch = seasonCmd.Int("m", 0, "match index")
	addCmd      = flag.NewFlagSet("add", flag.ExitOnError)
	addMatch    = addCmd.Int("m", 0, "match index")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "\tepify show name year tvdbid dir\n")
	fmt.Fprintf(os.Stderr, "\tepify movie name year tmdbid dir movie\n")
	fmt.Fprintf(os.Stderr, "\tepify season [-m index] seasonnum showdir episode...\n")
	fmt.Fprintf(os.Stderr, "\tepify add [-m index] seasondir episode...\n")
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
		show := epify.Media{
			Name: args[1],
			Year: args[2],
			ID:   args[3],
			Dir:  args[4],
		}
		if err := epify.MkShow(&show); err != nil {
			log.Fatal(err)
		}
	case "movie":
		if flag.NArg() != 6 {
			usage()
		}
		movie := epify.Movie{
			Media: epify.Media{
				Name: args[1],
				Year: args[2],
				ID:   args[3],
				Dir:  args[4],
			},
			File: args[5],
		}
		if err := epify.AddMovie(&movie); err != nil {
			log.Fatal(err)
		}
	case "season":
		if err := seasonCmd.Parse(args[1:]); err != nil {
			log.Fatal(err)
		}
		if seasonCmd.NArg() < 3 {
			usage()
		}
		args = seasonCmd.Args()
		s := epify.Season{
			N:          args[0],
			ShowDir:    args[1],
			Episodes:   args[2:],
			MatchIndex: *seasonMatch,
		}
		if err := epify.MkSeason(&s); err != nil {
			log.Fatal(err)
		}
	case "add":
		if err := addCmd.Parse(args[1:]); err != nil {
			log.Fatal(err)
		}
		if addCmd.NArg() < 2 {
			usage()
		}
		args = addCmd.Args()
		s := epify.SeasonAddition{
			SeasonDir:  args[0],
			Episodes:   args[1:],
			MatchIndex: *addMatch,
		}
		if err := epify.AddEpisodes(&s); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}
