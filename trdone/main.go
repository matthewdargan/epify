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

	"github.com/matthewdargan/epify/internal/torrent"
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
	if flag.NArg() != 1 {
		usage()
	}
	f := torrent.File{
		Dir:    os.Getenv("TR_TORRENT_DIR"),
		Name:   os.Getenv("TR_TORRENT_NAME"),
		DstDir: flag.Arg(0),
	}
	if err := torrent.Rename(&f); err != nil {
		log.Fatal(err)
	}
}
