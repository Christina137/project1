package controller

import (
	"douyin/src/common"
	"douyin/src/middleware"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRegisterResponse struct {
	common.Response
	service.UserIdTokenResponse
}

func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	registerResponse, err := service.UserRegister(username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response:            common.Response{StatusCode: 0},
		UserIdTokenResponse: registerResponse,
	})

	return
}

func UserLogin(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")

	loginResponse, err := service.UserLoginService(username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response:            common.Response{StatusCode: 0},
		UserIdTokenResponse: loginResponse,
	})

}

type UserInfoResponse struct {
	common.Response
	UserList service.UserInfoQueryResponse `json:"user"`
}

func UserInfo(c *gin.Context) {
	userId := c.Query("user_id")
	Response, err := service.UserInfoService(userId)
	token := c.Query("token")

	hostStruct, _ := middleware.CheckToken(token)
	Response.IsFollow = service.CheckIsFollow(userId, hostStruct.UserId)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "登录成功",
		},
		UserList: Response,
	})

}
