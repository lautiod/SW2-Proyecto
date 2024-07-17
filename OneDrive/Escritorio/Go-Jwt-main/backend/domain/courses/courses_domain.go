package courses

type CourseDetail struct {
	CourseID    uint   `json:"courseID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
	CreatorID   uint   `json:"creatorId"`
}

type CoursesDetail []CourseDetail

type FileRequest struct {
	Data     string `json:"data"`
	FileName string `json:"file_name"`
	CourseID uint   `json:"course_id"`
}
