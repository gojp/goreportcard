package handlers

import "testing"

var dirNameTests = []struct {
	url  string
	ver  string
	want string
}{
	{"github.com/foo/bar", "v0.1.0", "_repos/src/github.com/foo/bar@v0.1.0"},
}

func TestDirName(t *testing.T) {
	for _, tt := range dirNameTests {
		if got := dirName(tt.url, tt.ver); got != tt.want {
			t.Errorf("dirName(%q, %q) = %q, want %q", tt.url, tt.ver, got, tt.want)
		}
	}
}
