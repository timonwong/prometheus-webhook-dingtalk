package notifier

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/timonwong/prometheus-webhook-dingtalk/config"
	"github.com/timonwong/prometheus-webhook-dingtalk/models"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
)

func BuildDingTalkNotification(promMessage *models.WebhookMessage) (*models.DingTalkNotification, error) {
	promMessage.AlertTime = time.Now().Format("2006.01.02 15:04:05")
	title, err := template.ExecuteTextString(`{{ template "ding.link.title" . }}`, promMessage)
	if err != nil {
		return nil, err
	}
	content, err := template.ExecuteTextString(`{{ template "ding.link.content" . }}`, promMessage)
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
	return notification, nil
}

func SendDingTalkNotification(httpClient *http.Client, target config.Target, notification *models.DingTalkNotification) (*models.DingTalkNotificationResponse, error) {
	body, err := json.Marshal(&notification)
	if err != nil {
		return nil, errors.Wrap(err, "error encoding DingTalk request")
	}

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

	httpReq, err := http.NewRequest("POST", targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "error building DingTalk request")
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "error sending notification to DingTalk")
	}
	defer resp.Body.Close()

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
