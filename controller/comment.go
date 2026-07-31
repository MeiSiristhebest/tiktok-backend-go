package controller

import (
	"net/http"
	"strconv"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/service"
	"github.com/gin-gonic/gin"
)

func CommentAction(c *gin.Context) {
	userIDVal, _ := c.Get("current_user_id")
	userID := userIDVal.(int64)

	videoIDStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")
	commentText := c.Query("comment_text")
	commentIDStr := c.Query("comment_id")

	if videoIDStr == "" || actionTypeStr == "" {
		videoIDStr = c.PostForm("video_id")
		actionTypeStr = c.PostForm("action_type")
		commentText = c.PostForm("comment_text")
		commentIDStr = c.PostForm("comment_id")
	}

	videoID, err1 := strconv.ParseInt(videoIDStr, 10, 64)
	actionType, err2 := strconv.ParseInt(actionTypeStr, 10, 32)
	var commentID int64
	if commentIDStr != "" {
		commentID, _ = strconv.ParseInt(commentIDStr, 10, 64)
	}

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, model.CommentActionResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "Invalid parameters for comment action",
			},
		})
		return
	}

	commentDTO, err := service.DefaultCommentService.CommentAction(userID, videoID, int32(actionType), commentText, commentID)
	if err != nil {
		c.JSON(http.StatusOK, model.CommentActionResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.CommentActionResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Comment action completed successfully",
		},
		Comment: commentDTO,
	})
}

func CommentList(c *gin.Context) {
	videoIDStr := c.Query("video_id")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.CommentListResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "Invalid video_id parameter",
			},
		})
		return
	}

	var currentUserID int64
	if val, exists := c.Get("current_user_id"); exists {
		currentUserID = val.(int64)
	}

	comments, err := service.DefaultCommentService.GetCommentList(videoID, currentUserID)
	if err != nil {
		c.JSON(http.StatusOK, model.CommentListResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.CommentListResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Comment list fetched successfully",
		},
		CommentList: comments,
	})
}
