package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/timonwong/prometheus-webhook-dingtalk/chilog"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
	"github.com/timonwong/prometheus-webhook-dingtalk/webrouter"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	listenAddress    = kingpin.Flag("web.listen-address", "The address to listen on for web interface.").Default(":8060").String()
	dingTalkProfiles = DingTalkProfiles(kingpin.Flag("ding.profile", "Custom DingTalk profile (can be given multiple times, <profile>=<dingtalk-url>).").Required())
	requestTimeout   = kingpin.Flag("ding.timeout", "Timeout for invoking DingTalk webhook.").Default("5s").Duration()
	templateFileName = kingpin.Flag("template.file", "Customized template file (see template/default.tmpl for example)").Default("").String()
)

func main() {
	allowedLevel := promlog.AllowedLevel{}
	flag.AddFlags(kingpin.CommandLine, &allowedLevel)
	kingpin.Version(version.Print("prometheus-webhook-dingtalk"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(allowedLevel)
	level.Info(logger).Log("msg", "Starting prometheus-webhook-dingtalk", "version", version.Info())

	// Load & validate customized template file
	if *templateFileName != "" {
		l := log.With(logger, "filename", *templateFileName)

		b, err := ioutil.ReadFile(*templateFileName)
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

	// Print current profile configuration
	profiles := map[string]string(*dingTalkProfiles)
	level.Info(logger).Log("msg", fmt.Sprintf("Using following dingtalk profiles: %v", profiles))

	r := chi.NewRouter()
	r.Use(middleware.CloseNotify)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&chilog.KitLogger{Logger: logger}))
	r.Use(middleware.Recoverer)

	dingTalkResource := &webrouter.DingTalkResource{
		Logger:   logger,
		Profiles: profiles,
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
