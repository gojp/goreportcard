package handlers

import (
	"fmt"
	"io"
	"text/template"
)

func add(x, y int) int {
	return x + y
}

func formatScore(x float64) string {
	return fmt.Sprintf("%.2f", x)
}

func (gh *GRCHandler) loadTemplate(name string) (*template.Template, error) {
	f, err := gh.AssetsFS.Open(name)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New(name).Delims("[[", "]]").Funcs(template.FuncMap{
		"add":         add,
		"formatScore": formatScore,
	}).Parse(string(contents))
	if err != nil {
		return nil, err
	}

	if name == "/templates/report.html" {
		return tpl, nil
	}

	base, err := gh.AssetsFS.Open("/templates/base.html")
	if err != nil {
		return nil, err
	}

	defer base.Close()

	baseContents, err := io.ReadAll(base)
	if err != nil {
		return nil, err
	}

	return tpl.Parse(string(baseContents))
}
