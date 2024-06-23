// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epify

import (
	"os"
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
			dir, err := os.MkdirTemp("", "shows")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			tt.show.Dir = dir
			err = MkShow(tt.show)
			if (err != nil) != tt.wantErr {
				t.Errorf("MkShow(%v) error = %v", tt.show, err)
			}
		})
	}
}
