package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/timonwong/prometheus-webhook-dingtalk/chilog"
	"github.com/timonwong/prometheus-webhook-dingtalk/config"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
	"github.com/timonwong/prometheus-webhook-dingtalk/webrouter"
)

func main() {
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

	cfg, err := config.LoadFile(*configFile)
	if err != nil {
		level.Error(logger).Log("msg", "Error reading configuration file", "err", err)
		os.Exit(1)
	}

	// Load & validate customized template file
	if cfg.Template != "" {
		l := log.With(logger, "filename", cfg.Template)

		b, err := ioutil.ReadFile(cfg.Template)
		if err != nil {
			level.Error(l).Log("msg", "Error reading customizable template file", "err", err)
			os.Exit(1)
		}

		_, err = template.UpdateTemplate(string(b))
		if err != nil {
			level.Error(l).Log("msg", "Error parsing template file", "err", err)
			os.Exit(1)
		}

		level.Info(l).Log("msg", "Using customized template")
	} else {
		level.Info(logger).Log("msg", "Using default template")
	}

	// Print current targets configuration
	if l := level.Info(logger); l != nil {
		host, port, _ := net.SplitHostPort(*listenAddress)
		if host == "" {
			host = "localhost"
		}

		var paths []string
		for name := range cfg.Targets {
			paths = append(paths, fmt.Sprintf("http://%s:%s/dingtalk/%s/send", host, port, name))
		}
		l.Log("msg", "Webhook urls for prometheus alertmanager", "urls", strings.Join(paths, " "))
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&chilog.KitLogger{Logger: logger}))
	r.Use(middleware.Recoverer)

	dingTalkResource := &webrouter.DingTalkResource{
		Logger:  logger,
		Targets: cfg.Targets,
		HttpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				Proxy:             http.ProxyFromEnvironment,
				DisableKeepAlives: true,
			},
		},
	}
	r.Mount("/dingtalk", dingTalkResource.Routes())

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, r); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
