package courses

import (
	"backend/clients"
	_ "backend/dao"
	domain "backend/domain/courses"
	"errors"
)

func CheckUsers(userid uint, courseid uint) error {
	course, err := clients.SelectCourseByID(courseid)
	if err != nil {
		return errors.New("course not found")
	}

	if course.CreatorID != userid {
		return errors.New("not allowed")
	}

	return nil
}

func GetCourses() (domain.CoursesDetail, error) {

	courses := clients.GetCourses()
	var coursesDomain domain.CoursesDetail

	for _, course := range courses {
		courseDomain := domain.CourseDetail{
			CourseID:    course.CourseID,
			Name:        course.Name,
			Description: course.Description,
			Category:    course.Category,
			ImageURL:    course.ImageURL,
			CreatorID:   course.CreatorID,
		}

		coursesDomain = append(coursesDomain, courseDomain)
	}

	return coursesDomain, nil
}

func GetCourseById(id int) (domain.CourseDetail, error) {
	course := clients.GetCourseById(id)
	var courseDomain domain.CourseDetail

	if course.CourseID == 0 {
		return courseDomain, errors.New("course not found")
	}

	courseDomain = domain.CourseDetail{
		CourseID:    course.CourseID,
		Name:        course.Name,
		Description: course.Description,
		Category:    course.Category,
		ImageURL:    course.ImageURL,
		CreatorID:   course.CreatorID,
	}

	return courseDomain, nil
}

func GetCoursesBySearch(query string) (domain.CoursesDetail, error) {
	courses, err := clients.GetCoursesBySearch(query)
	if err != nil {
		return nil, err
	}

	var coursesDomain domain.CoursesDetail

	for _, course := range courses {
		courseDomain := domain.CourseDetail{
			CourseID:    course.CourseID,
			Name:        course.Name,
			Description: course.Description,
			Category:    course.Category,
			ImageURL:    course.ImageURL,
			CreatorID:   course.CreatorID,
		}

		coursesDomain = append(coursesDomain, courseDomain)
	}

	return coursesDomain, nil
}

func CreateCourse(course domain.CourseDetail, creatorID uint) error {
	err := clients.CreateCourse(course.Name, course.Description, course.Category, course.ImageURL, creatorID)

	if err != nil {
		return err
	}

	return nil
}

func AddFile(filereq domain.FileRequest) error {
	err := clients.AddFile(filereq.FileName, filereq.CourseID, filereq.Data)

	if err != nil {
		return err
	}
	return nil
}
