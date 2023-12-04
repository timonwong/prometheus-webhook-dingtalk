package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/VictoriaMetrics/metrics"
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
			"web.enable-ui",
			"Enable Web UI mounted on /ui path",
		).Default("false").Bool()
		enableLifecycle = kingpin.Flag(
			"web.enable-lifecycle",
			"Enable reload via HTTP request.",
		).Default("false").Bool()
		configFile = kingpin.Flag(
			"config.file",
			"Path to the configuration file.",
		).Default("config.yml").ExistingFile()
		// for push metrics
		extraLabel = kingpin.Flag(
			"pushmetrics.extraLabel",
			"extraLabel for push metrics.",
		).Default("").String()
		intervalForPushMetrics = kingpin.Flag(
			"pushmetrics.interval",
			"interval for push metrics.",
		).Default("15s").Duration()
		urlForPushMetrics = kingpin.Flag(
			"pushmetrics.url",
			"urls for push metrics.",
		).Default("").String()
		// aviod too many alert
		maxAlertCount = kingpin.Flag(
			"maxalertcount",
			"max alert count to send to ding talk.",
		).Default("30").Uint16()
	)

	// DO NOT REMOVE. For compatibility purpose
	kingpin.Flag("web.ui-enabled", "").Hidden().BoolVar(enableWebUI)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)

	kingpin.Version(version.Print("prometheus-webhook-dingtalk"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting prometheus-webhook-dingtalk", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", version.BuildContext())

	if len(*urlForPushMetrics) > 0 {
		metrics.InitPush(*urlForPushMetrics,
			*intervalForPushMetrics, *extraLabel, true)
	}

	flagsMap := map[string]string{}
	// Exclude kingpin default flags to expose only Prometheus ones.
	boilerplateFlags := kingpin.New("", "").Version("")
	for _, f := range kingpin.CommandLine.Model().Flags {
		if boilerplateFlags.GetFlag(f.Name) != nil {
			continue
		}

		// filter hidden flags (they are just reserved for compatibility purpose)
		if f.Hidden {
			continue
		}

		flagsMap[f.Name] = f.Value.String()
	}

	webHandler := web.New(log.With(logger, "component", "web"), &web.Options{
		ListenAddress:   *listenAddress,
		EnableWebUI:     *enableWebUI,
		EnableLifecycle: *enableLifecycle,
		Version: &web.VersionInfo{
			Version:   version.Version,
			Revision:  version.Revision,
			Branch:    version.Branch,
			BuildUser: version.BuildUser,
			BuildDate: version.BuildDate,
			GoVersion: version.GoVersion,
		},
		Flags:         flagsMap,
		MaxAlertCount: *maxAlertCount,
	})

	configLogger := log.With(logger, "component", "configuration")
	configCoordinator := config.NewCoordinator(*configFile, configLogger)
	configCoordinator.Subscribe(func(conf *config.Config) error {
		// Parse templates
		level.Info(configLogger).Log("msg", "Loading templates", "templates", strings.Join(conf.Templates, ";"))
		tmpl, err := template.FromGlobs(!conf.NoBuiltinTemplate, conf.Templates...)
		if err != nil {
			return fmt.Errorf("failed to parse templates: %w", err)
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
		reloadReady = make(chan struct{})
		hup         = make(chan os.Signal, 1)
		term        = make(chan os.Signal, 1)
	)
	signal.Notify(hup, syscall.SIGHUP)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-reloadReady
		for {
			select {
			case <-ctxWeb.Done():
				return
			case <-hup:
				// ignore error, already logged in `reload()`
				_ = configCoordinator.Reload()
			case rc := <-webHandler.Reload():
				if err := configCoordinator.Reload(); err != nil {
					rc <- err
				} else {
					rc <- nil
				}
			}
		}
	}()

	// Wait for reload or termination signals.
	close(reloadReady) // Unblock SIGHUP handler.
	webHandler.Ready()

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
