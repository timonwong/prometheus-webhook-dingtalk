package api

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/timonwong/prometheus-webhook-dingtalk/config"
	"github.com/timonwong/prometheus-webhook-dingtalk/pkg/models"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
)

type API struct {
	// Protect against config, template and http client
	mtx sync.RWMutex

	tmpl       *template.Template
	targets    map[string]config.Target
	httpClient *http.Client

	Logger log.Logger
}

func (api *API) Update(conf *config.Config, tmpl *template.Template) {
	api.mtx.Lock()
	defer api.mtx.Unlock()

	api.targets = conf.Targets
	api.httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}
	api.tmpl = tmpl
}

func (api *API) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/{name}/send", api.SendHandler)
	return r
}

func (api *API) SendHandler(w http.ResponseWriter, r *http.Request) {
	api.mtx.RLock()
	targets := api.targets
	tmpl := api.tmpl
	httpClient := api.httpClient
	api.mtx.RUnlock()

	targetName := chi.URLParam(r, "name")
	logger := log.With(api.Logger, "target", targetName)

	target, ok := targets[targetName]
	if !ok {
		level.Warn(logger).Log("msg", "target not found")
		http.NotFound(w, r)
		return
	}

	var promMessage models.WebhookMessage
	if err := json.NewDecoder(r.Body).Decode(&promMessage); err != nil {
		level.Error(logger).Log("msg", "Cannot decode prometheus webhook JSON request", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	notification, err := buildDingTalkNotification(tmpl, &target, &promMessage)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to build notification", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	robotResp, err := sendDingTalkNotification(httpClient, &target, notification)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to send notification", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if robotResp.ErrorCode != 0 {
		level.Error(logger).Log("msg", "Failed to send notification to DingTalk", "respCode", robotResp.ErrorCode, "respMsg", robotResp.ErrorMessage)
		http.Error(w, "Unable to talk to DingTalk", http.StatusBadRequest)
		return
	}

	io.WriteString(w, "OK") // nolint: errcheck
}
