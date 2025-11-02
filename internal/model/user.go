package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`           // 密码不返回给前端
	Salt      string         `json:"-" gorm:"size:32;not null;default:''"` // 密码盐值
	Role      int32          `json:"role" gorm:"not null;default:0"`
	Profile   string         `json:"profile" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // 软删除
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     int32  `json:"role" binding:"required"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
