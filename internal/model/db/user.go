package db

import "time"

const TableNameUser = "user"

// User 用户模型
type User struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"`           // 密码不返回给前端
	Salt      string    `json:"-" gorm:"size:32;not null;default:''"` // 密码盐值
	Role      int32     `json:"role" gorm:"not null;default:0"`
	Profile   string    `json:"profile" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*User) TableName() string {
	return TableNameUser
}
