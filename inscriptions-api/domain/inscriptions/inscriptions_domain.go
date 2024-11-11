package inscriptions

type Inscription struct {
	ID       string `json:"_id"`
	CourseID string `json:"course_id"`
	UserID   string `json:"user_id"`
}
