package controller

import (
	"net/http"
	"strconv"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/service"
	"github.com/gin-gonic/gin"
)

func Feed(c *gin.Context) {
	latestTimeStr := c.Query("latest_time")
	var latestTime int64
	if latestTimeStr != "" {
		latestTime, _ = strconv.ParseInt(latestTimeStr, 10, 64)
	}

	var currentUserID int64
	if val, exists := c.Get("current_user_id"); exists {
		currentUserID = val.(int64)
	}

	videos, nextTime, err := service.DefaultVideoService.GetFeed(latestTime, currentUserID)
	if err != nil {
		c.JSON(http.StatusOK, model.FeedResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.FeedResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Feed loaded successfully",
		},
		VideoList: videos,
		NextTime:  nextTime,
	})
}

func PublishAction(c *gin.Context) {
	userIDVal, _ := c.Get("current_user_id")
	userID := userIDVal.(int64)

	title := c.PostForm("title")
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, model.PublishActionResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "Uploaded video file payload not found: " + err.Error(),
			},
		})
		return
	}

	err = service.DefaultVideoService.PublishVideo(c, userID, file, title)
	if err != nil {
		c.JSON(http.StatusOK, model.PublishActionResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "Failed to publish video: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.PublishActionResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Video published successfully",
		},
	})
}

func PublishList(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.PublishListResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "Invalid user_id parameter",
			},
		})
		return
	}

	var currentUserID int64
	if val, exists := c.Get("current_user_id"); exists {
		currentUserID = val.(int64)
	}

	videos, err := service.DefaultVideoService.GetPublishList(targetUserID, currentUserID)
	if err != nil {
		c.JSON(http.StatusOK, model.PublishListResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.PublishListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Publish list fetched successfully",
		},
		VideoList: videos,
	})
}
