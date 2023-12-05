package middleware

import (
	"douyin/src/common"
	"douyin/src/conf"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type MyClaims struct {
	UserId   uint   `json:"user_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

var key = []byte(conf.Config.Jwt.Secret)

func CreateToken(userId uint, userName string) (string, error) {
	expireTime := time.Now().Add(24 * time.Hour)
	timeNow := time.Now()
	claims := MyClaims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "chris",
			IssuedAt:  timeNow.Unix(),
			Subject:   "JwtToken",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)

}

func CheckToken(tokenString string) (*MyClaims, bool) {
	var claims MyClaims
	token, _ := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, e error) {
		return key, nil
	})
	if token.Valid {
		return &claims, true
	}
	return nil, false
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.PostForm("token")

		if tokenStr == "" {
			tokenStr = c.Query("token")
		}

		// log.Println("token: ", tokenStr)

		if tokenStr == "" {
			c.JSON(http.StatusOK, common.Response{StatusCode: 401, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		//验证token
		tokenStruck, ok := CheckToken(tokenStr)
		if !ok {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 403,
				StatusMsg:  "token不正确",
			})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 402,
				StatusMsg:  "token过期",
			})
			c.Abort() //阻止执行
			return
		}
		c.Set("username", tokenStruck.UserName)
		c.Set("user_id", tokenStruck.UserId)

		c.Next()
	}
}
