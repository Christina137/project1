package service

import (
	"bytes"
	"douyin/src/dao"
	"douyin/src/model"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
	"time"
)

const videoNum = 2

func CreateVideo(video *model.Video) {
	dao.SqlSession.Table("videos").Create(&video)
}

func ExampleReadFrameAsJpeg(inFileName string, frameNum int) (io.Reader, int64) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf, int64(buf.Len())
}

func GetvideoList(userId uint) []model.Video {
	var videoList []model.Video
	if err := dao.SqlSession.Table("videos").Where("author_id = ?", userId).Find(&videoList).Error; err != nil {
		return nil
	}
	return videoList
}

func GetVideoAuthor(videoId uint) (uint, error) {
	var video model.Video
	if err := dao.SqlSession.Table("videos").Where("id = ?", videoId).Find(&video).Error; err != nil {
		return video.ID, err
	}
	return video.AuthorId, nil
}

func FeedGet(lastTime int64) ([]model.Video, error) {
	if lastTime == 0 { //没有传入参数或者视屏已经刷完
		lastTime = time.Now().Unix()
	}
	strTime := fmt.Sprint(time.Unix(lastTime, 0).Format("2006-01-02 15:04:05"))
	fmt.Println("查询的时间", strTime)
	var VideoList []model.Video
	VideoList = make([]model.Video, 0)
	err := dao.SqlSession.Table("videos").Where("created_at < ?", strTime).Order("created_at desc").Limit(videoNum).Find(&VideoList).Error
	return VideoList, err
}
