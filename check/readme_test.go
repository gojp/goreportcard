package check

import "testing"

func TestReadmePercentage(t *testing.T) {
	g := Readme{"testfiles", []string{}}
	p, _, err := g.Percentage()
	if err != nil {
		t.Fatal(err)
	}
	if p != 1.0 {
		t.Errorf("Readme check failed")
	}
}
