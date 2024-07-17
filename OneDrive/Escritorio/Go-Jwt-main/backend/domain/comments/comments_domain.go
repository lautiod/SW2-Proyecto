package comments

type CommentRequest struct {
	CourseID uint   `json:"course_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type CommentDetail struct {
	CommentID uint   `json:"comment_id"`
	CourseID  uint   `json:"course_id"`
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Content   string `json:"content"`
}

type CommentsDetail []CommentDetail
