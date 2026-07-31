package controller

import (
	"net/http"
	"strconv"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/service"
	"github.com/gin-gonic/gin"
)

func RelationAction(c *gin.Context) {
	userIDVal, _ := c.Get("current_user_id")
	userID := userIDVal.(int64)

	toUserIDStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type")

	if toUserIDStr == "" || actionTypeStr == "" {
		toUserIDStr = c.PostForm("to_user_id")
		actionTypeStr = c.PostForm("action_type")
	}

	toUserID, err1 := strconv.ParseInt(toUserIDStr, 10, 64)
	actionType, err2 := strconv.ParseInt(actionTypeStr, 10, 32)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, model.RelationActionResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "Invalid parameters for relation action",
			},
		})
		return
	}

	err := service.DefaultRelationService.RelationAction(userID, toUserID, int32(actionType))
	if err != nil {
		c.JSON(http.StatusOK, model.RelationActionResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.RelationActionResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Relation action processed successfully",
		},
	})
}

func FollowList(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.RelationUserListResponse{
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

	users, err := service.DefaultRelationService.GetFollowList(targetUserID, currentUserID)
	if err != nil {
		c.JSON(http.StatusOK, model.RelationUserListResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.RelationUserListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Follow list fetched successfully",
		},
		UserList: users,
	})
}

func FollowerList(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.RelationUserListResponse{
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

	users, err := service.DefaultRelationService.GetFollowerList(targetUserID, currentUserID)
	if err != nil {
		c.JSON(http.StatusOK, model.RelationUserListResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.RelationUserListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Follower list fetched successfully",
		},
		UserList: users,
	})
}

func FriendList(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.RelationFriendListResponse{
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

	friends, err := service.DefaultRelationService.GetFriendList(targetUserID, currentUserID)
	if err != nil {
		c.JSON(http.StatusOK, model.RelationFriendListResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.RelationFriendListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Friend list fetched successfully",
		},
		UserList: friends,
	})
}
