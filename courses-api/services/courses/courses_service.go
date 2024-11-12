package courses

import (
	"context"
	coursesDAO "courses-api/dao/courses"
	inscriptionsDAO "courses-api/dao/inscriptions"
	coursesDomain "courses-api/domain/courses"
	inscriptionsDomain "courses-api/domain/inscriptions"
	"fmt"
)

type Repository interface {
	GetCourses(ctx context.Context) (coursesDAO.Courses, error)
	GetCourseByID(ctx context.Context, id string) (coursesDAO.Course, error)
	CreateCourse(ctx context.Context, course coursesDAO.Course) (string, error)
	UpdateCourse(ctx context.Context, course coursesDAO.Course) error
	EnrollUser(ctx context.Context, inscription inscriptionsDAO.Inscription) (string, error)
	ValidateEnrrol(ctx context.Context, inscription inscriptionsDAO.Inscription) error
	GetInscriptionsByUserId(ctx context.Context, userID string) ([]inscriptionsDAO.Inscription, error)
}

type Queue interface {
	Publish(courseNew coursesDomain.CourseNew) error
}

type Service struct {
	mainRepository Repository
	eventsQueue    Queue
}

func NewService(mainRepository Repository, eventsQueue Queue) Service {
	return Service{
		mainRepository: mainRepository,
		eventsQueue:    eventsQueue,
	}
}

func (service Service) GetCourses(ctx context.Context) (coursesDomain.Courses, error) {
	results, err := service.mainRepository.GetCourses(ctx)
	if err != nil {
		return coursesDomain.Courses{}, fmt.Errorf("error getting courses from repository: %v", err)
	}

	// Convertimos la lista de cursos de DAO a DTO
	var courses coursesDomain.Courses
	for _, course := range results {
		courseDomain := coursesDomain.Course{
			ID:           course.ID,
			Name:         course.Name,
			Description:  course.Description,
			Professor:    course.Professor,
			ImageURL:     course.ImageURL,
			Requirement:  course.Requirement,
			Duration:     course.Duration,
			Availability: course.Availability,
		}
		courses = append(courses, courseDomain)
	}

	return courses, nil
}

func (service Service) GetCourseByID(ctx context.Context, id string) (coursesDomain.Course, error) {
	courseDAO, err := service.mainRepository.GetCourseByID(ctx, id)
	if err != nil {
		return coursesDomain.Course{}, fmt.Errorf("error getting course from repository: %v", err)
	}

	// Convert DAO to DTO
	return coursesDomain.Course{
		ID:           courseDAO.ID,
		Name:         courseDAO.Name,
		Description:  courseDAO.Description,
		Professor:    courseDAO.Professor,
		ImageURL:     courseDAO.ImageURL,
		Requirement:  courseDAO.Requirement,
		Duration:     courseDAO.Duration,
		Availability: courseDAO.Availability,
	}, nil
}

func (service Service) CreateCourse(ctx context.Context, course coursesDomain.Course) (string, error) {
	// Convert domain model to DAO model
	record := coursesDAO.Course{
		Name:         course.Name,
		Description:  course.Description,
		Professor:    course.Professor,
		ImageURL:     course.ImageURL,
		Requirement:  course.Requirement,
		Duration:     course.Duration,
		Availability: course.Availability,
	}
	id, err := service.mainRepository.CreateCourse(ctx, record)
	if err != nil {
		return "", fmt.Errorf("error creating course in main repository: %w", err)
	}
	// Set ID from main repository to use in the rest of the repositories
	record.ID = id

	if err := service.eventsQueue.Publish(coursesDomain.CourseNew{
		Operation: "CREATE",
		CourseID:  id,
	}); err != nil {
		return "", fmt.Errorf("error publishing course new: %w", err)
	}

	return id, nil
}

func (service Service) UpdateCourse(ctx context.Context, course coursesDomain.Course) error {
	// Convert domain model to DAO model
	record := coursesDAO.Course{
		ID:           course.ID,
		Name:         course.Name,
		Description:  course.Description,
		Professor:    course.Professor,
		ImageURL:     course.ImageURL,
		Requirement:  course.Requirement,
		Duration:     course.Duration,
		Availability: course.Availability,
	}

	// Update the hotel in the main repository
	err := service.mainRepository.UpdateCourse(ctx, record)
	if err != nil {
		return fmt.Errorf("error updating course in main repository: %w", err)
	}

	// Publish an event for the update operation
	if err := service.eventsQueue.Publish(coursesDomain.CourseNew{
		Operation: "UPDATE",
		CourseID:  course.ID,
	}); err != nil {
		return fmt.Errorf("error publishing course update: %w", err)
	}

	return nil
}

// ******************************** I N S C R I P T I O N S

func (service Service) EnrollUser(ctx context.Context, inscription inscriptionsDomain.Inscription) (string, error) {
	// Convert domain model to DAO model
	record := inscriptionsDAO.Inscription{
		CourseID: inscription.CourseID,
		UserID:   inscription.UserID,
	}

	err := service.mainRepository.ValidateEnrrol(ctx, record)
	if err == nil {
		return "", fmt.Errorf("error user is already enrolled in the course")
	}

	id, err := service.mainRepository.EnrollUser(ctx, record)
	if err != nil {
		return "", fmt.Errorf("error creating course in main repository: %w", err)
	}
	// Set ID from main repository to use in the rest of the repositories
	record.ID = id

	return id, nil
}

func (service Service) GetCoursesByUserID(ctx context.Context, userID string) ([]coursesDomain.Course, error) {
	// Obtener inscripciones del repositorio
	inscriptions, err := service.mainRepository.GetInscriptionsByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting inscriptions from repository: %v", err)
	}

	var coursesList []coursesDomain.Course
	for _, inscription := range inscriptions {
		// Obtener curso por ID usando el CourseID de cada inscripción
		courseDAO, err := service.mainRepository.GetCourseByID(ctx, inscription.CourseID)
		if err != nil {
			return nil, fmt.Errorf("error getting course from repository: %v", err)
		}

		// Mapear la información del curso desde courseDAO a coursesDomain.Course
		course := coursesDomain.Course{
			ID:           courseDAO.ID,
			Name:         courseDAO.Name,
			Description:  courseDAO.Description,
			Professor:    courseDAO.Professor,
			ImageURL:     courseDAO.ImageURL,
			Requirement:  courseDAO.Requirement,
			Duration:     courseDAO.Duration,
			Availability: courseDAO.Availability,
		}

		// Agregar el curso a la lista de cursos
		coursesList = append(coursesList, course)
	}

	return coursesList, nil
}
