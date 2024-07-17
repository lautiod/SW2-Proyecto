package dao

import "time"

type File struct {
	FileID    uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FileName  string
	Data      string `gorm:"type:text;not null"`
	CourseID  uint   `gorm:"not null"`
}
