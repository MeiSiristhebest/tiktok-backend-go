package service

import (
	"errors"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/repository"
)

type RelationService struct{}

var DefaultRelationService = &RelationService{}

func (s *RelationService) RelationAction(userID, toUserID int64, actionType int32) error {
	if userID == toUserID {
		return errors.New("cannot follow yourself")
	}

	if actionType == 1 {
		// 1 - 关注
		var count int64
		repository.DB.Model(&model.RelationEntity{}).
			Where("follower_id = ? AND followee_id = ?", userID, toUserID).
			Count(&count)
		if count > 0 {
			return nil
		}

		rel := model.RelationEntity{
			FollowerID: userID,
			FolloweeID: toUserID,
		}
		return repository.DB.Create(&rel).Error
	} else if actionType == 2 {
		// 2 - 取消关注
		return repository.DB.Where("follower_id = ? AND followee_id = ?", userID, toUserID).
			Delete(&model.RelationEntity{}).Error
	}

	return errors.New("invalid action_type")
}

func (s *RelationService) GetFollowList(targetUserID, currentUserID int64) ([]model.UserDTO, error) {
	var rels []model.RelationEntity
	if err := repository.DB.Where("follower_id = ?", targetUserID).Order("id DESC").Find(&rels).Error; err != nil {
		return nil, err
	}

	userDTOs := make([]model.UserDTO, 0, len(rels))
	for _, r := range rels {
		uDTO, err := DefaultUserService.GetUserDTO(r.FolloweeID, currentUserID)
		if err == nil {
			userDTOs = append(userDTOs, uDTO)
		}
	}

	return userDTOs, nil
}

func (s *RelationService) GetFollowerList(targetUserID, currentUserID int64) ([]model.UserDTO, error) {
	var rels []model.RelationEntity
	if err := repository.DB.Where("followee_id = ?", targetUserID).Order("id DESC").Find(&rels).Error; err != nil {
		return nil, err
	}

	userDTOs := make([]model.UserDTO, 0, len(rels))
	for _, r := range rels {
		uDTO, err := DefaultUserService.GetUserDTO(r.FollowerID, currentUserID)
		if err == nil {
			userDTOs = append(userDTOs, uDTO)
		}
	}

	return userDTOs, nil
}

func (s *RelationService) GetFriendList(targetUserID, currentUserID int64) ([]model.FriendUserDTO, error) {
	// 查找 targetUserID 关注的所有用户
	var follows []int64
	repository.DB.Model(&model.RelationEntity{}).
		Where("follower_id = ?", targetUserID).
		Pluck("followee_id", &follows)

	if len(follows) == 0 {
		return []model.FriendUserDTO{}, nil
	}

	// 过滤出同时也关注了 targetUserID 的用户 (即双向关注好友)
	var friendIDs []int64
	repository.DB.Model(&model.RelationEntity{}).
		Where("follower_id IN ? AND followee_id = ?", follows, targetUserID).
		Pluck("follower_id", &friendIDs)

	friendDTOs := make([]model.FriendUserDTO, 0, len(friendIDs))

	for _, fid := range friendIDs {
		uDTO, err := DefaultUserService.GetUserDTO(fid, currentUserID)
		if err != nil {
			continue
		}

		// 查最新的聊天记录
		var msg model.MessageEntity
		err = repository.DB.Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
			targetUserID, fid, fid, targetUserID).
			Order("created_at DESC").
			First(&msg).Error

		latestMsg := ""
		var msgType int64 = 0 // 0: 接收, 1: 发送

		if err == nil {
			latestMsg = msg.Content
			if msg.FromUserID == targetUserID {
				msgType = 1
			} else {
				msgType = 0
			}
		}

		friendDTOs = append(friendDTOs, model.FriendUserDTO{
			UserDTO: uDTO,
			Message: latestMsg,
			MsgType: msgType,
		})
	}

	return friendDTOs, nil
}
