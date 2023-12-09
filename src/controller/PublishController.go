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

func Publish(c *gin.Context) { 
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userId, fileName)
	saveFile := filepath.Join("./resources/static/video/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	coverName := strings.Replace(finalName, ".mp4", ".jpeg", 1)
	img, coverLen := service.ExampleReadFrameAsJpeg(saveFile, 2) 

	covImg := io.NopCloser(img)

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

	os.Remove(saveFile)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  finalName + "--uploaded successfully",
	})

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

	getHostId, _ := c.Get("user_id")
	var HostId uint
	if v, ok := getHostId.(uint); ok {
		HostId = v
	}
	
	getGuestId := c.Query("user_id")
	id, _ := strconv.Atoi(getGuestId)
	GuestId := uint(id)
	if GuestId == 0 || GuestId == HostId {
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
		} else { 
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
		} else { 
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
