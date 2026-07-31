package service

import (
	"errors"

	"github.com/MeiSiristhebest/tiktok-backend-go/middleware"
	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

var DefaultUserService = &UserService{}

func (s *UserService) Register(username, password string) (int64, string, error) {
	if len(username) == 0 || len(username) > 32 || len(password) == 0 || len(password) > 32 {
		return 0, "", errors.New("username or password length invalid (1-32 chars)")
	}

	var count int64
	repository.DB.Model(&model.UserEntity{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return 0, "", errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, "", errors.New("failed to encrypt password")
	}

	newUser := model.UserEntity{
		Username:        username,
		Password:        string(hashedPassword),
		Name:            username,
		Avatar:          "https://p3-passport.byteimg.com/img/user-avatar/7f42841bc6c1e33c94ed6ec8c75d40a0~120x120.awebp",
		BackgroundImage: "https://p3-passport.byteimg.com/img/user-avatar/7f42841bc6c1e33c94ed6ec8c75d40a0~120x120.awebp",
		Signature:       "追求卓越，极简抖音后端开发者。",
	}

	if err := repository.DB.Create(&newUser).Error; err != nil {
		return 0, "", err
	}

	token, err := middleware.GenerateToken(newUser.ID)
	if err != nil {
		return 0, "", err
	}

	return newUser.ID, token, nil
}

func (s *UserService) Login(username, password string) (int64, string, error) {
	var user model.UserEntity
	if err := repository.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, "", errors.New("incorrect password")
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}

func (s *UserService) GetUserDTO(targetUserID, currentUserID int64) (model.UserDTO, error) {
	var user model.UserEntity
	if err := repository.DB.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		return model.UserDTO{}, errors.New("user not found")
	}

	var followCount, followerCount, workCount, favoriteCount, totalFavorited int64

	// 关注总数
	repository.DB.Model(&model.RelationEntity{}).Where("follower_id = ?", targetUserID).Count(&followCount)
	// 粉丝总数
	repository.DB.Model(&model.RelationEntity{}).Where("followee_id = ?", targetUserID).Count(&followerCount)
	// 作品数量
	repository.DB.Model(&model.VideoEntity{}).Where("author_id = ?", targetUserID).Count(&workCount)
	// 点赞(喜欢)视频总数
	repository.DB.Model(&model.FavoriteEntity{}).Where("user_id = ?", targetUserID).Count(&favoriteCount)

	// 该用户发布的视频被赞的总次数
	var authorVideoIDs []int64
	repository.DB.Model(&model.VideoEntity{}).Where("author_id = ?", targetUserID).Pluck("id", &authorVideoIDs)
	if len(authorVideoIDs) > 0 {
		repository.DB.Model(&model.FavoriteEntity{}).Where("video_id IN ?", authorVideoIDs).Count(&totalFavorited)
	}

	// 当前登录用户是否已关注 targetUserID
	isFollow := false
	if currentUserID > 0 && currentUserID != targetUserID {
		var relCount int64
		repository.DB.Model(&model.RelationEntity{}).
			Where("follower_id = ? AND followee_id = ?", currentUserID, targetUserID).
			Count(&relCount)
		isFollow = relCount > 0
	}

	return model.UserDTO{
		ID:              user.ID,
		Name:            user.Name,
		FollowCount:     followCount,
		FollowerCount:   followerCount,
		IsFollow:        isFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  totalFavorited,
		WorkCount:       workCount,
		FavoriteCount:   favoriteCount,
	}, nil
}
