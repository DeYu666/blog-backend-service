package blog

import (
	"gorm.io/gorm"
)

type Option func(db *gorm.DB)

func setIdByUint(id uint) Option {
	return func(db *gorm.DB) {
		db.Where("`id` = ?", id)
	}
}
