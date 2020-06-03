package httpapi

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/yeqown/log"
)

// MakeHandler .
func MakeHandler(name string, fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validPath := regexp.MustCompile(fmt.Sprintf(`^/%s/([a-zA-Z0-9\-_\/\.]+)$`, name))

		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
			return
		}
		if len(m) < 1 || m[1] == "" {
			http.Error(w, "Please enter a repository", http.StatusBadRequest)
			return
		}

		repo := m[1]

		// for backwards-compatibility, we must support URLs formatted as
		//   /report/[org]/[repo]
		// and they will be assumed to be github.com URLs. This is because
		// at first Go Report Card only supported github.com URLs, and
		// took only the org name and repo name as parameters. This is no longer the
		// case, but we do not want external links to break.
		oldFormat := regexp.MustCompile(fmt.Sprintf(`^/%s/([a-zA-Z0-9\-_]+)/([a-zA-Z0-9\-_]+)$`, name))
		m2 := oldFormat.FindStringSubmatch(r.URL.Path)
		if m2 != nil {
			// old format is being used
			repo = "github.com/" + repo
			log.Infof("Assuming intended repo is %q, redirecting", repo)
			http.Redirect(w, r, fmt.Sprintf("/%s/%s", name, repo), http.StatusMovedPermanently)
			return
		}

		fn(w, r, repo)
	}
}
