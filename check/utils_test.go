package check

import (
	"reflect"
	"testing"
)

func TestGoFiles(t *testing.T) {
	files, err := GoFiles("testfiles/")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"testfiles/a.go", "testfiles/b.go", "testfiles/c.go"}
	if !reflect.DeepEqual(files, want) {
		t.Errorf("GoFiles(%q) = %v, want %v", "testfiles/", files, want)
	}
}
