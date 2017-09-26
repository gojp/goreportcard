package handlers

import (
	"path/filepath"
	"text/template"
)

const (
	delimL = "[["
	delimR = "]]"

	tmplDir = "templates"
	tmplExt = ".html"
)

// template base file names - used in resp. handlers
const (
	tmplHome   = "home"
	tmplAbout  = "about"
	tmplScores = "high_scores"
	tmplReport = "report"
	tmplError  = "404"

	tmplFoot = "footer" // standard footer - available for all pages
)

// parse returns the parsed template associated with the tmplBase name
// (incl. the standard footer),
// or panics iff fail
func parse(tmplBase string, funcMaps ...template.FuncMap) *template.Template {

	name := tmplBase + tmplExt
	foot := tmplFoot + tmplExt

	namePath := filepath.Join(tmplDir, name)
	footPath := filepath.Join(tmplDir, foot)

	t := template.New(name).Delims(delimL, delimR)

	for i := range funcMaps {
		t = t.Funcs(funcMaps[i])
	}

	return template.Must(t.ParseFiles(namePath, footPath))
}
