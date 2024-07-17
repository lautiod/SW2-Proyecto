package users

import (
	"backend/clients"
	"backend/dao"
	domainCourses "backend/domain/courses"
	domain "backend/domain/users"
	"errors"

	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(email string, password string, admin bool) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return errors.New("failed to hash password")
	}

	user := dao.User{Email: email, Password: string(hash), IsAdmin: admin}

	result := clients.Db.Create(&user)

	if result.Error != nil {
		return errors.New("failed to create user")
	}
	return nil
}

func Login(body domain.Login_Request) (domain.LoginResponse, string, error) {

	user, err := clients.SelectUserByEmail(body.Email)
	if err != nil {
		return domain.LoginResponse{}, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		return domain.LoginResponse{}, "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UserID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return domain.LoginResponse{}, "", errors.New("failed to create token")
	}

	userDomain := domain.LoginResponse{
		IsAdmin: user.IsAdmin,
		Email:   user.Email,
	}

	return userDomain, tokenString, nil
}

func GetCoursesByUserID(userID uint) (domainCourses.CoursesDetail, error) {
	courses, err := clients.GetCoursesByUserID(userID)
	if err != nil {
		return nil, err
	}

	var coursesDomain domainCourses.CoursesDetail

	for _, course := range courses {
		courseDomain := domainCourses.CourseDetail{
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

/*
func Search(query string) ([]domain.Course, error) {
	trimmed := strings.TrimSpace(query)

	courses, err := clients.SelectCoursesWithFilter(trimmed)
	if err != nil {
		return nil, fmt.Errorf("error getting courses from DB: %w", err)
	}

	results := make([]domain.Course, 0)
	for _, course := range courses {
		results = append(results, domain.Course{
			ID:           course.ID,
			Title:        course.Title,
			Description:  course.Description,
			Category:     course.Category,
			ImageURL:     course.ImageURL,
			CreationDate: course.CreationDate,
			LastUpdated:  course.LastUpdated,
		})
	}
	return results, nil
}
*/
