package models

type DingTalkNotificationResponse struct {
	ErrorMessage string `json:"errmsg"`
	ErrorCode    int    `json:"errcode"`
}

type DingTalkNotification struct {
	MessageType string                          `json:"msgtype"`
	Text        *DingTalkNotificationText       `json:"text,omitempty"`
	Link        *DingTalkNotificationLink       `json:"link,omitempty"`
	Markdown    *DingTalkNotificationMarkdown   `json:"markdown,omitempty"`
	ActionCard  *DingTalkNotificationActionCard `json:"actionCard,omitempty"`
	At          *DingTalkNotificationAt         `json:"at,omitempty"`
}

type DingTalkNotificationText struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DingTalkNotificationLink struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	MessageURL string `json:"messageUrl"`
	PictureURL string `json:"picUrl"`
}

type DingTalkNotificationMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingTalkNotificationAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

type DingTalkNotificationActionCard struct {
	Title             string                       `json:"title"`
	Text              string                       `json:"text"`
	HideAvatar        string                       `json:"hideAvatar"`
	ButtonOrientation string                       `json:"btnOrientation"`
	Buttons           []DingTalkNotificationButton `json:"btns,omitempty"`
	SingleTitle       string                       `json:"singleTitle,omitempty"`
	SingleURL         string                       `json:"singleURL"`
}

type DingTalkNotificationButton struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}
