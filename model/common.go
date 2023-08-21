package model

import (
	"gorm.io/gorm"
	"time"
)

// ID 自增ID主键
type ID struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`
}

// Timestamps 创建、更新时间
type Timestamps struct {
	CreatedTime  time.Time       `gorm:"type:time" json:"create_time"`
	ModifiedTime time.Time       `gorm:"type:time" json:"modified_time"`
}

// SoftDeletes 软删除
type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

//GRANT ALL PRIVILEGES ON *.* TO 'deyu666'@'%' IDENTIFIED BY 'WoshiRR.1314' WITH GRANT OPTION;
