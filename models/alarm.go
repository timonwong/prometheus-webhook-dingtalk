package models

type Alarm struct {
	Dingding   string `gorm:"column:dingding" json:"dingding"`
	Attendance string `gorm:"column:attendance" json:"attendance"`
	Content    string `gorm:"column:content" json:"content"`
	Title      string `gorm:"column:title" json:"title"`
	Status     string `gorm:"column:status" json:"status"`
}

func (Alarm) TableName() string {
	return "main_alarm"
}
