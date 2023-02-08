package models

import "time"

type LoveInfo struct {
	ID
	KnownTime      time.Time `gorm:"type:time" json:"known_time"`
	ConfessionTime time.Time `gorm:"type:time" json:"confession_time"`
	LoveName       string    `json:"love_name"`
	ExtraInfo      string    `gorm:"type:text" json:"extra_info"`
}

func (LoveInfo) GetTableName() string {
	return "love_infos"
}
