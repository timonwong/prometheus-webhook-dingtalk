package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/timonwong/prometheus-webhook-dingtalk/config"
	"github.com/timonwong/prometheus-webhook-dingtalk/pkg/chilog"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
	"github.com/timonwong/prometheus-webhook-dingtalk/webrouter"
)

func main() {
	os.Exit(run())
}

func run() int {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"The address to listen on for web interface.",
		).Default(":8060").String()
		configFile = kingpin.Flag(
			"config.file",
			"Path to the configuration file.",
		).Default("config.yml").ExistingFile()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("prometheus-webhook-dingtalk"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting prometheus-webhook-dingtalk", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", version.BuildContext())

	conf, err := config.LoadFile(*configFile)
	if err != nil {
		level.Error(logger).Log("msg", "Error reading configuration file", "err", err)
		return 1
	}

	// Parse templates
	level.Info(logger).Log("msg", "loading templates", "templates", strings.Join(conf.Templates, ";"))
	tmpl, err := template.FromGlobs(conf.Templates...)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to parse templates", "err", err)
		return 1
	}

	// Print current targets configuration
	if l := level.Info(logger); l != nil {
		host, port, _ := net.SplitHostPort(*listenAddress)
		if host == "" {
			host = "localhost"
		}

		var paths []string
		for name := range conf.Targets {
			paths = append(paths, fmt.Sprintf("http://%s:%s/dingtalk/%s/send", host, port, name))
		}
		l.Log("msg", "Webhook urls for prometheus alertmanager", "urls", strings.Join(paths, " "))
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&chilog.KitLogger{Logger: logger}))
	r.Use(middleware.Recoverer)

	dingTalkResource := &webrouter.DingTalkResource{
		Logger:   logger,
		Template: tmpl,
		Targets:  conf.Targets,
		HttpClient: &http.Client{
			Timeout: conf.Timeout,
			Transport: &http.Transport{
				Proxy:             http.ProxyFromEnvironment,
				DisableKeepAlives: true,
			},
		},
	}
	r.Mount("/dingtalk", dingTalkResource.Routes())

	srv := http.Server{Addr: *listenAddress, Handler: r}
	srvCh := make(chan struct{})

	go func() {
		level.Info(logger).Log("msg", "Listening on address", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
			close(srvCh)
		}
	}()

	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-termCh:
		level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		_ = srv.Shutdown(ctx)
		return 0
	case <-srvCh:
		return 1
	}
}
