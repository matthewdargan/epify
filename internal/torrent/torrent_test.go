// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package torrent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matthewdargan/epify/internal/test"
)

func TestRename(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		f          *File
		wantErr    bool
		cDir       bool
		cTorrent   bool
		showDirs   []string
		seasonDirs []string
		path       string
	}{
		{
			name:    "invalid directory",
			f:       &File{DstDir: "nonexistentdir"},
			wantErr: true,
		},
		{
			name:    "directory file",
			f:       &File{DstDir: "torrent.go"},
			wantErr: true,
		},
		{
			name:    "no shows",
			f:       &File{},
			wantErr: true,
			cDir:    true,
		},
		{
			name:     "show file",
			f:        &File{Name: "Cardcaptor Sakura 20.mkv"},
			wantErr:  true,
			cDir:     true,
			showDirs: []string{"Cardcaptor Sakura.mkv"},
		},
		{
			name:     "show directory missing year",
			f:        &File{Name: "Cardcaptor Sakura 20.mkv"},
			wantErr:  true,
			cDir:     true,
			showDirs: []string{"Cardcaptor Sakura [tvdbid-70668]"},
		},
		{
			name:     "no seasons",
			f:        &File{Name: "Cardcaptor Sakura 20.mkv"},
			wantErr:  true,
			cDir:     true,
			showDirs: []string{"Cardcaptor Sakura (1998) [tvdbid-70668]", "My Hero Academia (2016) [tvdbid-305074]"},
		},
		{
			name:       "season file",
			f:          &File{Name: "Jujustu Kaisen 15.mp4"},
			wantErr:    true,
			cDir:       true,
			showDirs:   []string{"Jujustu Kaisen (2020) [tvdbid-377543]"},
			seasonDirs: []string{"season1.mkv"},
		},
		{
			name:       "season directory without prefix",
			f:          &File{Name: "Jujustu Kaisen 15.mp4"},
			wantErr:    true,
			cDir:       true,
			showDirs:   []string{"Jujustu Kaisen (2020) [tvdbid-377543]"},
			seasonDirs: []string{"season1"},
		},
		{
			name:       "invalid season number",
			f:          &File{Name: "Jujustu Kaisen 15.mp4"},
			wantErr:    true,
			cDir:       true,
			showDirs:   []string{"Jujustu Kaisen (2020) [tvdbid-377543]"},
			seasonDirs: []string{"Season one"},
		},
		{
			name:       "invalid episode",
			f:          &File{Name: "Astro Boy 01.mkv"},
			wantErr:    true,
			cDir:       true,
			showDirs:   []string{"Astro Boy (1963) [tvdbid-71952]"},
			seasonDirs: []string{"Season 09", "Season 10", "Season 11"},
		},
		{
			name:       "valid episode",
			f:          &File{Name: "Knights of Sidonia 100.avi"},
			wantErr:    false,
			cDir:       true,
			cTorrent:   true,
			showDirs:   []string{"Knights of Sidonia (2014) [tvdbid-278154]"},
			seasonDirs: []string{"Season 101", "Season 102"},
			path:       "Knights of Sidonia (2014) [tvdbid-278154]/Season 102/Knights of Sidonia S102E01.avi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.cDir {
				dir, err := os.MkdirTemp("", "shows")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				tt.f.DstDir = dir
			}
			if tt.cTorrent {
				dir, err := os.MkdirTemp("", "downloads")
				if err != nil {
					t.Fatal(err)
				}
				defer os.RemoveAll(dir)
				tt.f.Dir = dir
				test.SetupFiles(t, tt.f.Dir, tt.f.Name)
			}
			tt.showDirs = test.SetupFiles(t, tt.f.DstDir, tt.showDirs...)
			// tt.showDirs[0] should be the valid show directory
			if len(tt.showDirs) > 0 {
				tt.seasonDirs = test.SetupFiles(t, tt.showDirs[0], tt.seasonDirs...)
			}
			err := Rename(tt.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rename(%v) error = %v", tt.f, err)
			}
			if !tt.wantErr {
				want := filepath.Join(tt.f.DstDir, tt.path)
				if _, err := os.Stat(want); os.IsNotExist(err) {
					t.Errorf("MkShow(%v) = %v, want %v", tt.f, err, want)
				}
			}
		})
	}
}
