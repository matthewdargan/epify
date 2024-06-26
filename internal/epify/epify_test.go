// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epify

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
		show    *Show
		wantErr bool
		path    string
	}{
		{
			name:    "empty name",
			show:    &Show{},
			wantErr: true,
		},
		{
			name:    "invalid year",
			show:    &Show{Name: "The Office", Year: "two thousand and five"},
			wantErr: true,
		},
		{
			name:    "invalid tvdbid",
			show:    &Show{Name: "The Office", Year: "2005", TVDBID: "seven three two four four"},
			wantErr: true,
		},
		{
			name: "valid show",
			show: &Show{Name: "The Office", Year: "2005", TVDBID: "73244"},
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
			tt.show.Dir = dir
			err = MkShow(tt.show)
			if (err != nil) != tt.wantErr {
				t.Errorf("MkShow(%v) error = %v", tt.show, err)
			}
			if !tt.wantErr {
				want := filepath.Join(tt.show.Dir, tt.path)
				if _, err := os.Stat(want); os.IsNotExist(err) {
					t.Errorf("MkShow(%v) = %v, want %v", tt.show, err, want)
				}
			}
		})
	}
}

func TestMkSeason(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		season    *Season
		wantErr   bool
		cDir      bool
		cEpisodes bool
	}{
		{
			name:    "invalid season number",
			season:  &Season{N: "three"},
			wantErr: true,
		},
		{
			name:    "invalid show directory",
			season:  &Season{N: "3", ShowDir: "nonexistentdir"},
			wantErr: true,
		},
		{
			name:    "show file",
			season:  &Season{N: "3", ShowDir: "doc.go"},
			wantErr: true,
		},
		{
			name:    "show directory missing name",
			season:  &Season{N: "3", ShowDir: "(2005) [tvdbid-73244]"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "show directory missing year",
			season:  &Season{N: "3", ShowDir: "The Office [tvdbid-73244]"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "show directory missing space before year",
			season:  &Season{N: "3", ShowDir: "The Office(2005) [tvdbid-73244]"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "no episodes",
			season:  &Season{N: "3", ShowDir: "The Office (2005) [tvdbid-73244]"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "invalid episode",
			season:  &Season{N: "3", ShowDir: "Game of Thrones (2011) [tvdbid-121361]", Episodes: []string{"nonexistent.mkv"}},
			wantErr: true,
			cDir:    true,
		},
		{
			name:      "episode directory",
			season:    &Season{N: "3", ShowDir: "Breaking Bad (2008) [tvdbid-81189]", Episodes: []string{"epdir"}},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
		},
		{
			name:      "episode without number",
			season:    &Season{N: "3", ShowDir: "One Piece (1999) [tvdbid-81797]", Episodes: []string{"epx.mkv"}},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
		},
		{
			name: "valid season 3",
			season: &Season{N: "3", ShowDir: "Dragon Ball (1986) [tvdbid-76666]", Episodes: []string{
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
			season:    &Season{N: "11", ShowDir: "Steins;Gate (2011) [tvdbid-244061]", Episodes: []string{"ep9.mp4", "ep10.mp4"}},
			cDir:      true,
			cEpisodes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.cDir {
				dir := filepath.Join(os.TempDir(), tt.season.ShowDir)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				tt.season.ShowDir = dir
			}
			if tt.cEpisodes {
				dir, err := os.MkdirTemp("", "season")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				createEpisodes(t, dir, tt.season.Episodes)
			}
			err := MkSeason(tt.season)
			if (err != nil) != tt.wantErr {
				t.Errorf("MkSeason(%v) error = %v", tt.season, err)
			}
			if !tt.wantErr {
				if len(tt.season.N) < 2 {
					tt.season.N = "0" + tt.season.N
				}
				seasonDir := filepath.Join(tt.season.ShowDir, fmt.Sprintf("Season %s", tt.season.N))
				if _, err := os.Stat(seasonDir); os.IsNotExist(err) {
					t.Errorf("MkSeason(%v) = %v, want %v", tt.season, err, seasonDir)
				}
				ents, err := os.ReadDir(seasonDir)
				if err != nil {
					t.Fatal(err)
				}
				got := len(ents)
				want := len(tt.season.Episodes)
				if got != want {
					t.Errorf("MkSeason(%v) = %v, want %v", tt.season, got, want)
				}
			}
		})
	}
}

func TestAddEpisodes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		add          *SeasonAddition
		wantErr      bool
		cDir         bool
		cEpisodes    bool
		showDir      string
		prevEpisodes []string
	}{
		{
			name:    "invalid season directory",
			add:     &SeasonAddition{SeasonDir: "nonexistentdir"},
			wantErr: true,
		},
		{
			name:    "season file",
			add:     &SeasonAddition{SeasonDir: "doc.go"},
			wantErr: true,
		},
		{
			name:    "season directory without prefix",
			add:     &SeasonAddition{SeasonDir: "noprefix"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "invalid season number",
			add:     &SeasonAddition{SeasonDir: "Season three"},
			wantErr: true,
			cDir:    true,
		},
		{
			name:    "show directory missing name",
			add:     &SeasonAddition{SeasonDir: "Season 03"},
			wantErr: true,
			cDir:    true,
			showDir: "(2011) [tvdbid-121361]",
		},
		{
			name:    "show directory missing year",
			add:     &SeasonAddition{SeasonDir: "Season 03"},
			wantErr: true,
			cDir:    true,
			showDir: "Game of Thrones [tvdbid-121361]",
		},
		{
			name:    "show directory missing space before year",
			add:     &SeasonAddition{SeasonDir: "Season 03"},
			wantErr: true,
			cDir:    true,
			showDir: "Game of Thrones(2011) [tvdbid-121361]",
		},
		{
			name:    "no episodes",
			add:     &SeasonAddition{SeasonDir: "Season 03"},
			wantErr: true,
			cDir:    true,
			showDir: "Game of Thrones (2011) [tvdbid-121361]",
		},
		{
			name:    "invalid episode",
			add:     &SeasonAddition{SeasonDir: "Season 03", Episodes: []string{"nonexistent.mkv"}},
			wantErr: true,
			cDir:    true,
			showDir: "Cowboy Bebop (1998) [tvdbid-76885]",
		},
		{
			name:      "episode directory",
			add:       &SeasonAddition{SeasonDir: "Season 03", Episodes: []string{"epdir"}},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
			showDir:   "Neon Genesis Evangelion (1995) [tvdbid-70350]",
		},
		{
			name:      "episode without number",
			add:       &SeasonAddition{SeasonDir: "Season 03", Episodes: []string{"epx.mkv"}},
			wantErr:   true,
			cDir:      true,
			cEpisodes: true,
			showDir:   "Yu Yu Hakusho (1992) [tvdbid-76665]",
		},
		{
			name:         "previous episode missing E",
			add:          &SeasonAddition{SeasonDir: "Season 10", Episodes: []string{"ep1.mkv"}},
			wantErr:      true,
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Fullmetal Alchemist (2003) [tvdbid-75579]",
			prevEpisodes: []string{"Fullmetal Alchemist S0301.mkv"},
		},
		{
			name:         "previous episode missing period",
			add:          &SeasonAddition{SeasonDir: "Season 10", Episodes: []string{"ep1.mkv"}},
			wantErr:      true,
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Fist of the North Star (1984) [tvdbid-79156]",
			prevEpisodes: []string{"Fist of the North Star S03E01mkv"},
		},
		{
			name:         "previous episode malformed",
			add:          &SeasonAddition{SeasonDir: "Season 10", Episodes: []string{"ep1.mkv"}},
			wantErr:      true,
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Berserk (1997) [tvdbid-73752]",
			prevEpisodes: []string{"Berserk S03.01Emkv"},
		},
		{
			name:         "previous episode invalid number",
			add:          &SeasonAddition{SeasonDir: "Season 10", Episodes: []string{"ep1.mkv"}},
			wantErr:      true,
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Vinland Saga (2019) [tvdbid-359274]",
			prevEpisodes: []string{"Vinland Saga S03E0A.mkv"},
		},
		{
			name: "add to season 3",
			add: &SeasonAddition{SeasonDir: "Season 3", Episodes: []string{
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
			add:          &SeasonAddition{SeasonDir: "Season 199", Episodes: []string{"ep102.avi", "ep103.mkv", "ep104.mp4"}},
			cDir:         true,
			cEpisodes:    true,
			showDir:      "Dragon Ball GT (1996) [tvdbid-79275]",
			prevEpisodes: []string{"Dragon Ball GT S199E01.mp4", "Dragon Ball GT S199E02.mkv", "Dragon Ball GT S199E03.avi"},
		},
		{
			name:      "new episodes",
			add:       &SeasonAddition{SeasonDir: "Season 00", Episodes: []string{"ep9.avi", "ep10.avi"}},
			cDir:      true,
			cEpisodes: true,
			showDir:   "Dragon Ball Super (2015) [tvdbid-295068]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.cDir {
				dir := filepath.Join(os.TempDir(), tt.showDir, tt.add.SeasonDir)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				base := dir
				if tt.showDir != "" {
					base = filepath.Join(os.TempDir(), tt.showDir)
				}
				defer os.RemoveAll(base)
				tt.add.SeasonDir = dir
				createEpisodes(t, tt.add.SeasonDir, tt.prevEpisodes)
			}
			if tt.cEpisodes {
				dir, err := os.MkdirTemp("", "season")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				createEpisodes(t, dir, tt.add.Episodes)
			}
			err := AddEpisodes(tt.add)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddEpisodes(%v) error = %v", tt.add, err)
			}
			if !tt.wantErr {
				if _, err := os.Stat(tt.add.SeasonDir); os.IsNotExist(err) {
					t.Errorf("AddEpisodes(%v) = %v, want %v", tt.add, err, tt.add.SeasonDir)
				}
				ents, err := os.ReadDir(tt.add.SeasonDir)
				if err != nil {
					t.Fatal(err)
				}
				got := len(ents)
				want := len(tt.prevEpisodes) + len(tt.add.Episodes)
				if got != want {
					t.Errorf("AddEpisodes(%v) = %v, want %v", tt.add, got, want)
				}
			}
		})
	}
}

func createEpisodes(t *testing.T, seasonDir string, eps []string) {
	for i, e := range eps {
		path := filepath.Join(seasonDir, e)
		eps[i] = path
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
}
