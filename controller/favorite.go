package controller

import (
	"net/http"
	"strconv"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/service"
	"github.com/gin-gonic/gin"
)

func FavoriteAction(c *gin.Context) {
	userIDVal, _ := c.Get("current_user_id")
	userID := userIDVal.(int64)

	videoIDStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")

	if videoIDStr == "" || actionTypeStr == "" {
		videoIDStr = c.PostForm("video_id")
		actionTypeStr = c.PostForm("action_type")
	}

	videoID, err1 := strconv.ParseInt(videoIDStr, 10, 64)
	actionType, err2 := strconv.ParseInt(actionTypeStr, 10, 32)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, model.FavoriteActionResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "Invalid parameters for favorite action",
			},
		})
		return
	}

	err := service.DefaultFavoriteService.FavoriteAction(userID, videoID, int32(actionType))
	if err != nil {
		c.JSON(http.StatusOK, model.FavoriteActionResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.FavoriteActionResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Favorite action processed successfully",
		},
	})
}

func FavoriteList(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.FavoriteListResponse{
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

	videos, err := service.DefaultFavoriteService.GetFavoriteList(targetUserID, currentUserID)
	if err != nil {
		c.JSON(http.StatusOK, model.FavoriteListResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.FavoriteListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Favorite list fetched successfully",
		},
		VideoList: videos,
	})
}
