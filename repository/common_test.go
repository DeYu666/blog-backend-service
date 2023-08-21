package repository

import (
	"context"

	"gorm.io/gorm"
)

type TxFunc func(ctx context.Context, tx *gorm.DB) (interface{}, error)
