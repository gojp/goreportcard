package download

import (
	"os"
	"path/filepath"
	"testing"
)

var testDownloadDir = "test_downloads"

func TestRepoRootForImportPath(t *testing.T) {
	t.Skip()
	cases := []struct {
		giveURL  string
		wantPath string
		wantVCS  string
	}{
		{"github.com/gojp/goreportcard", "github.com/gojp/goreportcard", "git"},
		{"https://github.com/boltdb/bolt", "github.com/boltdb/bolt", "git"},
		{"https://bitbucket.org/rickb777/go-talk", "bitbucket.org/rickb777/go-talk", "hg"},
		{"ssh://hg@bitbucket.org/rickb777/go-talk", "bitbucket.org/rickb777/go-talk", "hg"},
	}

	for _, tt := range cases {
		root, err := Download(tt.giveURL, testDownloadDir)
		if err != nil {
			t.Fatalf("Error calling Download(%q): %v", tt.giveURL, err)
		}

		if root.Root != tt.wantPath {
			t.Errorf("Download(%q): root.Repo = %q, want %q", tt.giveURL, root.Repo, tt.wantPath)
		}

		wantPath := filepath.Join(testDownloadDir, tt.wantPath)
		ex, _ := exists(wantPath)
		if !ex {
			t.Errorf("Download(%q): %q was not created", tt.giveURL, wantPath)
		}
	}

	// clean up the test
	os.RemoveAll(testDownloadDir)
}
