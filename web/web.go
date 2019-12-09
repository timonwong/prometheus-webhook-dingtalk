// Copyright 2013 The Prometheus Authors
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

package web

import (
	"context"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/server"

	"github.com/timonwong/prometheus-webhook-dingtalk/config"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
	"github.com/timonwong/prometheus-webhook-dingtalk/web/apiv1"
	"github.com/timonwong/prometheus-webhook-dingtalk/web/dingtalk"
	"github.com/timonwong/prometheus-webhook-dingtalk/web/ui"
)

var (
	// Paths that are handled by the React / Reach router that should all be served the main React app's index.html.
	reactRouterPaths = []string{
		"/",
		"/playground",
		"/config",
		"/flags",
		"/status",
	}
)

// Options for the web Handler.
type Options struct {
	ListenAddress string
	EnableWebUI   bool
	Version       *VersionInfo
	Flags         map[string]string
}

type VersionInfo = apiv1.VersionInfo

type Handler struct {
	mtx    sync.RWMutex
	logger log.Logger

	apiV1    *apiv1.API
	dingTalk *dingtalk.API

	router      chi.Router
	options     *Options
	config      *config.Config
	tmpl        *template.Template
	versionInfo *VersionInfo
	birth       time.Time
	cwd         string
}

func New(logger log.Logger, o *Options) *Handler {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "<error retrieving current working directory>"
	}

	h := &Handler{
		logger: logger,

		options:     o,
		versionInfo: o.Version,
		birth:       time.Now(),
		cwd:         cwd,
	}

	h.apiV1 = apiv1.NewAPI(
		logger,
		func() *config.Config {
			h.mtx.RLock()
			defer h.mtx.RUnlock()
			return h.config
		},
		func() *template.Template {
			h.mtx.RLock()
			defer h.mtx.RUnlock()
			return h.tmpl
		},
		o.Flags,
		h.versionInfo,
		h.runtimeInfo,
	)
	h.dingTalk = dingtalk.NewAPI(logger)

	router := chi.NewRouter()
	h.router = router

	fs := server.StaticFileServer(ui.Assets)
	router.Mount("/dingtalk", h.dingTalk.Routes())

	if o.EnableWebUI {
		router.Mount("/api/v1", h.apiV1.Routes())
		router.Get("/static/*", fs.ServeHTTP)
		// Make sure that "/ui" is redirected to "/ui/" and
		// not just the naked "/ui/"
		router.Get("/ui", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/ui/", http.StatusFound)
		})

		router.Get("/ui/*", func(w http.ResponseWriter, r *http.Request) {
			p := strings.TrimPrefix(r.URL.Path, "/ui")
			// For paths that the React/Reach router handles, we want to serve the
			// index.html, but with replaced path prefix placeholder.
			for _, rp := range reactRouterPaths {
				if p != rp {
					continue
				}

				f, err := ui.Assets.Open("/static/react/index.html")
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error opening React index.html: %v", err)
					return
				}
				idx, err := ioutil.ReadAll(f)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Error reading React index.html: %v", err)
					return
				}
				w.Write(idx) // nolint: errcheck
				return
			}

			// For all other paths, serve auxiliary assets.
			r.URL.Path = path.Join("/static/react/", p)
			fs := server.StaticFileServer(ui.Assets)
			fs.ServeHTTP(w, r)
		})
	}

	return h
}

// ApplyConfig updates the config field of the Handler struct
func (h *Handler) ApplyConfig(conf *config.Config, tmpl *template.Template) error {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	h.config = conf
	h.tmpl = tmpl
	h.dingTalk.Update(conf, tmpl)
	return nil
}

// Run serves the HTTP endpoints.
func (h *Handler) Run(ctx context.Context) error {
	level.Info(h.logger).Log("msg", "Start listening for connections", "address", h.options.ListenAddress)
	listener, err := net.Listen("tcp", h.options.ListenAddress)
	if err != nil {
		return err
	}

	errlog := stdlog.New(log.NewStdlibAdapter(level.Error(h.logger)), "", 0)
	httpSrv := &http.Server{
		Handler:  h.router,
		ErrorLog: errlog,
	}

	errCh := make(chan error)
	go func() {
		errCh <- httpSrv.Serve(listener)
	}()

	select {
	case e := <-errCh:
		return e
	case <-ctx.Done():
		return nil
	}
}

func (h *Handler) runtimeInfo() (*apiv1.RuntimeInfo, error) {
	status := &apiv1.RuntimeInfo{
		StartTime:  h.birth,
		CWD:        h.cwd,
		GOMAXPROCS: runtime.GOMAXPROCS(0),
		GOGC:       os.Getenv("GOGC"),
		GODEBUG:    os.Getenv("GODEBUG"),
	}
	return status, nil
}
