package dao

import "time"

type User struct {
	UserID    uint   `gorm:"primarykey"`
	Email     string `gorm:"unique"`
	Password  string `gorm:"not null"`
	IsAdmin   bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
