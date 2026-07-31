package service

import (
	"errors"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/repository"
)

type FavoriteService struct{}

var DefaultFavoriteService = &FavoriteService{}

func (s *FavoriteService) FavoriteAction(userID, videoID int64, actionType int32) error {
	if actionType == 1 {
		// 1 - 点赞
		var count int64
		repository.DB.Model(&model.FavoriteEntity{}).
			Where("user_id = ? AND video_id = ?", userID, videoID).
			Count(&count)
		if count > 0 {
			return nil // 已点赞，直接返回成功
		}
		fav := model.FavoriteEntity{
			UserID:  userID,
			VideoID: videoID,
		}
		return repository.DB.Create(&fav).Error
	} else if actionType == 2 {
		// 2 - 取消点赞
		return repository.DB.Where("user_id = ? AND video_id = ?", userID, videoID).
			Delete(&model.FavoriteEntity{}).Error
	}
	return errors.New("invalid action_type")
}

func (s *FavoriteService) GetFavoriteList(targetUserID, currentUserID int64) ([]model.VideoDTO, error) {
	var favs []model.FavoriteEntity
	if err := repository.DB.Where("user_id = ?", targetUserID).Order("id DESC").Find(&favs).Error; err != nil {
		return nil, err
	}

	dtos := make([]model.VideoDTO, 0, len(favs))
	for _, f := range favs {
		var v model.VideoEntity
		if err := repository.DB.Where("id = ?", f.VideoID).First(&v).Error; err != nil {
			continue
		}

		authorDTO, _ := DefaultUserService.GetUserDTO(v.AuthorID, currentUserID)

		var favCount, comCount int64
		repository.DB.Model(&model.FavoriteEntity{}).Where("video_id = ?", v.ID).Count(&favCount)
		repository.DB.Model(&model.CommentEntity{}).Where("video_id = ?", v.ID).Count(&comCount)

		isFav := false
		if currentUserID > 0 {
			var cnt int64
			repository.DB.Model(&model.FavoriteEntity{}).
				Where("user_id = ? AND video_id = ?", currentUserID, v.ID).
				Count(&cnt)
			isFav = cnt > 0
		}

		dtos = append(dtos, model.VideoDTO{
			ID:            v.ID,
			Author:        authorDTO,
			PlayURL:       v.PlayURL,
			CoverURL:      v.CoverURL,
			FavoriteCount: favCount,
			CommentCount:  comCount,
			IsFavorite:    isFav,
			Title:         v.Title,
		})
	}

	return dtos, nil
}
