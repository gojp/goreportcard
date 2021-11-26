package download

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModuleName(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL
		switch u.Path {
		case "/github.com/user/module/@latest":
			fmt.Fprintf(w, `{"Version":"v0.1.0","Time":"2019-08-07T08:30:46Z"}`)
			return
		case "/github.com/user/module/@v/v0.1.0.mod":
			fmt.Fprintf(w, `module github.com/user/module`)
			return
		}
	}))
	defer ts.Close()

	c := NewProxyClient(ts.URL)

	got, err := c.ModuleName("github.com/user/module")
	if err != nil {
		t.Fatal(err)
	}

	want := "github.com/user/module"
	if got != want {
		t.Errorf("got module name = %q, want %q", got, want)
	}
}

func TestLatestVersion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"Version":"v0.1.0","Time":"2019-08-07T08:30:46Z"}`)
	}))
	defer ts.Close()

	c := NewProxyClient(ts.URL)

	got, err := c.LatestVersion("github.com/user/module")
	if err != nil {
		t.Fatal(err)
	}

	want := "v0.1.0"
	if got != want {
		t.Errorf("got latest version = %q, want %q", got, want)
	}
}
