package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	Model
	Username     string `gorm:"comment:登录名;type:varchar(128)" json:"username"`
	Salt         string `json:"-"`
	Password     string `gorm:"type:varchar(128)" json:"-"`
	Level        uint   `json:"level"`
	Unauthorized bool   `json:"unauthorized"`
	Name         string `gorm:"comment:用户姓名" json:"name"`
	Company      string `gorm:"comment:公司" json:"company"`
	Position     string `gorm:"comment:职位" json:"position"`
}
