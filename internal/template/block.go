package template

import (
	"bufio"
	"bytes"
	"text/template"
)

// ParseBlock parses a block of text using a template.
func ParseBlock(block string, api map[string]interface{}, data interface{}) (string, error) {
	helperAPI := defaultAPI
	for k, v := range api {
		helperAPI[k] = v
	}

	tpl, err := template.New("custom").Funcs(helperAPI).Parse(block)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	if err := tpl.Execute(w, data); err != nil {
		return "", err
	}

	if err := w.Flush(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
