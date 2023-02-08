package models

// Memo 备忘录
type Memo struct {
	ID
	Content string `gorm:"type:text" json:"content"`
	Status  int    `json:"status"` // 0 表示未完成； 1 表示已完成；2 表示全部
	Timestamps
}

func (Memo) GetTableName() string {
	return "memos"
}
