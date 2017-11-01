package handlers

import (
	"path/filepath"
	"testing"
)

var dirNameTests = []struct {
	url  string
	want string
}{
	{"", "_repos/src"},
	{"github.com/foo/bar", "_repos/src/github.com/foo/bar"},
}

func TestDirName(t *testing.T) {
	for _, tt := range dirNameTests {
		want := filepath.FromSlash(tt.want)
		if got := DirName(tt.url); got != want {
			t.Errorf("dirName(%q) = %q, want %q", tt.url, got, want)
		}
	}
}
