package comments

import (
	client "backend/clients"
	domain "backend/domain/comments"
)

func AddComment(userID uint, request domain.CommentRequest) error {
	err := client.AddComment(userID, request.CourseID, request.Content)

	if err != nil {
		return err
	}

	return nil
}

func GetCommentsByCourseID(courseID uint) (domain.CommentsDetail, error) {
	comments, err := client.GetCommentsByCourseID(courseID)
	if err != nil {
		return nil, err
	}

	var commentsDomain domain.CommentsDetail

	for _, comment := range comments {
		user, err := client.SelectUserByID(comment.UserID)
		if err != nil {
			return nil, err
		}
		commentDomain := domain.CommentDetail{
			CommentID: comment.CommentID,
			CourseID:  comment.CourseID,
			UserID:    comment.UserID,
			Email:     user.Email,
			Content:   comment.Content,
		}

		commentsDomain = append(commentsDomain, commentDomain)
	}

	return commentsDomain, nil
}
