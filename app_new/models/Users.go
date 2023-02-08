package models

import "strconv"

// AuthUser 用户表
type AuthUser struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Password    string `json:"password,omitempty"`
	LastLogin   string `json:"last_login" json:"-"`
	IsSuperuser string `json:"-"`
	Username    string `json:"username"`
	FirstName   string `json:"-"`
	Email       string `json:"email,omitempty"`
	IsStaff     string `json:"-"`
	IsActive    string `json:"-"`
	DateJoined  string `json:"dataJoined,omitempty"`
	LastName    string `json:"-"`
}

func (AuthUser) GetTableName() string {
	return "auth_users"
}

func (user AuthUser) GetUid() string {
	return strconv.Itoa(user.ID)
}
