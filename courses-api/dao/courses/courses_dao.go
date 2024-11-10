package courses

type Course struct {
	ID           string  `bson:"_id,omitempty"`
	Name         string  `bson:"name"`
	Description  string  `bson:"description"`
	Professor    string  `bson:"professor"`
	ImageURL     string  `bson:"image_url"`
	Duration     float64 `bson:"duration"`
	Requirement  string  `bson:"requirement"`
	Availability float64 `bson:"availability"`
}

type Courses []Course
