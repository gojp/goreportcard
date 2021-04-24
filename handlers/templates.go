package handlers

import (
	"fmt"
	"io/ioutil"
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

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	footer, err := gh.AssetsFS.Open("/templates/footer.html")
	if err != nil {
		return nil, err
	}

	defer footer.Close()

	footerContents, err := ioutil.ReadAll(footer)
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

	return tpl.Parse(string(footerContents))
}
