package service

import (
	"errors"
	"time"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/repository"
)

type MessageService struct{}

var DefaultMessageService = &MessageService{}

func (s *MessageService) SendMessage(fromUserID, toUserID int64, content string) error {
	if len(content) == 0 {
		return errors.New("message content cannot be empty")
	}

	msg := model.MessageEntity{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	return repository.DB.Create(&msg).Error
}

func (s *MessageService) GetChatHistory(fromUserID, toUserID int64, preMsgTime int64) ([]model.MessageDTO, error) {
	db := repository.DB.Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
		fromUserID, toUserID, toUserID, fromUserID)

	if preMsgTime > 0 {
		// 转换毫秒或秒戳
		var filterTime time.Time
		if preMsgTime > 1e11 {
			filterTime = time.UnixMilli(preMsgTime)
		} else {
			filterTime = time.Unix(preMsgTime, 0)
		}
		db = db.Where("created_at > ?", filterTime)
	}

	var msgs []model.MessageEntity
	if err := db.Order("created_at ASC").Find(&msgs).Error; err != nil {
		return nil, err
	}

	dtos := make([]model.MessageDTO, 0, len(msgs))
	for _, m := range msgs {
		dtos = append(dtos, model.MessageDTO{
			ID:         m.ID,
			ToUserID:   m.ToUserID,
			FromUserID: m.FromUserID,
			Content:    m.Content,
			CreateTime: m.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return dtos, nil
}
