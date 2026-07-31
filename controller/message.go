package controller

import (
	"net/http"
	"strconv"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/service"
	"github.com/gin-gonic/gin"
)

func MessageChat(c *gin.Context) {
	userIDVal, _ := c.Get("current_user_id")
	fromUserID := userIDVal.(int64)

	toUserIDStr := c.Query("to_user_id")
	preMsgTimeStr := c.Query("pre_msg_time")

	toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.MessageChatResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "Invalid to_user_id parameter",
			},
		})
		return
	}

	var preMsgTime int64
	if preMsgTimeStr != "" {
		preMsgTime, _ = strconv.ParseInt(preMsgTimeStr, 10, 64)
	}

	messages, err := service.DefaultMessageService.GetChatHistory(fromUserID, toUserID, preMsgTime)
	if err != nil {
		c.JSON(http.StatusOK, model.MessageChatResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.MessageChatResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Message history fetched successfully",
		},
		MessageList: messages,
	})
}

func MessageAction(c *gin.Context) {
	userIDVal, _ := c.Get("current_user_id")
	fromUserID := userIDVal.(int64)

	toUserIDStr := c.Query("to_user_id")
	actionTypeStr := c.Query("action_type")
	content := c.Query("content")

	if toUserIDStr == "" || actionTypeStr == "" {
		toUserIDStr = c.PostForm("to_user_id")
		actionTypeStr = c.PostForm("action_type")
		content = c.PostForm("content")
	}

	toUserID, err1 := strconv.ParseInt(toUserIDStr, 10, 64)
	actionType, err2 := strconv.ParseInt(actionTypeStr, 10, 32)

	if err1 != nil || err2 != nil || actionType != 1 {
		c.JSON(http.StatusOK, model.MessageActionResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  "Invalid message action parameters",
			},
		})
		return
	}

	err := service.DefaultMessageService.SendMessage(fromUserID, toUserID, content)
	if err != nil {
		c.JSON(http.StatusOK, model.MessageActionResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.MessageActionResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Message sent successfully",
		},
	})
}
