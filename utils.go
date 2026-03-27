package main

import (
	"fmt"
	"text/template"
)

func loadTemplates() (map[string]*template.Template, error) {
	tmpls := make(map[string]*template.Template)
	
	// Base layouts
	baseFiles := []string{
		"web/templates/base.tmpl",
		"web/templates/auth_layout.tmpl",
		"web/templates/app_layout.tmpl",
	}
	
	// Component files
	componentFiles := []string{
		"web/templates/components/file_card.tmpl",
		"web/templates/components/modal_share.tmpl",
		"web/templates/components/toast.tmpl",
		"web/templates/components/passkey_login.tmpl",
	}
	
	// Auth pages
	authPages := map[string][]string{
		"login":    append(baseFiles, "web/templates/login.tmpl"),
		"register": append(baseFiles, "web/templates/register.tmpl"),
		"passkey_begin": append(baseFiles, "web/templates/passkey_begin.tmpl"),
	}
	
	// App pages
	appPages := map[string][]string{
		"dashboard": append(append(baseFiles, componentFiles...), "web/templates/dashboard.tmpl"),
	}
	
	// Parse auth pages
	for name, files := range authPages {
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %w", name, err)
		}
		tmpls[name] = tmpl
	}
	
	// Parse app pages
	for name, files := range appPages {
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %w", name, err)
		}
		tmpls[name] = tmpl
	}
	
	return tmpls, nil
}
