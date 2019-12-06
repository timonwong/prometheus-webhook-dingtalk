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
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/timonwong/prometheus-webhook-dingtalk/api"
	"github.com/timonwong/prometheus-webhook-dingtalk/config"
	"github.com/timonwong/prometheus-webhook-dingtalk/pkg/chilog"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
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

	var tmpl *template.Template
	api := &api.API{
		Logger: log.With(logger, "component", "api"),
	}

	configLogger := log.With(logger, "component", "configuration")
	configCoordinator := config.NewCoordinator(
		*configFile,
		configLogger,
	)
	configCoordinator.Subscribe(func(conf *config.Config) error {
		// Parse templates
		var err error
		level.Info(logger).Log("msg", "Loading templates", "templates", strings.Join(conf.Templates, ";"))
		tmpl, err = template.FromGlobs(conf.Templates...)
		if err != nil {
			return errors.Wrap(err, "failed to parse templates")
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

		api.Update(conf, tmpl)
		return nil
	})

	if err := configCoordinator.Reload(); err != nil {
		return 1
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&chilog.KitLogger{Logger: logger}))
	r.Use(middleware.Recoverer)
	r.Mount("/dingtalk", api.Routes())

	srv := http.Server{Addr: *listenAddress, Handler: r}
	srvCh := make(chan struct{})

	go func() {
		level.Info(logger).Log("msg", "Listening on address", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
			close(srvCh)
		}
	}()

	var (
		hup      = make(chan os.Signal, 1)
		hupReady = make(chan bool)
		term     = make(chan os.Signal, 1)
	)
	signal.Notify(hup, syscall.SIGHUP)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-hupReady
		for {
			select {
			case <-hup:
				// ignore error, already logged in `reload()`
				_ = configCoordinator.Reload()
			case <-term:
				return
			case <-srvCh:
				return
			}
		}
	}()

	// Wait for reload or termination signals.
	close(hupReady) // Unblock SIGHUP handler.

	select {
	case <-term:
		level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		_ = srv.Shutdown(ctx)
		return 0
	case <-srvCh:
		return 1
	}
}
