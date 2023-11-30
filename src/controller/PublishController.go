package controller

import (
	"context"
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/model"
	"douyin/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ReturnUserInfo struct {
	AuthorId      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `mandatory:"false" json:"is_follow"`
}

type ReturnVideo struct {
	VideoId       uint           `json:"video_id"`
	Author        ReturnUserInfo `json:"author"`
	PlayUrl       string         `json:"play_url"`
	CoverUrl      string         `json:"cover_url"`
	FavoriteCount uint           `json:"favorite_count"`
	CommentCount  uint           `json:"comment_count"`
	IsFavorite    bool           `json:"is_favorite"`
	Title         string         `json:"title"`
}

type VideoListResponse struct {
	common.Response
	VideoList []ReturnVideo `json:"video_list"`
}

func Publish(c *gin.Context) { //上传视频方法
	//1.中间件验证token后，获取userId
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	//2.接收请求参数信息
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//3.返回至前端页面的展示信息
	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userId, fileName)
	//先存储到本地文件夹，再保存到云端，获取封面后最后删除
	saveFile := filepath.Join("./resources/static/video/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	coverName := strings.Replace(finalName, ".mp4", ".jpeg", 1)
	img, coverLen := service.ExampleReadFrameAsJpeg(saveFile, 2) //获取第2帧封面

	covImg := io.NopCloser(img)

	//4.保存到云端
	objClient := dao.GetObjectStorageClient()

	//upload cover
	err = dao.PutObject(context.Background(), objClient, dao.Namespace, dao.BucketName, coverName, coverLen, "image/jpeg", covImg, nil)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//upload video
	video, err := os.Open(saveFile)
	if err != nil {
		log.Println(err)
	}

	err = dao.PutObject(context.Background(), objClient, dao.Namespace, dao.BucketName, finalName, data.Size, "video/mp4", video, nil)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//5.删除本地文件
	os.Remove(saveFile)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  finalName + "--uploaded successfully",
	})

	//4.保存发布信息至数据库,刚开始发布，喜爱和评论默认为0
	store_video_info := model.Video{
		Model:         gorm.Model{},
		AuthorId:      userId,
		PlayUrl:       dao.Url + finalName,
		CoverUrl:      dao.Url + coverName,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	service.CreateVideo(&store_video_info)

}

func PublishList(c *gin.Context) {
	//1.中间件鉴权token
	getHostId, _ := c.Get("user_id")
	var HostId uint
	if v, ok := getHostId.(uint); ok {
		HostId = v
	}
	//2.查询要查看用户的id的所有视频，返回页面
	getGuestId := c.Query("user_id")
	id, _ := strconv.Atoi(getGuestId)
	GuestId := uint(id)
	if GuestId == 0 || GuestId == HostId {
		//根据token-id查找用户
		getUser, err := service.GetUser(HostId)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "can't find this user.",
			})
			c.Abort()
			return
		}

		returnMyself := ReturnUserInfo{
			AuthorId:      getUser.ID,
			Name:          getUser.Name,
			FollowCount:   getUser.FollowCount,
			FollowerCount: getUser.FollowerCount,
		}
		//根据用户id查找 所有相关视频信息
		var videoList []model.Video
		videoList = service.GetvideoList(HostId)
		if len(videoList) == 0 {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					StatusCode: 1,
					StatusMsg:  "null",
				},
				VideoList: nil,
			})
		} else { //需要展示的列表信息
			var returnVideoList2 []ReturnVideo
			for i := 0; i < len(videoList); i++ {
				returnVideo2 := ReturnVideo{
					VideoId:       videoList[i].ID,
					Author:        returnMyself,
					PlayUrl:       videoList[i].PlayUrl,
					CoverUrl:      videoList[i].CoverUrl,
					FavoriteCount: videoList[i].FavoriteCount,
					CommentCount:  videoList[i].CommentCount,
					IsFavorite:    service.CheckFavorite(HostId, videoList[i].ID),
					Title:         videoList[i].Title,
				}
				returnVideoList2 = append(returnVideoList2, returnVideo2)
			}
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					StatusCode: 0,
					StatusMsg:  "success",
				},
				VideoList: returnVideoList2,
			})
		}
	} else {
		//根据传入id查找用户
		getUser, err := service.GetUser(GuestId)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "Not find this person.",
			})
			c.Abort()
			return
		}
		returnAuthor := ReturnUserInfo{
			AuthorId:      getUser.ID,
			Name:          getUser.Name,
			FollowCount:   getUser.FollowCount,
			FollowerCount: getUser.FollowerCount,
			IsFollow:      service.IsFollowing(HostId, GuestId),
		}
		//根据用户id查找 所有相关视频信息
		var videoList []model.Video
		videoList = service.GetvideoList(GuestId)
		if len(videoList) == 0 {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					StatusCode: 1,
					StatusMsg:  "null",
				},
				VideoList: nil,
			})
		} else { //需要展示的列表信息
			var returnVideoList []ReturnVideo
			for i := 0; i < len(videoList); i++ {
				returnVideo := ReturnVideo{
					VideoId:       videoList[i].ID,
					Author:        returnAuthor,
					PlayUrl:       videoList[i].PlayUrl,
					CoverUrl:      videoList[i].CoverUrl,
					FavoriteCount: videoList[i].FavoriteCount,
					CommentCount:  videoList[i].CommentCount,
					IsFavorite:    service.CheckFavorite(HostId, videoList[i].ID),
					Title:         videoList[i].Title,
				}
				returnVideoList = append(returnVideoList, returnVideo)
			}
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					StatusCode: 0,
					StatusMsg:  "success",
				},
				VideoList: returnVideoList,
			})
		}
	}
}
