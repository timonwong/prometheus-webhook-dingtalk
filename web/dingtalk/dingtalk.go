package dingtalk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/VictoriaMetrics/metrics"
	"github.com/timonwong/prometheus-webhook-dingtalk/config"
	"github.com/timonwong/prometheus-webhook-dingtalk/notifier"
	"github.com/timonwong/prometheus-webhook-dingtalk/pkg/chilog"
	"github.com/timonwong/prometheus-webhook-dingtalk/pkg/models"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
)

type API struct {
	// Protect against config, template and http client
	mtx sync.RWMutex

	conf       *config.Config
	tmpl       *template.Template
	targets    map[string]config.Target
	httpClient *http.Client
	logger     log.Logger

	MaxAlertCount uint16 // to avoid too many alert
}

func NewAPI(logger log.Logger) *API {
	return &API{
		logger: logger,
	}
}

func (api *API) Update(conf *config.Config, tmpl *template.Template) {
	api.mtx.Lock()
	defer api.mtx.Unlock()

	api.conf = conf
	api.tmpl = tmpl
	api.targets = conf.Targets
	api.httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}
}

func (api *API) Routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestLogger(&chilog.KitLogger{Logger: api.logger}))
	router.Use(middleware.Recoverer)
	router.Post("/{name}/send", api.serveSend)
	return router
}

func (api *API) serveSend(w http.ResponseWriter, r *http.Request) {
	metrics.GetOrCreateCounter("http_requests_total{path=\"" + r.RequestURI + "\"}").Inc()
	api.mtx.RLock()
	targets := api.targets
	conf := api.conf
	tmpl := api.tmpl
	httpClient := api.httpClient
	api.mtx.RUnlock()

	targetName := chi.URLParam(r, "name")
	logger := log.With(api.logger, "target", targetName)

	target, ok := targets[targetName]
	if !ok {
		const reason = "target not found"
		level.Warn(logger).Log("msg", reason)
		http.NotFound(w, r)
		metrics.GetOrCreateCounter(
			fmt.Sprintf(`http_requests_error{path="%s",reason="%s"}`, r.RequestURI, reason)).Inc()
		return
	}

	var promMessage models.WebhookMessage
	if err := json.NewDecoder(r.Body).Decode(&promMessage); err != nil {
		const reason = "Cannot decode prometheus webhook JSON request"
		level.Error(logger).Log("msg", reason, "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		metrics.GetOrCreateCounter(
			fmt.Sprintf(`http_requests_error{path="%s",reason="%s"}`, r.RequestURI, reason)).Inc()
		return
	}

	metrics.GetOrCreateCounter(
		fmt.Sprintf(`alert_count{path="%s"}`, r.RequestURI)).Add(len(promMessage.Alerts))
	if api.MaxAlertCount > 0 && len(promMessage.Alerts) > int(api.MaxAlertCount) {
		metrics.GetOrCreateCounter(
			fmt.Sprintf(`alert_drop_count{path="%s"}`, r.RequestURI)).Add(len(promMessage.Alerts) - int(api.MaxAlertCount))
		promMessage.Alerts = promMessage.Alerts[:api.MaxAlertCount]
	}

	builder := notifier.NewDingNotificationBuilder(tmpl, conf, &target)
	notification, err := builder.Build(&promMessage)
	if err != nil {
		const reason = "Failed to build notification"
		level.Error(logger).Log("msg", reason, "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		metrics.GetOrCreateCounter(
			fmt.Sprintf(`http_requests_error{path="%s",reason="%s"}`, r.RequestURI, reason)).Inc()
		return
	}

	robotResp, err := notifier.SendNotification(notification, httpClient, &target)
	if err != nil {
		const reason = "Failed to send notification"
		level.Error(logger).Log("msg", reason, "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		metrics.GetOrCreateCounter(
			fmt.Sprintf(`http_requests_error{path="%s",reason="%s"}`, r.RequestURI, reason)).Inc()
		return
	}

	if robotResp.ErrorCode != 0 {
		const reason = "Failed to send notification to DingTalk"
		level.Error(logger).Log("msg", reason, "respCode", robotResp.ErrorCode, "respMsg", robotResp.ErrorMessage)
		http.Error(w, "Unable to talk to DingTalk", http.StatusBadRequest)
		metrics.GetOrCreateCounter(
			fmt.Sprintf(`http_requests_error{path="%s",reason="%s",respCode="%d",respMsg="%s"}`,
				r.RequestURI, reason, robotResp.ErrorCode, robotResp.ErrorMessage)).Inc()
		return
	}

	io.WriteString(w, "OK")
}
