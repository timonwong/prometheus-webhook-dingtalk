package notifier

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

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
	m.AlertTime = time.Now().Format("2006.01.02 15:04:05")
	title, err := r.renderTitle(m)
	if err != nil {
		return nil, err
	}
	content, err := r.renderText(m)
	if err != nil {
		return nil, err
	}

	notification := &models.DingTalkNotification{
		MessageType: "text",
		Text: &models.DingTalkNotificationText{
			Title:   title,
			Content: content,
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
		return nil, errors.Wrap(err, "error encoding DingTalk request")
	}

	httpReq, err := http.NewRequest("POST", targetURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "error building DingTalk request")
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
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
