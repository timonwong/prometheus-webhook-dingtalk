package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/timonwong/prometheus-webhook-dingtalk/config"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
	"github.com/timonwong/prometheus-webhook-dingtalk/web"
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
		enableWebUI = kingpin.Flag(
			"web.ui-enabled",
			"Enable Web UI mounted on /ui path",
		).Default("false").Bool()
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

	flagsMap := map[string]string{}
	// Exclude kingpin default flags to expose only Prometheus ones.
	boilerplateFlags := kingpin.New("", "").Version("")
	for _, f := range kingpin.CommandLine.Model().Flags {
		if boilerplateFlags.GetFlag(f.Name) != nil {
			continue
		}

		flagsMap[f.Name] = f.Value.String()
	}

	webHandler := web.New(log.With(logger, "component", "web"), &web.Options{
		ListenAddress: *listenAddress,
		EnableWebUI:   *enableWebUI,
		Version: &web.VersionInfo{
			Version:   version.Version,
			Revision:  version.Revision,
			Branch:    version.Branch,
			BuildUser: version.BuildUser,
			BuildDate: version.BuildDate,
			GoVersion: version.GoVersion,
		},
		Flags: flagsMap,
	})

	configLogger := log.With(logger, "component", "configuration")
	configCoordinator := config.NewCoordinator(
		*configFile,
		configLogger,
	)
	configCoordinator.Subscribe(func(conf *config.Config) error {
		// Parse templates
		level.Info(configLogger).Log("msg", "Loading templates", "templates", strings.Join(conf.Templates, ";"))
		tmpl, err := template.FromGlobs(conf.Templates...)
		if err != nil {
			return errors.Wrap(err, "failed to parse templates")
		}

		// Print current targets configuration
		host, port, _ := net.SplitHostPort(*listenAddress)
		if host == "" {
			host = "localhost"
		}

		var paths []string
		for name := range conf.Targets {
			paths = append(paths, fmt.Sprintf("http://%s:%s/dingtalk/%s/send", host, port, name))
		}
		configLogger.Log("msg", "Webhook urls for prometheus alertmanager", "urls", strings.Join(paths, " "))

		return webHandler.ApplyConfig(conf, tmpl)
	})

	if err := configCoordinator.Reload(); err != nil {
		return 1
	}

	ctxWeb, cancelWeb := context.WithCancel(context.Background())
	defer cancelWeb()

	srvCh := make(chan error, 1)
	go func() {
		defer close(srvCh)

		if err := webHandler.Run(ctxWeb); err != nil {
			level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
			srvCh <- err
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
			case <-ctxWeb.Done():
				return
			case <-hup:
				// ignore error, already logged in `reload()`
				_ = configCoordinator.Reload()
			}
		}
	}()

	// Wait for reload or termination signals.
	close(hupReady) // Unblock SIGHUP handler.

	for {
		select {
		case <-term:
			level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
			cancelWeb()
		case err := <-srvCh:
			if err != nil {
				return 1
			}

			return 0
		}
	}
}
