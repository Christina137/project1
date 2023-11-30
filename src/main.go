package main

import (
	"douyin/src/conf"
	"douyin/src/dao"
	"douyin/src/routes"
)

func main() {
	conf.InitConfig()
	err := dao.InitDB()
	if err != nil {
		panic(err)
	}

	defer dao.CloseDB()

	r := routes.InitRouter()
	//启动端口为8080的项目
	errRun := r.Run(":8080")
	if errRun != nil {
		return
	}

}
