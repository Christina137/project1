package routes

import (
	"douyin/src/controller"
	"douyin/src/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	douyinGroup := r.Group("/douyin")
	{
		userGroup := douyinGroup.Group("/user")
		{
			userGroup.POST("/register/", controller.UserRegister)
			userGroup.POST("/login/", controller.UserLogin)
			userGroup.GET("/", middleware.JwtMiddleware(), controller.UserInfo)
		}
		publishGroup := douyinGroup.Group("/publish")
		{
			publishGroup.POST("/action/", middleware.JwtMiddleware(), controller.Publish)
			publishGroup.GET("/list/", middleware.JwtMiddleware(), controller.PublishList)

		}

		// feed
		douyinGroup.GET("/feed/", controller.Feed)

		favoriteGroup := douyinGroup.Group("favorite")
		{
			favoriteGroup.POST("/action/", middleware.JwtMiddleware(), controller.Favorite)
			favoriteGroup.GET("/list/", middleware.JwtMiddleware(), controller.FavoriteList)
		}

		// comment路由组
		commentGroup := douyinGroup.Group("/comment")
		{
			commentGroup.POST("/action/", middleware.JwtMiddleware(), controller.CommentAction)
			commentGroup.GET("/list/", middleware.JwtMiddleware(), controller.CommentList)
		}

		// relation路由组
		relationGroup := douyinGroup.Group("relation")
		{
			relationGroup.POST("/action/", middleware.JwtMiddleware(), controller.RelationAction)
			relationGroup.GET("/follow/list/", middleware.JwtMiddleware(), controller.FollowList)
			relationGroup.GET("/follower/list/", middleware.JwtMiddleware(), controller.FollowerList)
		}
	}

	return r
}
