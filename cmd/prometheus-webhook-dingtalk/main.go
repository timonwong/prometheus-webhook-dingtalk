package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/timonwong/prometheus-webhook-dingtalk/chilog"
	"github.com/timonwong/prometheus-webhook-dingtalk/webrouter"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	listenAddress    = kingpin.Flag("web.listen-address", "The address to listen on for web interface.").Default(":8060").String()
	dingTalkProfiles = DingTalkProfiles(kingpin.Flag("ding.profile", "Custom DingTalk profile (can specify multiple times, <profile>=<dingtalk-url>).").Required())
	requestTimeout   = kingpin.Flag("ding.timeout", "Timeout for invoking DingTalk webhook.").Default("5s").Duration()
)

func main() {
	allowedLevel := promlog.AllowedLevel{}
	flag.AddFlags(kingpin.CommandLine, &allowedLevel)
	kingpin.Version(version.Print("prometheus-webhook-dingtalk"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(allowedLevel)
	level.Info(logger).Log("msg", "Starting prometheus-webhook-dingtalk", "version", version.Info())

	r := chi.NewRouter()
	r.Use(middleware.CloseNotify)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&chilog.KitLogger{Logger: logger}))
	r.Use(middleware.Recoverer)

	dingTalkResource := &webrouter.DingTalkResource{
		Logger:   logger,
		Profiles: map[string]string(*dingTalkProfiles),
		HttpClient: &http.Client{
			Timeout: *requestTimeout,
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
