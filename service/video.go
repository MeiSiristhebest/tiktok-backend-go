package service

import (
	"fmt"

	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/MeiSiristhebest/tiktok-backend-go/config"
	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/repository"
	"github.com/gin-gonic/gin"
)

type VideoService struct{}

var DefaultVideoService = &VideoService{}

func (s *VideoService) GetFeed(latestTime int64, currentUserID int64) ([]model.VideoDTO, int64, error) {
	var latestTimeVal time.Time
	if latestTime > 0 {
		latestTimeVal = time.Unix(latestTime, 0)
	} else {
		latestTimeVal = time.Now()
	}

	var entities []model.VideoEntity
	err := repository.DB.Where("created_at <= ?", latestTimeVal).
		Order("created_at DESC").
		Limit(30).
		Find(&entities).Error

	if err != nil {
		return nil, 0, err
	}

	dtos := make([]model.VideoDTO, 0, len(entities))
	var nextTime int64 = time.Now().Unix()

	for _, v := range entities {
		authorDTO, _ := DefaultUserService.GetUserDTO(v.AuthorID, currentUserID)

		// 点赞数与评论数
		var favCount, comCount int64
		repository.DB.Model(&model.FavoriteEntity{}).Where("video_id = ?", v.ID).Count(&favCount)
		repository.DB.Model(&model.CommentEntity{}).Where("video_id = ?", v.ID).Count(&comCount)

		// 是否已点赞
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

		if v.CreatedAt.Unix() < nextTime {
			nextTime = v.CreatedAt.Unix()
		}
	}

	return dtos, nextTime, nil
}

func (s *VideoService) PublishVideo(c *gin.Context, userID int64, file *multipart.FileHeader, title string) error {
	staticDir := config.GlobalConfig.StaticDir
	videoDir := filepath.Join(staticDir, "video")
	coverDir := filepath.Join(staticDir, "cover")

	os.MkdirAll(videoDir, 0755)
	os.MkdirAll(coverDir, 0755)

	filename := fmt.Sprintf("%d_%d_%s", userID, time.Now().Unix(), filepath.Base(file.Filename))
	savePath := filepath.Join(videoDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return fmt.Errorf("failed to save uploaded video: %w", err)
	}

	// 生成 URL
	playURL := fmt.Sprintf("%s/static/video/%s", config.GlobalConfig.BaseURL, filename)
	// 默认提供高性能质感封面
	coverURL := "https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=500&auto=format&fit=crop&q=60"

	video := model.VideoEntity{
		AuthorID:  userID,
		PlayURL:   playURL,
		CoverURL:  coverURL,
		Title:     title,
		CreatedAt: time.Now(),
	}

	return repository.DB.Create(&video).Error
}

func (s *VideoService) GetPublishList(targetUserID, currentUserID int64) ([]model.VideoDTO, error) {
	var entities []model.VideoEntity
	err := repository.DB.Where("author_id = ?", targetUserID).
		Order("created_at DESC").
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	dtos := make([]model.VideoDTO, 0, len(entities))
	authorDTO, _ := DefaultUserService.GetUserDTO(targetUserID, currentUserID)

	for _, v := range entities {
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
