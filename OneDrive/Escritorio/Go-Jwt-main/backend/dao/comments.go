package dao

import "time"

type Comment struct {
	CommentID uint   `gorm:"primaryKey"`
	CourseID  uint   `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comments []Comment
