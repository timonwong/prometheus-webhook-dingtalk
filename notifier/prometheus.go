package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/timonwong/prometheus-webhook-dingtalk/models"
	"github.com/timonwong/prometheus-webhook-dingtalk/template"

	mysql "github.com/timonwong/prometheus-webhook-dingtalk/middleware"
)

func BuildDingTalkNotification(dingding string, promMessage *models.WebhookMessage) (*models.DingTalkNotification, error) {
	title, err := template.ExecuteTextString(`{{ template "ding.link.title" . }}`, promMessage)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	content, err := template.ExecuteTextString(`{{ template "ding.link.content" . }}`, promMessage)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	buf, _ := json.Marshal(*promMessage.Data)
	fmt.Println(string(buf))

	var buttons []models.DingTalkNotificationButton
	for i, alert := range promMessage.Alerts.Firing() {
		buttons = append(buttons, models.DingTalkNotificationButton{
			Title:     fmt.Sprintf("Graph for alert #%d", i+1),
			ActionURL: alert.GeneratorURL,
		})
	}

	notification := &models.DingTalkNotification{
		MessageType: "text",
		Text: &models.DingTalkNotificationText{
			Title:   title,
			Content: content,
		},
	}

	fmt.Println(promMessage.Status, promMessage.CommonLabels, "!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	notification.At = new(models.DingTalkNotificationAt)
	alarm := models.Alarm{Dingding: dingding, Title: title, Content: content, Status: promMessage.Status}
	users := []string{}
	if v, ok := map[string]string(promMessage.CommonLabels)["at"]; ok {

		uri := fmt.Sprintf("%v?name=%v&on_duty_date=%v", MonitorCoreAddress, v, fmt.Sprintf("%v", time.Now().Day()))
		fmt.Println("request info:", uri)
		response, _ := httpClient(uri, "GET", nil, "") //请求monitor-core获取电话号码

		var resp Response
		if err := json.Unmarshal(response, &resp); err != nil {
			fmt.Println("request err:", err)
		} else {
			fmt.Println("response data:", resp.Data)
		}

		head := map[string]string{"Servicetoken": LinkedseeToken, "Content-Type": "application/json"}
		for _, recver := range resp.Data {
			//短信告警
			notification.At.AtMobiles = append(notification.At.AtMobiles, recver.Phone)
			users = append(users, recver.Name)

			if _, err := httpClient(LinkedseeUrl, "GET", head, SendAlarm{
				Receiver: recver.Phone,
				Type:     "sms",
				Title:    "alarm_sms",
				Content:  content,
			}); err != nil {
				fmt.Println(err.Error())
			}

			//告警回执，不打电话
			if strings.ToUpper(promMessage.Status) == "RESOLVED" {
				continue
			}

			//电话告警
			if _, err := httpClient(LinkedseeUrl, "GET", head, SendAlarm{
				Receiver: recver.Phone,
				Type:     "phone",
				Title:    "alarm_phone",
				Content:  string(content),
			}); err != nil {
				fmt.Println(err.Error())
			}

		}
	}
	alarm.Attendance = strings.Join(users, ",")
	if mysql.SAVE_TO_MYSQL {
		if err := mysql.GormDB.Create(alarm).Error; err != nil {
			fmt.Println(err.Error())
		}
	}

	return notification, nil
}

func SendDingTalkNotification(httpClient *http.Client, webhookURL string, notification *models.DingTalkNotification) (*models.DingTalkNotificationResponse, error) {
	body, err := json.Marshal(&notification)
	if err != nil {
		return nil, errors.Wrap(err, "error encoding DingTalk request")
	}

	httpReq, err := http.NewRequest("POST", webhookURL, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "error building DingTalk request")
	}
	httpReq.Header.Set("Content-Type", "application/json")

	req, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "error sending notification to DingTalk")
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		return nil, errors.Errorf("unacceptable response code %d", req.StatusCode)
	}

	var robotResp models.DingTalkNotificationResponse
	enc := json.NewDecoder(req.Body)
	if err := enc.Decode(&robotResp); err != nil {
		return nil, errors.Wrap(err, "error decoding response from DingTalk")
	}

	return &robotResp, nil
}
