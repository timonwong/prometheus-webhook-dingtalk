package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/timonwong/prometheus-webhook-dingtalk/webrouter"
)

func main() {
	if err := parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return
		}
		log.Fatalf("Parse error: %s", err)
	}

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// When a client closes their connection midway through a request, the
	// http.CloseNotifier will cancel the request context (ctx).
	r.Use(middleware.CloseNotify)

	dingTalkResource := &webrouter.DingTalkResource{
		Profiles: cfg.dingTalkProfiles.profiles,
		HttpClient: &http.Client{
			Timeout: cfg.requestTimeout,
			Transport: &http.Transport{
				Proxy:             http.ProxyFromEnvironment,
				DisableKeepAlives: true,
			},
		},
	}
	r.Mount("/dingtalk", dingTalkResource.Routes())

	log.Printf("Starting webserver on %s", cfg.listenAddress)
	if err := http.ListenAndServe(cfg.listenAddress, r); err != nil {
		log.Panicf("Failed to serve: %s", err)
	}
}
