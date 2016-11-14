package check

import "testing"

func TestPercentage(t *testing.T) {
	g := License{"testfiles", []string{}}
	p, _, err := g.Percentage()
	if err != nil {
		t.Fatal(err)
	}
	if p != 1.0 {
		t.Errorf("License check failed")
	}
}
