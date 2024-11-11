package courses

import (
	"context"
	"fmt"
	inscriptionsDAO "inscriptions-api/dao/inscriptions"
	inscriptionsDomain "inscriptions-api/domain/inscriptions"
)

type Repository interface {
	EnrollUser(ctx context.Context, inscription inscriptionsDAO.Inscription) (string, error)
	ValidateEnrrol(ctx context.Context, inscription inscriptionsDAO.Inscription) error
}

type Service struct {
	mainRepository Repository
}

func NewService(mainRepository Repository) Service {
	return Service{
		mainRepository: mainRepository,
	}
}

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
