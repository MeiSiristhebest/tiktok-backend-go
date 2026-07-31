package service

import (
	"errors"
	"time"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/repository"
)

type CommentService struct{}

var DefaultCommentService = &CommentService{}

func (s *CommentService) CommentAction(userID, videoID int64, actionType int32, commentText string, commentID int64) (*model.CommentDTO, error) {
	if actionType == 1 {
		// 1 - 发表评论
		if len(commentText) == 0 {
			return nil, errors.New("comment content cannot be empty")
		}

		comment := model.CommentEntity{
			UserID:    userID,
			VideoID:   videoID,
			Content:   commentText,
			CreatedAt: time.Now(),
		}

		if err := repository.DB.Create(&comment).Error; err != nil {
			return nil, err
		}

		userDTO, _ := DefaultUserService.GetUserDTO(userID, userID)

		return &model.CommentDTO{
			ID:         comment.ID,
			User:       userDTO,
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("01-02"),
		}, nil

	} else if actionType == 2 {
		// 2 - 删除评论
		var c model.CommentEntity
		if err := repository.DB.Where("id = ?", commentID).First(&c).Error; err != nil {
			return nil, errors.New("comment not found")
		}
		if c.UserID != userID {
			return nil, errors.New("unauthorized to delete this comment")
		}

		if err := repository.DB.Delete(&c).Error; err != nil {
			return nil, err
		}

		return nil, nil
	}

	return nil, errors.New("invalid action_type")
}

func (s *CommentService) GetCommentList(videoID, currentUserID int64) ([]model.CommentDTO, error) {
	var comments []model.CommentEntity
	err := repository.DB.Where("video_id = ?", videoID).
		Order("created_at DESC").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}

	dtos := make([]model.CommentDTO, 0, len(comments))
	for _, c := range comments {
		userDTO, _ := DefaultUserService.GetUserDTO(c.UserID, currentUserID)
		dtos = append(dtos, model.CommentDTO{
			ID:         c.ID,
			User:       userDTO,
			Content:    c.Content,
			CreateDate: c.CreatedAt.Format("01-02"),
		})
	}

	return dtos, nil
}
