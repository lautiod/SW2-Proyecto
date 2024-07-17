package clients

import (
	"backend/dao"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectToDb() {
	var err error
	dsn := os.Getenv("Db")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to DB")
	}
}

func SyncDb() {
	Db.AutoMigrate(&dao.User{}, &dao.Course{}, &dao.Sub{}, &dao.Comment{}, &dao.File{})
}

func SelectUserByEmail(email string) (dao.User, error) {

	var user dao.User

	Db.First(&user, "email = ?", email)

	if user.UserID == 0 {
		return dao.User{}, errors.New("failed to query")
	}
	return user, nil
}

func SelectUserByID(id uint) (dao.User, error) {
	var user dao.User
	result := Db.First(&user, id)
	if result.Error != nil {
		return dao.User{}, errors.New("user not found")
	}
	return user, nil
}

func SelectCourseByID(id uint) (dao.Course, error) {
	var course dao.Course
	result := Db.First(&course, id)
	if result.Error != nil {
		return dao.Course{}, errors.New("course not found")
	}
	return course, nil
}

func Subscribe(user uint, course uint) error {

	sub := dao.Sub{UserID: user, CourseID: course}

	result := Db.Create(&sub)
	if result.Error != nil {
		return errors.New("failed to create sub")
	}
	return nil
}

func ValidateSub(user uint, course uint) error {

	var sub dao.Sub
	err := Db.Where("user_id = ? AND course_id = ?", user, course).First(&sub).Error
	if err == nil {
		return errors.New("user is already subscribed")
	}
	return nil
}

func GetCourseIDsByUserID(userID uint) ([]uint, error) {
	var subs []dao.Sub

	// Encontrar todas las subscripciones del usuario
	if err := Db.Where("user_id = ?", userID).Find(&subs).Error; err != nil {
		return nil, errors.New("failed to get courses")
	}

	// Obtener los IDs de los cursos a los que el usuario está inscrito
	courseIDs := make([]uint, len(subs))
	for i, sub := range subs {
		courseIDs[i] = sub.CourseID
	}

	return courseIDs, nil
}

func GetUserIDsByCourseID(courseID uint) ([]uint, error) {
	var subs []dao.Sub

	// Encontrar todas las subscripciones del usuario
	if err := Db.Where("course_id = ?", courseID).Find(&subs).Error; err != nil {
		return nil, errors.New("failed to get users")
	}

	// Obtener los IDs de los cursos a los que el usuario está inscrito
	usersIDs := make([]uint, len(subs))
	for i, sub := range subs {
		usersIDs[i] = sub.UserID
	}

	return usersIDs, nil
}

func GetCourses() dao.Courses {
	var courses dao.Courses
	Db.Find(&courses)

	log.Debug("Courses: ", courses)

	return courses
}

func GetCourseById(id int) dao.Course {
	var course dao.Course

	Db.Where("course_id = ?", id).First(&course)
	log.Debug("Course: ", course)

	return course
}

func GetCoursesBySearch(query string) (dao.Courses, error) {
	var courses dao.Courses

	err := Db.Where("name LIKE ? OR category LIKE ?", "%"+query+"%", "%"+query+"%").Find(&courses).Error
	if err != nil {
		return nil, errors.New("failed to get courses  by search")
	}

	return courses, nil
}

func GetCoursesByUserID(userID uint) (dao.Courses, error) {
	var courses []dao.Course
	err := Db.Joins("JOIN subs ON subs.course_id = courses.course_id").
		Where("subs.user_id = ?", userID).
		Find(&courses).Error
	if err != nil {
		return nil, errors.New("failed to get courses  by search")
	}

	log.Debug("Courses for user ID ", userID, ": ", courses)
	return courses, nil
}

func AddComment(userID uint, courseID uint, content string) error {
	comment := dao.Comment{
		CourseID: courseID,
		UserID:   userID,
		Content:  content,
	}

	if err := Db.Create(&comment).Error; err != nil {
		log.Error("Error adding comment: ", err)
		return err
	}

	log.Debug("Added comment: ", comment)
	return nil
}

func GetCommentsByCourseID(courseID uint) (dao.Comments, error) {
	var comments []dao.Comment
	err := Db.Where("course_id = ?", courseID).Find(&comments).Error
	if err != nil {
		log.Error("Error retrieving comments for course ID ", courseID, ": ", err)
		return nil, err
	}

	log.Debug("Comments for course ID ", courseID, ": ", comments)
	return comments, nil
}

func CreateCourse(name string, description string, category string, imageUrl string, creatorId uint) error {
	course := dao.Course{
		Name:        name,
		Description: description,
		Category:    category,
		ImageURL:    imageUrl,
		CreatorID:   creatorId,
	}

	if err := Db.Create(&course).Error; err != nil {
		log.Error("Error adding course: ", err)
		return err
	}

	log.Debug("Added course: ", course)
	return nil

}

func AddFile(filename string, courseId uint, data string) error {

	file := dao.File{
		FileName: filename,
		CourseID: courseId,
		Data:     data,
	}

	if err := Db.Create(&file).Error; err != nil {
		log.Error("Error adding file: ", err)
		return err
	}

	log.Debug("Added file: ", file)
	return nil

}
