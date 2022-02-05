package download

import "testing"

func TestClean(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"github.com/foo/bar", "github.com/foo/bar"},
		{"https://github.com/foo/bar", "github.com/foo/bar"},
		{"https://user@github.com/foo/bar", "github.com/foo/bar"},
		{"github.com/foo/bar/v2", "github.com/foo/bar/v2"},
	}

	for _, tt := range cases {
		got, err := Clean(tt.path)
		if err != nil {
			t.Fatal(err)
		}

		if got != tt.want {
			t.Errorf("Clean(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}
