package webrouter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/timonwong/prometheus-webhook-dingtalk/models"
	"github.com/timonwong/prometheus-webhook-dingtalk/notifier"
)

type DingTalkResource struct {
	Profiles   map[string]string
	HttpClient *http.Client
}

func (rs *DingTalkResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/:profile/send", rs.SendNotification)
	return r
}

func (rs *DingTalkResource) SendNotification(w http.ResponseWriter, r *http.Request) {
	profile := chi.URLParam(r, "profile")
	webhookURL, ok := rs.Profiles[profile]
	if !ok || webhookURL == "" {
		http.NotFound(w, r)
		return
	}

	var promMessage models.WebhookMessage
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&promMessage); err != nil {
		log.Printf("Cannot decode prometheus webhook JSON request: %s", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	notification, err := notifier.BuildDingTalkNotification(&promMessage)
	if err != nil {
		log.Printf("Failed to build notification: %s", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	robotResp, err := notifier.SendDingTalkNotification(rs.HttpClient, webhookURL, notification)
	if err != nil {
		log.Printf("Failed to send notification: %s", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

	if robotResp.ErrorCode != 0 {
		log.Printf("Failed to send notification to DingTalk: [%d] %s", robotResp.ErrorCode, robotResp.ErrorMessage)
		return
	}

	log.Println("Successfully send notification to DingTalk")
}
