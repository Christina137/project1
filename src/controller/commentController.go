package controller

import (
	"douyin/src/common"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	//1 数据处理
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 10)

	// 2 判断评论操作类型：1代表发布评论，2代表删除评论
	//2.1 非合法操作类型
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 405,
			StatusMsg:  "Unsupported actionType",
		})
		c.Abort()
		return
	}

	//2.2 合法操作类型
	if actionType == "1" { // 发布评论
		text := c.Query("comment_text")
		var commentResponse service.CommentActionResponse
		commentResponse, _ = service.PostComment(userId, text, uint(videoId))
		c.JSON(http.StatusOK, commentResponse)
	} else if actionType == "2" { //删除评论
		commentIdStr := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(commentIdStr, 10, 10)
		var commentResponse common.Response
		commentResponse, _ = service.DeleteComment(uint(videoId), uint(commentId))
		c.JSON(http.StatusOK, commentResponse)

	}
}

// CommentList 获取评论表
func CommentList(c *gin.Context) {
	//1 数据处理
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 10)

	//2.调用service层获取指定videoid的评论表
	commentList, err := service.GetCommentList(uint(videoId))

	//2.1 评论表不存在
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 403,
			StatusMsg:  "Failed to get commentList",
		})
		c.Abort()
		return
	}
	//2.2 评论表存在
	var responseCommentList []service.CommentResponse
	for i := 0; i < len(commentList); i++ {
		getUser, err1 := service.GetUser(commentList[i].UserId)

		if err1 != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 403,
				StatusMsg:  "Failed to get commentList.",
			})
			c.Abort()
			return
		}
		responseComment := service.CommentResponse{
			ID:         commentList[i].ID,
			Content:    commentList[i].Content,
			CreateDate: commentList[i].CreatedAt.Format("01-02"), // mm-dd
			User: service.UserResponse{
				ID:            getUser.ID,
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      service.IsFollowing(userId, commentList[i].ID),
			},
		}
		responseCommentList = append(responseCommentList, responseComment)

	}

	//响应返回
	c.JSON(http.StatusOK, service.CommentListResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "Successfully obtained the comment list.",
		},
		CommentList: responseCommentList,
	})

}
