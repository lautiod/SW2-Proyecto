package dao

import "time"

type Course struct {
	CourseID    uint   `gorm:"primarykey"`
	Name        string `gorm:"unique"`
	Description string `gorm:"not null"`
	Category    string `gorm:"not null"`
	ImageURL    string
	CreatorID   uint `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Courses []Course
