package check

import (
	"reflect"
	"testing"
)

func TestGoFiles(t *testing.T) {
	files, skipped, err := GoFiles("testfiles/")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"testfiles/a.go", "testfiles/b.go", "testfiles/c.go"}
	if !reflect.DeepEqual(files, want) {
		t.Errorf("GoFiles(%q) = %v, want %v", "testfiles/", files, want)
	}

	wantSkipped := []string{"testfiles/a.pb.go", "testfiles/a.pb.gw.go"}
	if !reflect.DeepEqual(skipped, wantSkipped) {
		t.Errorf("GoFiles(%q) skipped = %v, want %v", "testfiles/", skipped, wantSkipped)
	}
}

var goToolTests = []struct {
	name      string
	dir       string
	filenames []string
	tool      []string
	percent   float64
	failed    []FileSummary
	wantErr   bool
}{
	{"go vet", "testfiles/", []string{"testfiles/a.go", "testfiles/b.go", "testfiles/c.go"}, []string{"go", "tool", "vet"}, 1, []FileSummary{}, false},
}

func TestGoTool(t *testing.T) {
	for _, tt := range goToolTests {
		f, fs, err := GoTool(tt.dir, tt.filenames, tt.tool)
		if err != nil && !tt.wantErr {
			t.Fatal(err)
		}
		if f != tt.percent {
			t.Errorf("[%s] GoTool percent = %f, want %f", tt.name, f, tt.percent)
		}
		if !reflect.DeepEqual(fs, tt.failed) {
			t.Errorf("[%s] GoTool failed = %v, want %v", tt.name, fs, tt.failed)
		}
	}
}

func TestMakeFilename(t *testing.T) {
	cases := []struct {
		fn   string
		want string
	}{
		{"/github.com/foo/bar/baz.go", "bar/baz.go"},
	}

	for _, tt := range cases {
		if got := makeFilename(tt.fn); got != tt.want {
			t.Errorf("makeFilename(%q) = %q, want %q", tt.fn, got, tt.want)
		}
	}
}

func TestFileURL(t *testing.T) {
	cases := []struct {
		dir  string
		fn   string
		want string
	}{
		{"_repos/src/github.com/foo/testrepo@v0.1.0/bar/baz.go", "/github.com/foo/testrepo@v0.1.0/bar/baz.go", "https://github.com/foo/testrepo/blob/v0.1.0/bar/baz.go"},
		{"_repos/src/github.com/foo/testrepo@v0.0.0-20211126063219-a5e10ccf946a/bar/baz.go", "/github.com/foo/testrepo@v0.0.0-20211126063219-a5e10ccf946a/bar/baz.go", "https://github.com/foo/testrepo/blob/a5e10ccf946a/bar/baz.go"},
	}

	for _, tt := range cases {
		if got := fileURL(tt.dir, tt.fn); got != tt.want {
			t.Errorf("fileURL(%q, %q) = %q, want %q", tt.dir, tt.fn, got, tt.want)
		}
	}

}

func TestDisplayFilename(t *testing.T) {
	cases := []struct {
		fn   string
		want string
	}{
		{"foo@v0.1.0/bar/baz.go", "bar/baz.go"},
		{"foo@v0.1.0/a/b/c/d/baz.go", "a/b/c/d/baz.go"},
	}

	for _, tt := range cases {
		if got := displayFilename(tt.fn); got != tt.want {
			t.Errorf("displayFilename(%q) = %q, want %q", tt.fn, got, tt.want)
		}
	}
}
