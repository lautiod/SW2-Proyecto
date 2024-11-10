package courses

type Course struct {
	ID           string  `json:"_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Professor    string  `json:"professor"`
	ImageURL     string  `json:"image_url"`
	Duration     float64 `json:"duration"`
	Requirement  string  `json:"requirement"`
	Availability float64 `json:"availability"`
}

type CourseNew struct {
	Operation string `json:"operation"`
	CourseID  string `json:"course_id"`
}
