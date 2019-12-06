package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

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

	r.Post("/{name}/send", api.SendNotification)
	return r
}

func (api *API) SendNotification(w http.ResponseWriter, r *http.Request) {
	logger := api.Logger
	targets := api.getTargets()
	targetName := chi.URLParam(r, "name")
	target, ok := targets[targetName]
	if !ok {
		level.Warn(logger).Log("msg", fmt.Sprintf("target %s not found", targetName))
		http.NotFound(w, r)
		return
	}

	if target.URL == "" {
		level.Warn(logger).Log("msg", fmt.Sprintf("target %s url is empty", targetName))
		http.NotFound(w, r)
		return
	}

	var promMessage models.WebhookMessage
	if err := json.NewDecoder(r.Body).Decode(&promMessage); err != nil {
		level.Error(logger).Log("msg", "Cannot decode prometheus webhook JSON request", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	notification, err := api.buildDingTalkNotification(&target, &promMessage)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to build notification", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	robotResp, err := api.sendDingTalkNotification(&target, notification)
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

func (api *API) buildDingTalkNotification(target *config.Target, m *models.WebhookMessage) (*models.DingTalkNotification, error) {
	tmpl := api.getTemplate()

	title, err := tmpl.ExecuteTextString(target.Message.Title, m)
	if err != nil {
		return nil, err
	}
	content, err := tmpl.ExecuteTextString(target.Message.Text, m)
	if err != nil {
		return nil, err
	}

	notification := &models.DingTalkNotification{
		MessageType: "markdown",
		Markdown: &models.DingTalkNotificationMarkdown{
			Title: title,
			Text:  content,
		},
	}

	// Build mention
	if target.Mention != nil {
		notification.At = &models.DingTalkNotificationAt{
			IsAtAll:   target.Mention.All,
			AtMobiles: target.Mention.Mobiles,
		}
	}

	return notification, nil
}

func (api *API) sendDingTalkNotification(target *config.Target, notification *models.DingTalkNotification) (*models.DingTalkNotificationResponse, error) {
	targetURL := target.URL
	// Calculate signature when secret is provided
	if target.Secret != "" {
		timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		stringToSign := []byte(timestamp + "\n" + target.Secret)

		mac := hmac.New(sha256.New, []byte(target.Secret))
		mac.Write(stringToSign) // nolint: errcheck
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		u, err := url.Parse(targetURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse target url")
		}

		qs := u.Query()
		qs.Set("timestamp", timestamp)
		qs.Set("sign", signature)
		u.RawQuery = qs.Encode()

		targetURL = u.String()
	}

	body, err := json.Marshal(&notification)
	if err != nil {
		return nil, errors.Wrap(err, "error encoding DingTalk request")
	}

	httpReq, err := http.NewRequest("POST", targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "error building DingTalk request")
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := api.getHttpClient().Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "error sending notification to DingTalk")
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("unacceptable response code %d", resp.StatusCode)
	}

	var robotResp models.DingTalkNotificationResponse
	enc := json.NewDecoder(resp.Body)
	if err := enc.Decode(&robotResp); err != nil {
		return nil, errors.Wrap(err, "error decoding response from DingTalk")
	}

	return &robotResp, nil
}

func (api *API) getTemplate() *template.Template {
	api.mtx.RLock()
	defer api.mtx.RUnlock()

	return api.tmpl
}

func (api *API) getTargets() map[string]config.Target {
	api.mtx.RLock()
	defer api.mtx.RUnlock()

	return api.targets
}

func (api *API) getHttpClient() *http.Client {
	api.mtx.RLock()
	defer api.mtx.RUnlock()

	return api.httpClient
}
