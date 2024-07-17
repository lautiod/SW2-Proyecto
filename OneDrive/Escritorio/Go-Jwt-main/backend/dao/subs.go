package dao

import "time"

type Sub struct {
	SubID     uint `gorm:"primarykey"`
	UserID    uint `gorm:"not null"`
	CourseID  uint `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
