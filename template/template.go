// Copyright 2015 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package template

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type Template struct {
	tmpl *template.Template
}

// FromGlobs calls ParseGlob on all path globs provided and returns the
// resulting tmpl.
func FromGlobs(loadBuiltinTemplate bool, paths ...string) (*Template, error) {
	tmpl := template.New("").
		Option("missingkey=zero").
		Funcs(defaultFuncs).
		Funcs(sprig.TxtFuncMap())

	if loadBuiltinTemplate {
		f, err := Assets.Open("/templates/default.tmpl")
		if err != nil {
			return nil, err
		}

		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}

		if _, err := tmpl.Parse(string(b)); err != nil {
			return nil, err
		}
	}

	for _, tp := range paths {
		// ParseGlob in the template packages errors if not at least one file is
		// matched. We want to allow empty matches that may be populated later on.
		p, err := filepath.Glob(tp)
		if err != nil {
			return nil, err
		}
		if len(p) > 0 {
			if _, err := tmpl.ParseGlob(tp); err != nil {
				return nil, err
			}
		}
	}

	return &Template{tmpl: tmpl}, nil
}

func (t *Template) ExecuteTextString(text string, data interface{}) (string, error) {
	if text == "" {
		return "", nil
	}
	tmpl, err := t.tmpl.Clone()
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
