package tpl

import (
	"bytes"
	"strings"
	"text/template"
)

const (
	alertTemplateText = `
{{ define "__subject" }}[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] {{ .GroupLabels.SortedPairs.Values | join " " }} {{ if gt (len .CommonLabels) (len .GroupLabels) }}({{ with .CommonLabels.Remove .GroupLabels.Names }}{{ .Values | join " " }}{{ end }}){{ end }}{{ end }}
{{ define "__alertmanagerURL" }}{{ .ExternalURL }}/#/alerts?receiver={{ .Receiver }}{{ end }}

{{ define "__text_alert_list" }}{{ range . }}
**Labels**
{{ range .Labels.SortedPairs }}> - {{ .Name }}: {{ .Value | markdown | html }}
{{ end }}
**Annotations**
{{ range .Annotations.SortedPairs }}> - {{ .Name }}: {{ .Value | markdown | html }}
{{ end }}
**Source:** {{ .GeneratorURL }}

{{ end }}{{ end }}

{{ define "ding.link.title" }}{{ template "__subject" . }}{{ end }}
{{ define "ding.link.content" }}#### \[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}\] **[{{ index .GroupLabels "alertname" }}]({{ template "__alertmanagerURL" . }})**
{{ template "__text_alert_list" .Alerts.Firing }}
{{ end }}
`
)

var (
	alertTemplate = template.Must(template.New("alert").Funcs(defaultFuncs).Option("missingkey=zero").Parse(alertTemplateText))
	defaultFuncs  = map[string]interface{}{
		"toUpper": strings.ToUpper,
		"toLower": strings.ToLower,
		"title":   strings.Title,
		// join is equal to strings.Join but inverts the argument order
		// for easier pipelining in templates.
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
		"markdown": markdownEscapeString,
	}
	isMarkdownSpecial [128]bool
)

func init() {
	for _, c := range "_*`" {
		isMarkdownSpecial[c] = true
	}
}

func markdownEscapeString(s string) string {
	b := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(b)

	for _, c := range s {
		if c < 128 && isMarkdownSpecial[c] {
			buf.WriteByte('\\')
		}
		buf.WriteRune(c)
	}
	return buf.String()
}

func ExecuteTextString(text string, data interface{}) (string, error) {
	if text == "" {
		return "", nil
	}
	tmpl, err := alertTemplate.Clone()
	if err != nil {
		return "", err
	}
	tmpl, err = tmpl.New("").Option("missingkey=zero").Parse(text)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
