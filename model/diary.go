package model

// Diary 工作经历
type Diary struct {
	ID
	Content string `gorm:"type:text" json:"content"`
	IsOpen  bool   `json:"is_open"`
	Timestamps
}

func (Diary) GetTableName() string {
	return "diaries"
}

type DiaryPs struct {
	ID
	Password string `json:"password"`
}

func (DiaryPs) GetTableName() string {
	return "diary_ps"
}
