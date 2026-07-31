package controller

import (
	"net/http"
	"strconv"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/MeiSiristhebest/tiktok-backend-go/service"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if username == "" || password == "" {
		username = c.PostForm("username")
		password = c.PostForm("password")
	}

	userID, token, err := service.DefaultUserService.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, model.UserRegisterResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.UserRegisterResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "User registered successfully",
		},
		UserID: userID,
		Token:  token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if username == "" || password == "" {
		username = c.PostForm("username")
		password = c.PostForm("password")
	}

	userID, token, err := service.DefaultUserService.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.UserLoginResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "Login successful",
		},
		UserID: userID,
		Token:  token,
	})
}

func UserInfo(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.UserResponse{
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

	userDTO, err := service.DefaultUserService.GetUserDTO(targetUserID, currentUserID)
	if err != nil {
		c.JSON(http.StatusOK, model.UserResponse{
			BaseResponse: model.BaseResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.UserResponse{
		BaseResponse: model.BaseResponse{
			StatusCode: 0,
			StatusMsg:  "User info fetched successfully",
		},
		User: userDTO,
	})
}
