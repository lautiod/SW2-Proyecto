package subscriptions

import (
	"backend/clients"
	domain "backend/domain/subscriptions"
)

func Subscribe(userId uint, request domain.SubRequest) error {
	if _, err := clients.SelectUserByID(userId); err != nil {
		return err
	}

	if _, err := clients.SelectCourseByID(request.CourseID); err != nil {
		return err
	}

	err := clients.ValidateSub(userId, request.CourseID)
	if err != nil {
		return err
	}

	err = clients.Subscribe(userId, request.CourseID)

	if err != nil {
		return err
	}
	return nil
}
