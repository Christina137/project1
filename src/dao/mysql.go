package dao

import (
	"douyin/src/conf"
	"douyin/src/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const DRIVER = "mysql"

var SqlSession *gorm.DB

func InitDB() (err error) {
	var c = conf.Config.Mysql
	url := c.UserName + ":" + c.PassWord + "@tcp(" + c.Url + ":" + c.Port + ")/" + c.DbName + "?charset=utf8&parseTime=True&loc=Local"
	fmt.Println()
	fmt.Println(url)
	SqlSession, err = gorm.Open(DRIVER, url)
	if err != nil {
		panic(err)
	}

	SqlSession.AutoMigrate(&model.User{}, &model.Video{}, &model.Comment{}, &model.Following{}, &model.Followers{}, &model.Favorite{})
	return SqlSession.DB().Ping()
}

func CloseDB() {
	err := SqlSession.Close()
	if err != nil {
		return
	}
}
