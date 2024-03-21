package templates

import (
	"bufio"
	"bytes"
	"text/template"
)

func ParseBlock(block string, api map[string]interface{}, data interface{}) (string, error) {
	helperApi := defaultApi
	for k, v := range api {
		helperApi[k] = v
	}

	tpl, err := template.New("custom").Funcs(helperApi).Parse(block)
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
