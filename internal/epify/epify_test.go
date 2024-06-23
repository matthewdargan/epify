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
			name:    "valid show",
			show:    &Show{Name: "The Office", Year: "2005", TVDBID: "73244"},
			wantErr: false,
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
				path := fmt.Sprintf("%s (%s) [tvdbid-%s]", tt.show.Name, tt.show.Year, tt.show.TVDBID)
				want := filepath.Join(tt.show.Dir, path)
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
		name    string
		season  *Season
		wantErr bool
	}{
		{
			name:    "invalid season number",
			season:  &Season{N: "three"},
			wantErr: true,
		},
		{
			name:    "invalid show directory",
			season:  &Season{N: "3", ShowDir: "nonexistent-dir"},
			wantErr: true,
		},
		{
			name:    "show directory not a directory",
			season:  &Season{N: "3", ShowDir: "doc.go"},
			wantErr: true,
		},
		{
			name:    "no episodes",
			season:  &Season{N: "3"},
			wantErr: true,
		},
		{
			name:    "invalid episode",
			season:  &Season{N: "3", Episodes: []string{"nonexistent.mkv"}},
			wantErr: true,
		},
		{
			name:    "invalid episode directory",
			season:  &Season{N: "3", Episodes: []string{"nonexistent-dir"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.season.ShowDir == "" {
				dir, err := os.MkdirTemp("", "show")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				tt.season.ShowDir = dir
			}
			err := MkSeason(tt.season)
			if (err != nil) != tt.wantErr {
				t.Errorf("MkSeason(%v) error = %v", tt.season, err)
			}
		})
	}
}
