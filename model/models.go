package model

import "time"

// ==================== API Response Base ====================

type BaseResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// ==================== Domain Models (GORM Database Entities) ====================

type UserEntity struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username        string    `gorm:"uniqueIndex;type:varchar(64);not null" json:"username"`
	Password        string    `gorm:"type:varchar(128);not null" json:"-"`
	Name            string    `gorm:"type:varchar(64)" json:"name"`
	Avatar          string    `gorm:"type:varchar(255)" json:"avatar"`
	BackgroundImage string    `gorm:"type:varchar(255)" json:"background_image"`
	Signature       string    `gorm:"type:varchar(255)" json:"signature"`
	CreatedAt       time.Time `json:"created_at"`
}

func (UserEntity) TableName() string {
	return "users"
}

type VideoEntity struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AuthorID  int64     `gorm:"index:idx_author;not null" json:"author_id"`
	PlayURL   string    `gorm:"type:varchar(255);not null" json:"play_url"`
	CoverURL  string    `gorm:"type:varchar(255);not null" json:"cover_url"`
	Title     string    `gorm:"type:varchar(128);not null" json:"title"`
	CreatedAt time.Time `gorm:"index:idx_video_created_at" json:"created_at"`
}

func (VideoEntity) TableName() string {
	return "videos"
}

type FavoriteEntity struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"uniqueIndex:uk_user_video;not null" json:"user_id"`
	VideoID   int64     `gorm:"uniqueIndex:uk_user_video;index:idx_video;not null" json:"video_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (FavoriteEntity) TableName() string {
	return "favorites"
}

type CommentEntity struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null" json:"user_id"`
	VideoID   int64     `gorm:"index:idx_video_id;not null" json:"video_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (CommentEntity) TableName() string {
	return "comments"
}

type RelationEntity struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID int64     `gorm:"uniqueIndex:uk_follower_followee;index:idx_follower;not null" json:"follower_id"`
	FolloweeID int64     `gorm:"uniqueIndex:uk_follower_followee;index:idx_followee;not null" json:"followee_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (RelationEntity) TableName() string {
	return "relations"
}

type MessageEntity struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FromUserID int64     `gorm:"index:idx_from_to;not null" json:"from_user_id"`
	ToUserID   int64     `gorm:"index:idx_from_to;not null" json:"to_user_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time `gorm:"index:idx_msg_created_at" json:"created_at"`
}

func (MessageEntity) TableName() string {
	return "messages"
}

// ==================== DTO / API Structural Payloads ====================

type UserDTO struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
}

type VideoDTO struct {
	ID            int64   `json:"id"`
	Author        UserDTO `json:"author"`
	PlayURL       string  `json:"play_url"`
	CoverURL      string  `json:"cover_url"`
	FavoriteCount int64   `json:"favorite_count"`
	CommentCount  int64   `json:"comment_count"`
	IsFavorite    bool    `json:"is_favorite"`
	Title         string  `json:"title"`
}

type CommentDTO struct {
	ID         int64   `json:"id"`
	User       UserDTO `json:"user"`
	Content    string  `json:"content"`
	CreateDate string  `json:"create_date"` // 格式 mm-dd
}

type FriendUserDTO struct {
	UserDTO
	Message string `json:"message,omitempty"`
	MsgType int64  `json:"msgType"` // 0 => 当前用户接收的消息, 1 => 当前用户发送的消息
}

type MessageDTO struct {
	ID         int64  `json:"id"`
	ToUserID   int64  `json:"to_user_id"`
	FromUserID int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}

// ==================== API Request / Response Structs ====================

type UserRegisterResponse struct {
	BaseResponse
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserLoginResponse struct {
	BaseResponse
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserResponse struct {
	BaseResponse
	User UserDTO `json:"user"`
}

type FeedResponse struct {
	BaseResponse
	VideoList []VideoDTO `json:"video_list"`
	NextTime  int64      `json:"next_time,omitempty"`
}

type PublishActionResponse struct {
	BaseResponse
}

type PublishListResponse struct {
	BaseResponse
	VideoList []VideoDTO `json:"video_list"`
}

type FavoriteActionResponse struct {
	BaseResponse
}

type FavoriteListResponse struct {
	BaseResponse
	VideoList []VideoDTO `json:"video_list"`
}

type CommentActionResponse struct {
	BaseResponse
	Comment *CommentDTO `json:"comment,omitempty"`
}

type CommentListResponse struct {
	BaseResponse
	CommentList []CommentDTO `json:"comment_list"`
}

type RelationActionResponse struct {
	BaseResponse
}

type RelationUserListResponse struct {
	BaseResponse
	UserList []UserDTO `json:"user_list"`
}

type RelationFriendListResponse struct {
	BaseResponse
	UserList []FriendUserDTO `json:"user_list"`
}

type MessageChatResponse struct {
	BaseResponse
	MessageList []MessageDTO `json:"message_list"`
}

type MessageActionResponse struct {
	BaseResponse
}
