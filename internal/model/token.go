package model

import "time"

// APIToken 用户生成的 API Token 模型
type APIToken struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64      `json:"user_id" gorm:"index;not null"`
	Name      string     `json:"name" gorm:"size:100;not null"`
	TokenHash string     `json:"-" gorm:"size:64;uniqueIndex;not null"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
