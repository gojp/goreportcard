package check

import (
	"testing"
)

func TestRun(t *testing.T) {
	cr, err := Run("testdata/testrepo@v0.1.0", false)
	if err != nil {
		t.Fatal(err)
	}

	if cr.Issues != 2 {
		t.Errorf("got cr.Issues = %d, want %d", cr.Issues, 2)
	}
}
