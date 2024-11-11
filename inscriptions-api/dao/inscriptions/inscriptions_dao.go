package inscriptions

type Inscription struct {
	ID       string `bson:"_id,omitempty"`
	CourseID string `bson:"course_id"`
	UserID   string `bson:"user_id"`
}
