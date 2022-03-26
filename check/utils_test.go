package check

import (
	"reflect"
	"testing"
)

func TestGoFiles(t *testing.T) {
	files, skipped, err := GoFiles("testdata/testfiles/")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"testdata/testfiles/a.go", "testdata/testfiles/b.go", "testdata/testfiles/c.go", "testdata/testfiles/d.go"}
	if !reflect.DeepEqual(files, want) {
		t.Errorf("GoFiles(%q) = %v, want %v", "testdata/testfiles/", files, want)
	}

	wantSkipped := []string{"testdata/testfiles/a.pb.go", "testdata/testfiles/a.pb.gw.go"}
	if !reflect.DeepEqual(skipped, wantSkipped) {
		t.Errorf("GoFiles(%q) skipped = %v, want %v", "testdata/testfiles/", skipped, wantSkipped)
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
	{"go vet", "testdata/testfiles", []string{"testdata/testfiles/a.go", "testdata/testfiles/b.go", "testdata/testfiles/c.go"}, []string{"go", "vet"}, 1, []FileSummary{}, false},
	{
		"staticcheck",
		"testdata/testfiles/",
		[]string{"testdata/testfiles/a.go", "testdata/testfiles/b.go", "testdata/testfiles/c.go", "testdata/testfiles/d.go"},
		[]string{"staticcheck", "./..."},
		0.75,
		[]FileSummary{
			{
				Filename: "testdata/testfiles/d.go", FileURL: "",
				Errors: []Error{
					{LineNumber: 8, ErrorString: " func foo is unused (U1000)"},
					{LineNumber: 10, ErrorString: " should use time.Until instead of t.Sub(time.Now()) (S1024)"}},
			},
		},
		false,
	},
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
			t.Errorf("[%s] GoTool failed = %#v, want %v", tt.name, fs, tt.failed)
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
		fn   string
		want string
	}{
		{"/github.com/foo/testrepo@v0.1.0/bar/baz.go", "https://github.com/foo/testrepo/blob/v0.1.0/bar/baz.go"},
		{"/github.com/foo/testrepo@v0.0.0-20211126063219-a5e10ccf946a/bar/baz.go", "https://github.com/foo/testrepo/blob/a5e10ccf946a/bar/baz.go"},
		{"/github.com/foo/testrepo@v20.10.11+incompatible/bar/baz.go", "https://github.com/foo/testrepo/blob/v20.10.11/bar/baz.go"},
		{"/github.com/foo/testrepo@v0.1.0-alpha/bar/baz.go", "https://github.com/foo/testrepo/blob/v0.1.0-alpha/bar/baz.go"},
	}

	for _, tt := range cases {
		if got := fileURL(tt.fn); got != tt.want {
			t.Errorf("fileURL(%q) = %q, want %q", tt.fn, got, tt.want)
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
