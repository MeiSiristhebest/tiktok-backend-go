package main

import (
	"fmt"

	"log"

	"github.com/MeiSiristhebest/tiktok-backend-go/config"
	"github.com/MeiSiristhebest/tiktok-backend-go/controller"
	"github.com/MeiSiristhebest/tiktok-backend-go/middleware"
	"github.com/MeiSiristhebest/tiktok-backend-go/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	repository.InitDB()

	r := gin.Default()

	// 注册全局中间件
	r.Use(middleware.Cors())

	// 静态文件服务 (视频、封面)
	r.Static("/static", config.GlobalConfig.StaticDir)

	// API 路由组
	douyin := r.Group("/douyin")
	{
		// 1. 基础接口
		douyin.GET("/feed/", middleware.OptionalAuthMiddleware(), controller.Feed)
		douyin.POST("/user/register/", controller.Register)
		douyin.POST("/user/login/", controller.Login)
		douyin.GET("/user/", middleware.OptionalAuthMiddleware(), controller.UserInfo)
		douyin.POST("/publish/action/", middleware.AuthMiddleware(), controller.PublishAction)
		douyin.GET("/publish/list/", middleware.OptionalAuthMiddleware(), controller.PublishList)

		// 2. 互动接口
		douyin.POST("/favorite/action/", middleware.AuthMiddleware(), controller.FavoriteAction)
		douyin.GET("/favorite/list/", middleware.OptionalAuthMiddleware(), controller.FavoriteList)
		douyin.POST("/comment/action/", middleware.AuthMiddleware(), controller.CommentAction)
		douyin.GET("/comment/list/", middleware.OptionalAuthMiddleware(), controller.CommentList)

		// 3. 社交接口
		douyin.POST("/relation/action/", middleware.AuthMiddleware(), controller.RelationAction)
		douyin.GET("/relation/follow/list/", middleware.OptionalAuthMiddleware(), controller.FollowList)
		douyin.GET("/relatioin/follow/list/", middleware.OptionalAuthMiddleware(), controller.FollowList) // 兼容青训营官方错别字文档路径
		douyin.GET("/relation/follower/list/", middleware.OptionalAuthMiddleware(), controller.FollowerList)
		douyin.GET("/relation/friend/list/", middleware.OptionalAuthMiddleware(), controller.FriendList)

		// 消息模块
		douyin.GET("/message/chat/", middleware.AuthMiddleware(), controller.MessageChat)
		douyin.POST("/message/action/", middleware.AuthMiddleware(), controller.MessageAction)
	}

	addr := fmt.Sprintf(":%s", config.GlobalConfig.ServerPort)
	log.Printf("TikTok Backend Server is running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
