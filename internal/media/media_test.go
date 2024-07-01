// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package media

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMkShow(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		s       *Show
		wantErr bool
		path    string
	}{
		{
			name:    "empty name",
			s:       &Show{},
			wantErr: true,
		},
		{
			name:    "invalid year",
			s:       &Show{Name: "The Office", Year: "two thousand and five"},
			wantErr: true,
		},
		{
			name:    "invalid tvdbid",
			s:       &Show{Name: "The Office", Year: "2005", ID: "seven three two four four"},
			wantErr: true,
		},
		{
			name: "valid show",
			s:    &Show{Name: "The Office", Year: "2005", ID: "73244"},
			path: "The Office (2005) [tvdbid-73244]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir, err := os.MkdirTemp("", "show")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			tt.s.Dir = dir
			err = MkShow(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("MkShow(%v) error = %v", tt.s, err)
			}
			if !tt.wantErr {
				want := filepath.Join(tt.s.Dir, tt.path)
				if _, err := os.Stat(want); os.IsNotExist(err) {
					t.Errorf("MkShow(%v) = %v, want %v", tt.s, err, want)
				}
			}
		})
	}
}

func TestAddMovie(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		m       *Movie
		wantErr bool
		cDir    bool
		cMovie  bool
		path    string
	}{
		{
			name:    "empty name",
			m:       &Movie{},
			wantErr: true,
		},
		{
			name:    "invalid year",
			m:       &Movie{Show: Show{Name: "Braveheart", Year: "nineteen ninety five"}},
			wantErr: true,
		},
		{
			name:    "invalid tmdbid",
			m:       &Movie{Show: Show{Name: "Braveheart", Year: "2005", ID: "one nine seven"}},
			wantErr: true,
		},
		{
			name:    "invalid directory",
			m:       &Movie{Show: Show{Name: "Braveheart", Year: "2005", ID: "197", Dir: "nonexistentdir"}},
			wantErr: true,
		},
		{
			name:    "directory file",
			m:       &Movie{Show: Show{Name: "Braveheart", Year: "2005", ID: "197", Dir: "doc.go"}},
			wantErr: true,
		},
		{
			name:    "invalid movie",
			m:       &Movie{Show: Show{Name: "Braveheart", Year: "2005", ID: "197"}, File: "nonexistent.mkv"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "movie directory",
			m:       &Movie{Show: Show{Name: "Braveheart", Year: "2005", ID: "197"}, File: "moviedir"},
			wantErr: true,
			cDir:    true,
			cMovie:  true,
		},
		{
			name:   "valid movie",
			m:      &Movie{Show: Show{Name: "Braveheart", Year: "2005", ID: "197"}, File: "braveheart.mkv"},
			cDir:   true,
			cMovie: true,
			path:   "Braveheart (2005) [tmdbid-197].mkv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.cDir {
				dir, err := os.MkdirTemp("", "movie")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				tt.m.Dir = dir
			}
			if tt.cMovie {
				dir, err := os.MkdirTemp("", "download")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				tt.m.File = createMedia(t, dir, tt.m.File)[0]
			}
			err := AddMovie(tt.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddMovie(%v) error = %v", tt.m, err)
			}
			if !tt.wantErr {
				want := filepath.Join(tt.m.Dir, tt.path)
				if _, err := os.Stat(want); os.IsNotExist(err) {
					t.Errorf("AddMovie(%v) = %v, want %v", tt.m, err, want)
				}
			}
		})
	}
}

func TestMkSeason(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		s         *Season
		wantErr   bool
		cDir      bool
		cEpisodes bool
	}{
		{
			name:    "invalid season number",
			s:       &Season{N: "three"},
			wantErr: true,
		},
		{
			name:    "invalid directory",
			s:       &Season{N: "3", ShowDir: "nonexistentdir"},
			wantErr: true,
		},
		{
			name:    "show file",
			s:       &Season{N: "3", ShowDir: "doc.go"},
			wantErr: true,
		},
		{
			name:    "directory missing name",
			s:       &Season{N: "3", ShowDir: "(2005) [tvdbid-73244]"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "directory missing year",
			s:       &Season{N: "3", ShowDir: "The Office [tvdbid-73244]"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "directory missing space before year",
			s:       &Season{N: "3", ShowDir: "The Office(2005) [tvdbid-73244]"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "no episodes",
			s:       &Season{N: "3", ShowDir: "The Office (2005) [tvdbid-73244]"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "invalid episode",
			s:       &Season{N: "3", ShowDir: "Game of Thrones (2011) [tvdbid-121361]", Episodes: []string{"nonexistent.mkv"}},
			wantErr: true,
			cDir:    true,
		},
		{
			name:      "episode directory",
			s:         &Season{N: "3", ShowDir: "Breaking Bad (2008) [tvdbid-81189]", Episodes: []string{"epdir"}},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
		},
		{
			name:      "episode without number",
			s:         &Season{N: "3", ShowDir: "One Piece (1999) [tvdbid-81797]", Episodes: []string{"epx.mkv"}},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
		},
		{
			name:      "negative match index",
			s:         &Season{N: "0", ShowDir: "Naruto (2002) [tvdbid-78857]", Episodes: []string{"ep1.mkv"}, MatchIndex: -1},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
		},
		{
			name:      "match index 1 out of range",
			s:         &Season{N: "3", ShowDir: "Naruto Shippuden (2007) [tvdbid-79824]", Episodes: []string{"ep1.mkv"}, MatchIndex: 1},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
		},
		{
			name:      "match index 2 out of range",
			s:         &Season{N: "300", ShowDir: "Samurai Champloo (2004) [tvdbid-79089]", Episodes: []string{"s1ep2.mkv"}, MatchIndex: 2},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
		},
		{
			name: "valid season 3",
			s: &Season{N: "3", ShowDir: "Dragon Ball (1986) [tvdbid-76666]", Episodes: []string{
				"ep1.mkv", "ep2.mkv", "ep3.mkv", "ep4.mkv", "ep5.mkv",
				"ep6.mkv", "ep7.mkv", "ep8.mkv", "ep9.mkv", "ep10.mkv",
				"ep11.mkv", "ep12.mkv", "ep13.mkv", "ep14.mkv", "ep15.mkv",
				"ep16.mkv", "ep17.mkv", "ep18.mkv", "ep19.mkv", "ep20.mkv",
				"ep21.mkv", "ep22.mkv", "ep23.mkv", "ep24.mkv", "ep25.mkv",
				"ep26.mkv", "ep27.mkv", "ep28.mkv", "ep29.mkv", "ep30.mkv",
				"ep31.mkv", "ep32.mkv", "ep33.mkv", "ep34.mkv", "ep35.mkv",
				"ep36.mkv", "ep37.mkv", "ep38.mkv", "ep39.mkv", "ep40.mkv",
				"ep41.mkv", "ep42.mkv", "ep43.mkv", "ep44.mkv", "ep45.mkv",
				"ep46.mkv", "ep47.mkv", "ep48.mkv", "ep49.mkv", "ep50.mkv",
				"ep51.mkv", "ep52.mkv", "ep53.mkv", "ep54.mkv", "ep55.mkv",
				"ep56.mkv", "ep57.mkv", "ep58.mkv", "ep59.mkv", "ep60.mkv",
				"ep61.mkv", "ep62.mkv", "ep63.mkv", "ep64.mkv", "ep65.mkv",
				"ep66.mkv", "ep67.mkv", "ep68.mkv", "ep69.mkv", "ep70.mkv",
				"ep71.mkv", "ep72.mkv", "ep73.mkv", "ep74.mkv", "ep75.mkv",
				"ep76.mkv", "ep77.mkv", "ep78.mkv", "ep79.mkv", "ep80.mkv",
				"ep81.mkv", "ep82.mkv", "ep83.mkv", "ep84.mkv", "ep85.mkv",
				"ep86.mkv", "ep87.mkv", "ep88.mkv", "ep89.mkv", "ep90.mkv",
				"ep91.mkv", "ep92.mkv", "ep93.mkv", "ep94.mkv", "ep95.mkv",
				"ep96.mkv", "ep97.mkv", "ep98.mkv", "ep99.mkv", "ep100.mkv",
				"ep101.mkv",
			}},
			cDir:      true,
			cEpisodes: true,
		},
		{
			name:      "valid season 11",
			s:         &Season{N: "11", ShowDir: "Steins;Gate (2011) [tvdbid-244061]", Episodes: []string{"ep9.mp4", "ep10.mp4"}},
			cDir:      true,
			cEpisodes: true,
		},
		{
			name: "match index 1",
			s: &Season{
				N:          "0",
				ShowDir:    "Attack on Titan (2013) [tvdbid-514059]",
				Episodes:   []string{"Attack on Titan S00E16.mkv", "Attack on Titan S00E15.mkv", "Attack on Titan S00E14.mkv"},
				MatchIndex: 1,
			},
			cDir:      true,
			cEpisodes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.cDir {
				dir := filepath.Join(os.TempDir(), tt.s.ShowDir)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				tt.s.ShowDir = dir
			}
			if tt.cEpisodes {
				dir, err := os.MkdirTemp("", "season")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				tt.s.Episodes = createMedia(t, dir, tt.s.Episodes...)
			}
			err := MkSeason(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("MkSeason(%v) error = %v", tt.s, err)
			}
			if !tt.wantErr {
				if len(tt.s.N) < 2 {
					tt.s.N = "0" + tt.s.N
				}
				seasonDir := filepath.Join(tt.s.ShowDir, fmt.Sprintf("Season %s", tt.s.N))
				if _, err := os.Stat(seasonDir); os.IsNotExist(err) {
					t.Errorf("MkSeason(%v) = %v, want %v", tt.s, err, seasonDir)
				}
				ents, err := os.ReadDir(seasonDir)
				if err != nil {
					t.Fatal(err)
				}
				got := len(ents)
				want := len(tt.s.Episodes)
				if got != want {
					t.Errorf("MkSeason(%v) = %v, want %v", tt.s, got, want)
				}
			}
		})
	}
}

func TestAddEpisodes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		a            *Addition
		wantErr      bool
		cDir         bool
		cEpisodes    bool
		showDir      string
		prevEpisodes []string
	}{
		{
			name:    "invalid season directory",
			a:       &Addition{SeasonDir: "nonexistentdir"},
			wantErr: true,
		},
		{
			name:    "season file",
			a:       &Addition{SeasonDir: "doc.go"},
			wantErr: true,
		},
		{
			name:    "season directory without prefix",
			a:       &Addition{SeasonDir: "noprefix"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "invalid season number",
			a:       &Addition{SeasonDir: "Season three"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "show directory missing name",
			a:       &Addition{SeasonDir: "Season 03"},
			wantErr: true,
			cDir:    true,
			showDir: "(2011) [tvdbid-121361]",
		},
		{
			name:    "show directory missing year",
			a:       &Addition{SeasonDir: "Season 03"},
			wantErr: true,
			cDir:    true,
			showDir: "Game of Thrones [tvdbid-121361]",
		},
		{
			name:    "show directory missing space before year",
			a:       &Addition{SeasonDir: "Season 03"},
			wantErr: true,
			cDir:    true,
			showDir: "Game of Thrones(2011) [tvdbid-121361]",
		},
		{
			name:    "no episodes",
			a:       &Addition{SeasonDir: "Season 03"},
			wantErr: true,
			cDir:    true,
			showDir: "Game of Thrones (2011) [tvdbid-121361]",
		},
		{
			name:    "invalid episode",
			a:       &Addition{SeasonDir: "Season 03", Episodes: []string{"nonexistent.mkv"}},
			wantErr: true,
			cDir:    true,
			showDir: "Cowboy Bebop (1998) [tvdbid-76885]",
		},
		{
			name:      "episode directory",
			a:         &Addition{SeasonDir: "Season 03", Episodes: []string{"epdir"}},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
			showDir:   "Neon Genesis Evangelion (1995) [tvdbid-70350]",
		},
		{
			name:      "episode without number",
			a:         &Addition{SeasonDir: "Season 03", Episodes: []string{"epx.mkv"}},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
			showDir:   "Yu Yu Hakusho (1992) [tvdbid-76665]",
		},
		{
			name:      "negative match index",
			a:         &Addition{SeasonDir: "Season 00", Episodes: []string{"ep10.mkv"}, MatchIndex: -1},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
			showDir:   "Hunter x Hunter (1999) [tvdbid-79076]",
		},
		{
			name:      "match index 1 out of range",
			a:         &Addition{SeasonDir: "Season 03", Episodes: []string{"ep100.mkv"}, MatchIndex: 1},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
			showDir:   "Sailor Moon (1992) [tvdbid-78500]",
		},
		{
			name:      "match index 2 out of range",
			a:         &Addition{SeasonDir: "Season 300", Episodes: []string{"s300ep2.mkv"}, MatchIndex: 2},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
			showDir:   "One-Punch Man (2015) [tvdbid-293088]",
		},
		{
			name:         "previous episode missing E",
			a:            &Addition{SeasonDir: "Season 10", Episodes: []string{"ep1.mkv"}},
			wantErr:      true,
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Fullmetal Alchemist (2003) [tvdbid-75579]",
			prevEpisodes: []string{"Fullmetal Alchemist S1001.mkv"},
		},
		{
			name:         "previous episode missing period",
			a:            &Addition{SeasonDir: "Season 10", Episodes: []string{"ep1.mkv"}},
			wantErr:      true,
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Fist of the North Star (1984) [tvdbid-79156]",
			prevEpisodes: []string{"Fist of the North Star S10E01mkv"},
		},
		{
			name:         "previous episode malformed",
			a:            &Addition{SeasonDir: "Season 10", Episodes: []string{"ep1.mkv"}},
			wantErr:      true,
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Berserk (1997) [tvdbid-73752]",
			prevEpisodes: []string{"Berserk S10.01Emkv"},
		},
		{
			name:         "previous episode invalid number",
			a:            &Addition{SeasonDir: "Season 10", Episodes: []string{"ep1.mkv"}},
			wantErr:      true,
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Vinland Saga (2019) [tvdbid-359274]",
			prevEpisodes: []string{"Vinland Saga S10E0A.mkv"},
		},
		{
			name: "add to season 3",
			a: &Addition{SeasonDir: "Season 3", Episodes: []string{
				"ep1.mp4", "ep2.mp4", "ep3.mp4", "ep4.mp4", "ep5.mp4",
				"ep6.mp4", "ep7.mp4", "ep8.mp4", "ep9.mp4", "ep10.mp4",
				"ep11.mp4", "ep12.mp4", "ep13.mp4", "ep14.mp4", "ep15.mp4",
				"ep16.mp4", "ep17.mp4", "ep18.mp4", "ep19.mp4", "ep20.mp4",
				"ep21.mp4", "ep22.mp4", "ep23.mp4", "ep24.mp4", "ep25.mp4",
				"ep26.mp4", "ep27.mp4", "ep28.mp4", "ep29.mp4", "ep30.mp4",
				"ep31.mp4", "ep32.mp4", "ep33.mp4", "ep34.mp4", "ep35.mp4",
				"ep36.mp4", "ep37.mp4", "ep38.mp4", "ep39.mp4", "ep40.mp4",
				"ep41.mp4", "ep42.mp4", "ep43.mp4", "ep44.mp4", "ep45.mp4",
				"ep46.mp4", "ep47.mp4", "ep48.mp4", "ep49.mp4", "ep50.mp4",
				"ep51.mp4", "ep52.mp4", "ep53.mp4", "ep54.mp4", "ep55.mp4",
				"ep56.mp4", "ep57.mp4", "ep58.mp4", "ep59.mp4", "ep60.mp4",
				"ep61.mp4", "ep62.mp4", "ep63.mp4", "ep64.mp4", "ep65.mp4",
				"ep66.mp4", "ep67.mp4", "ep68.mp4", "ep69.mp4", "ep70.mp4",
				"ep71.mp4", "ep72.mp4", "ep73.mp4", "ep74.mp4", "ep75.mp4",
				"ep76.mp4", "ep77.mp4", "ep78.mp4", "ep79.mp4", "ep80.mp4",
				"ep81.mp4", "ep82.mp4", "ep83.mp4", "ep84.mp4", "ep85.mp4",
				"ep86.mp4", "ep87.mp4", "ep88.mp4", "ep89.mp4", "ep90.mp4",
				"ep91.mp4", "ep92.mp4", "ep93.mp4", "ep94.mp4", "ep95.mp4",
				"ep96.mp4", "ep97.mp4", "ep98.mp4", "ep99.mp4", "ep100.mp4",
				"ep101.mp4",
			}},
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Dragon Ball Z (1989) [tvdbid-81472]",
			prevEpisodes: []string{"Dragon Ball Z S03E01.mkv"},
		},
		{
			name:         "add to season 199",
			a:            &Addition{SeasonDir: "Season 199", Episodes: []string{"ep102.avi", "ep103.mkv", "ep104.mp4"}},
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Defenders of the Earth (1986) [tvdbid-70824]",
			prevEpisodes: []string{"Defenders of the Earth S199E01.mp4", "Defenders of the Earth S199E02.mkv", "Defenders of the Earth S199E03.avi"},
		},
		{
			name:      "new episodes",
			a:         &Addition{SeasonDir: "Season 00", Episodes: []string{"ep9.avi", "ep10.avi"}},
			cDir:      true,
			cEpisodes: true,
			showDir:   "Dragon Ball Super (2015) [tvdbid-295068]",
		},
		{
			name: "match index 2",
			a: &Addition{
				SeasonDir:  "Season 30",
				Episodes:   []string{"Bleach 1S30E04.mkv", "Bleach 2S30E03.mkv", "Bleach 3S30E02.mkv", "Bleach 4S30E01.mkv"},
				MatchIndex: 2,
			},
			cDir:      true,
			cEpisodes: true,
			showDir:   "Bleach (2004) [tvdbid-74796]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.cDir {
				dir := filepath.Join(os.TempDir(), tt.showDir, tt.a.SeasonDir)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				base := dir
				if tt.showDir != "" {
					base = filepath.Join(os.TempDir(), tt.showDir)
				}
				defer os.RemoveAll(base)
				tt.a.SeasonDir = dir
				tt.prevEpisodes = createMedia(t, tt.a.SeasonDir, tt.prevEpisodes...)
			}
			if tt.cEpisodes {
				dir, err := os.MkdirTemp("", "season")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				tt.a.Episodes = createMedia(t, dir, tt.a.Episodes...)
			}
			err := AddEpisodes(tt.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddEpisodes(%v) error = %v", tt.a, err)
			}
			if !tt.wantErr {
				if _, err := os.Stat(tt.a.SeasonDir); os.IsNotExist(err) {
					t.Errorf("AddEpisodes(%v) = %v, want %v", tt.a, err, tt.a.SeasonDir)
				}
				ents, err := os.ReadDir(tt.a.SeasonDir)
				if err != nil {
					t.Fatal(err)
				}
				got := len(ents)
				want := len(tt.prevEpisodes) + len(tt.a.Episodes)
				if got != want {
					t.Errorf("AddEpisodes(%v) = %v, want %v", tt.a, got, want)
				}
			}
		})
	}
}

func createMedia(t *testing.T, dir string, ms ...string) []string {
	ps := make([]string, len(ms))
	for i, m := range ms {
		path := filepath.Join(dir, m)
		ps[i] = path
		if filepath.Ext(path) == "" {
			if err := os.Mkdir(path, 0o755); err != nil {
				t.Fatal(err)
			}
			continue
		}
		f, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}
	return ps
}
