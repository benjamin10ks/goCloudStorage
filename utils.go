package main

import (
	"fmt"
	"text/template"
)

func loadTemplates() (map[string]*template.Template, error) {
	tmpls := make(map[string]*template.Template)
	pages := []string{"login", "dashboard", "passkey_begin"}

	for _, page := range pages {
		tmpl, err := template.ParseFiles(
			"web/templates/layout.tmpl",
			fmt.Sprintf("web/templates/%s.tmpl", page),
		)
		if err != nil {
			return nil, err
		}
		tmpls[page] = tmpl
	}
	return tmpls, nil
}
