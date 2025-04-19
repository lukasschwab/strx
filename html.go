package main

import (
	_ "embed"

	"fmt"
	"html/template"
	"io"
)

//go:embed templates/html.tmpl
var htmlTemplate string

type Row struct {
	Alias string
	URL   string
}

func HTML(all map[string]string, to io.Writer) error {
	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("error parsing HTML template: %w", err)
	}

	rows := make([]Row, 0, len(all))
	for alias, url := range all {
		rows = append(rows, Row{
			Alias: alias,
			URL:   url,
		})
	}

	err = tmpl.Execute(to, rows)
	if err != nil {
		return fmt.Errorf("error executing HTML template: %w", err)
	}
	return nil
}
