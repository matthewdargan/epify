// Copyright 2024 Matthew P. Dargan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package test provides utilities for testing.
package test

import (
	"os"
	"path/filepath"
	"testing"
)

// SetupFiles creates files and directories for testing.
func SetupFiles(t *testing.T, dir string, fs ...string) []string {
	t.Helper()
	ps := make([]string, len(fs))
	for i, m := range fs {
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
