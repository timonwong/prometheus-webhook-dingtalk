package notifier

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/timonwong/prometheus-webhook-dingtalk/config"
	"github.com/timonwong/prometheus-webhook-dingtalk/pkg/models"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"
)

type DingNotificationBuilder struct {
	tmpl     *template.Template
	target   *config.Target
	titleTpl string
	textTpl  string
}

func NewDingNotificationBuilder(tmpl *template.Template, conf *config.Config, target *config.Target) *DingNotificationBuilder {
	// Message template from the following order:
	//   target level > config global level > builtin global level

	var (
		defaultMessage = conf.GetDefaultMessage()
		titleTpl       = defaultMessage.Title
		textTpl        = defaultMessage.Text
	)

	if target.Message != nil {
		titleTpl = target.Message.Title
		textTpl = target.Message.Text
	}

	return &DingNotificationBuilder{
		tmpl:     tmpl,
		target:   target,
		titleTpl: titleTpl,
		textTpl:  textTpl,
	}
}

func (r *DingNotificationBuilder) renderTitle(data interface{}) (string, error) {
	return r.tmpl.ExecuteTextString(r.titleTpl, data)
}

func (r *DingNotificationBuilder) renderText(data interface{}) (string, error) {
	return r.tmpl.ExecuteTextString(r.textTpl, data)
}

func (r *DingNotificationBuilder) Build(m *models.WebhookMessage) (*models.DingTalkNotification, error) {
	if r.target.Mention != nil {
		m.AtMobiles = append(m.AtMobiles, r.target.Mention.Mobiles...)
	}

	title, err := r.renderTitle(m)
	if err != nil {
		return nil, err
	}
	content, err := r.renderText(m)
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
	if r.target.Mention != nil {
		notification.At = &models.DingTalkNotificationAt{
			IsAtAll:   r.target.Mention.All,
			AtMobiles: r.target.Mention.Mobiles,
		}
	}

	return notification, nil
}

func SendNotification(notification *models.DingTalkNotification, httpClient *http.Client, target *config.Target) (*models.DingTalkNotificationResponse, error) {
	targetURL := *target.URL
	// Calculate signature when secret is provided
	if target.Secret != "" {
		timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		stringToSign := []byte(timestamp + "\n" + string(target.Secret))

		mac := hmac.New(sha256.New, []byte(target.Secret))
		mac.Write(stringToSign) // nolint: errcheck
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		qs := targetURL.Query()
		qs.Set("timestamp", timestamp)
		qs.Set("sign", signature)
		targetURL.RawQuery = qs.Encode()
	}

	body, err := json.Marshal(&notification)
	if err != nil {
		return nil, fmt.Errorf("error encoding DingTalk request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", targetURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error building DingTalk request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error sending notification to DingTalk: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unacceptable response code %d", resp.StatusCode)
	}

	var robotResp models.DingTalkNotificationResponse
	enc := json.NewDecoder(resp.Body)
	if err := enc.Decode(&robotResp); err != nil {
		return nil, fmt.Errorf("error decoding response from DingTalk: %w", err)
	}

	return &robotResp, nil
}
