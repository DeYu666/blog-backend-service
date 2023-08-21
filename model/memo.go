package model

import "time"

// Memo 备忘录
type Memo struct {
	ID
	Content string `gorm:"type:text" json:"content"`
	Status  int    `json:"status"` // 0 表示未完成； 1 表示已完成；-1 表示全部
	// Status  int    `json:"status"` // 1 表示未完成； 2 表示已完成；-1 表示全部
	Timestamps
}

func (Memo) GetTableName() string {
	return "memos"
}

// ShowMemo 往外展示的备忘录结构
type ShowMemo struct {
	CreateTime time.Time `json:"create_time"`
	Content    []Memo    `json:"content"`
}
