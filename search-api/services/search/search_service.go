package search

import (
	"context"
	"fmt"
	coursesDAO "search-api/dao/courses"
	coursesDomain "search-api/domain/courses"
)

type Repository interface {
	Index(ctx context.Context, course coursesDAO.Course) (string, error)
	Update(ctx context.Context, course coursesDAO.Course) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string, limit int, offset int) ([]coursesDAO.Course, error) // Updated signature
}

type ExternalRepository interface {
	GetCourseByID(ctx context.Context, id string) (coursesDomain.Course, error)
}

type Service struct {
	repository Repository
	coursesAPI ExternalRepository
}

func NewService(repository Repository, coursesAPI ExternalRepository) Service {
	return Service{
		repository: repository,
		coursesAPI: coursesAPI,
	}
}

func (service Service) Search(ctx context.Context, query string, offset int, limit int) ([]coursesDomain.Course, error) {
	// Call the repository's Search method
	coursesDAOList, err := service.repository.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error searching courses: %w", err)
	}

	// Convert the dao layer courses to domain layer courses
	coursesDomainList := make([]coursesDomain.Course, 0)
	for _, course := range coursesDAOList {
		coursesDomainList = append(coursesDomainList, coursesDomain.Course{
			ID:           course.ID,
			Name:         course.Name,
			Description:  course.Description,
			Professor:    course.Professor,
			ImageURL:     course.ImageURL,
			Requirement:  course.Requirement,
			Duration:     course.Duration,
			Availability: course.Availability,
		})
	}

	return coursesDomainList, nil
}

func (service Service) HandleCourseNew(courseNew coursesDomain.CourseNew) {
	switch courseNew.Operation {
	case "CREATE", "UPDATE":
		// Fetch course details from the local service
		course, err := service.coursesAPI.GetCourseByID(context.Background(), courseNew.CourseID)
		if err != nil {
			fmt.Printf("Error getting course (%s) from API: %v\n", courseNew.CourseID, err)
			return
		}

		courseDAO := coursesDAO.Course{
			ID:           course.ID,
			Name:         course.Name,
			Description:  course.Description,
			Professor:    course.Professor,
			ImageURL:     course.ImageURL,
			Requirement:  course.Requirement,
			Duration:     course.Duration,
			Availability: course.Availability,
		}

		// Handle Index operation
		if courseNew.Operation == "CREATE" {
			if _, err := service.repository.Index(context.Background(), courseDAO); err != nil {
				fmt.Printf("Error indexing course (%s): %v\n", courseNew.CourseID, err)
			} else {
				fmt.Println("Course indexed successfully:", courseNew.CourseID)
			}
		} else { // Handle Update operation
			if err := service.repository.Update(context.Background(), courseDAO); err != nil {
				fmt.Printf("Error updating course (%s): %v\n", courseNew.CourseID, err)
			} else {
				fmt.Println("Course updated successfully:", courseNew.CourseID)
			}
		}

	default:
		fmt.Printf("Unknown operation: %s\n", courseNew.Operation)
	}
}