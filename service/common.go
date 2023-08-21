package service

import (
	"errors"
	"time"
)

var (
	ErrGenerateHasCate = errors.New("该总类下仍有分类存在，不能删除")
)

func generateUnionId() uint {
	return uint(time.Now().UnixNano())
}
