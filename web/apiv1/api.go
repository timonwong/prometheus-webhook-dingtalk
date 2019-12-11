// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiv1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/timonwong/prometheus-webhook-dingtalk/config"
	"github.com/timonwong/prometheus-webhook-dingtalk/notifier"
	"github.com/timonwong/prometheus-webhook-dingtalk/pkg/models"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
)

type API struct {
	logger      log.Logger
	config      func() *config.Config
	tmpl        func() *template.Template
	flagsMap    map[string]string
	versionInfo *VersionInfo
	runtimeInfo func() (*RuntimeInfo, error)
}

func NewAPI(logger log.Logger,
	config func() *config.Config,
	tmpl func() *template.Template,
	flagsMap map[string]string,
	versionInfo *VersionInfo,
	runtimeInfo func() (*RuntimeInfo, error)) *API {

	return &API{
		logger:      logger,
		config:      config,
		tmpl:        tmpl,
		flagsMap:    flagsMap,
		versionInfo: versionInfo,
		runtimeInfo: runtimeInfo,
	}
}

func (api *API) Routes() chi.Router {
	wrap := func(f apiFunc) http.HandlerFunc {
		hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := f(r)
			if result.err != nil {
				api.respondError(w, result.err, result.data)
			} else if result.data != nil {
				api.respond(w, result.data)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		})
		return hf
	}

	router := chi.NewRouter()
	router.Get("/status/templates", wrap(api.serveTemplates))
	router.Post("/status/templates/render", wrap(api.serveRenderTemplate))
	router.Get("/status/config", wrap(api.serveConfig))
	router.Get("/status/runtimeinfo", wrap(api.serveRuntimeInfo))
	router.Get("/status/buildinfo", wrap(api.serveBuildInfo))
	router.Get("/status/flags", wrap(api.serveFlags))
	return router
}

type status string

const (
	statusSuccess status = "success"
	statusError   status = "error"
)

type errorType string

const (
	errorTimeout  errorType = "timeout"
	errorCanceled errorType = "canceled"
	errorExec     errorType = "execution"
	errorBadData  errorType = "bad_data"
	errorInternal errorType = "internal"
	errorNotFound errorType = "not_found"
)

type apiError struct {
	typ errorType
	err error
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s: %s", e.typ, e.err)
}

type response struct {
	Status    status      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	ErrorType errorType   `json:"errorType,omitempty"`
	Error     string      `json:"error,omitempty"`
}

type apiFuncResult struct {
	data interface{}
	err  *apiError
}

type apiFunc func(r *http.Request) apiFuncResult

func (api *API) serveTemplates(r *http.Request) apiFuncResult {
	type templateInfo struct {
		Target string `json:"name"`
		Title  string `json:"title"`
		Text   string `json:"text"`
	}

	type templatesInfo struct {
		Templates []templateInfo `json:"templates"`
	}

	conf := api.config()
	defaultMessage := conf.GetDefaultMessage()
	templates := []templateInfo{
		{
			Target: "<default>",
			Title:  defaultMessage.Title,
			Text:   defaultMessage.Text,
		},
	}
	for name, target := range conf.Targets {
		if target.Message == nil {
			templates = append(templates, templateInfo{
				Target: name,
				Title:  defaultMessage.Title,
				Text:   defaultMessage.Text,
			})
		} else {
			templates = append(templates, templateInfo{
				Target: name,
				Title:  target.Message.Title,
				Text:   target.Message.Text,
			})
		}
	}

	info := &templatesInfo{Templates: templates}
	return apiFuncResult{info, nil}
}

func (api *API) serveRenderTemplate(r *http.Request) apiFuncResult {
	var req struct {
		Title         string `json:"title"`
		Text          string `json:"text"`
		DemoAlertJSON string `json:"demoAlertJSON"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return apiFuncResult{nil, &apiError{errorBadData, err}}
	}

	var webhookMessage models.WebhookMessage
	err = json.Unmarshal([]byte(req.DemoAlertJSON), &webhookMessage)
	if err != nil {
		return apiFuncResult{nil, &apiError{errorBadData, err}}
	}

	// Construct a fake "target"
	target := &config.Target{
		Message: &config.TargetMessage{
			Title: req.Title,
			Text:  req.Text,
		},
	}
	builder := notifier.NewDingNotificationBuilder(api.tmpl(), api.config(), target)
	notification, err := builder.Build(&webhookMessage)
	if err != nil {
		return apiFuncResult{nil, &apiError{errorBadData, err}}
	}

	resp := struct {
		Markdown string `json:"markdown"`
	}{
		Markdown: notification.Markdown.Text,
	}
	return apiFuncResult{&resp, nil}
}

func (api *API) serveConfig(r *http.Request) apiFuncResult {
	cfg := &struct {
		YAML string `json:"yaml"`
	}{
		YAML: api.config().String(),
	}
	return apiFuncResult{cfg, nil}
}

type RuntimeInfo struct {
	StartTime time.Time `json:"startTime"`
	CWD       string    `json:"CWD"`
	//ReloadConfigSuccess bool      `json:"reloadConfigSuccess"`
	//LastConfigTime      time.Time `json:"lastConfigTime"`
	GoroutineCount int    `json:"goroutineCount"`
	GOMAXPROCS     int    `json:"GOMAXPROCS"`
	GOGC           string `json:"GOGC"`
	GODEBUG        string `json:"GODEBUG"`
}

func (api *API) serveRuntimeInfo(r *http.Request) apiFuncResult {
	status, err := api.runtimeInfo()
	if err != nil {
		return apiFuncResult{status, &apiError{errorInternal, err}}
	}

	return apiFuncResult{status, nil}
}

type VersionInfo struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	Branch    string `json:"branch"`
	BuildUser string `json:"buildUser"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

func (api *API) serveBuildInfo(r *http.Request) apiFuncResult {
	return apiFuncResult{api.versionInfo, nil}
}

func (api *API) serveFlags(r *http.Request) apiFuncResult {
	return apiFuncResult{api.flagsMap, nil}
}

func (api *API) respond(w http.ResponseWriter, data interface{}) {
	statusMessage := statusSuccess
	b, err := json.Marshal(&response{
		Status: statusMessage,
		Data:   data,
	})
	if err != nil {
		level.Error(api.logger).Log("msg", "error marshaling json response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if n, err := w.Write(b); err != nil {
		level.Error(api.logger).Log("msg", "error writing response", "bytesWritten", n, "err", err)
	}
}

func (api *API) respondError(w http.ResponseWriter, apiErr *apiError, data interface{}) {
	b, err := json.Marshal(&response{
		Status:    statusError,
		ErrorType: apiErr.typ,
		Error:     apiErr.err.Error(),
		Data:      data,
	})

	if err != nil {
		level.Error(api.logger).Log("msg", "error marshaling json response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var code int
	switch apiErr.typ {
	case errorBadData:
		code = http.StatusBadRequest
	case errorExec:
		code = 422
	case errorCanceled, errorTimeout:
		code = http.StatusServiceUnavailable
	case errorInternal:
		code = http.StatusInternalServerError
	case errorNotFound:
		code = http.StatusNotFound
	default:
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if n, err := w.Write(b); err != nil {
		level.Error(api.logger).Log("msg", "error writing response", "bytesWritten", n, "err", err)
	}
}
