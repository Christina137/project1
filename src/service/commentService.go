package service

import (
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/model"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

// CommentListResponse 评论表的响应结构体
type CommentListResponse struct {
	common.Response
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}

// CommentActionResponse 评论操作的响应结构体
type CommentActionResponse struct {
	common.Response
	Comment CommentResponse `json:"comment,omitempty"`
}

// UserResponse 用户信息的响应结构体
type UserResponse struct {
	ID            uint   `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   uint   `json:"follow_count,omitempty"`
	FollowerCount uint   `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// CommentResponse 评论信息的响应结构体
type CommentResponse struct {
	ID         uint         `json:"id,omitempty"`
	Content    string       `json:"content,omitempty"`
	CreateDate string       `json:"create_date,omitempty"`
	User       UserResponse `json:"user,omitempty"`
}

// GetCommentList 获取指定videoId的评论表
func GetCommentList(videoId uint) ([]model.Comment, error) {
	var commentList []model.Comment
	if err := dao.SqlSession.Table("comments").Where("video_id=?", videoId).Find(&commentList).Error; err != nil {
		return commentList, err
	}
	return commentList, nil
}

// PostComment2DB 发布评论
func PostComment2DB(comment model.Comment) error {
	if err := dao.SqlSession.Table("comments").Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

// DeleteComment2DB 删除指定commentId的评论
func DeleteComment2DB(commentId uint) error {
	if err := dao.SqlSession.Table("comments").Where("id = ?", commentId).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

// AddCommentCount add comment_count
func AddCommentCount(videoId uint) error {

	if err := dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		return err
	}
	return nil
}

// ReduceCommentCount reduce comment_count
func ReduceCommentCount(videoId uint) error {

	if err := dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		return err
	}
	return nil
}

// PostComment 发布评论
func PostComment(userId uint, text string, videoId uint) (CommentActionResponse, error) {
	//1 准备数据模型
	newComment := model.Comment{
		VideoId: videoId,
		UserId:  userId,
		Content: text,
	}

	//2 调用service层发布评论并改变评论数量，获取video作者信息
	err1 := dao.SqlSession.Transaction(func(db *gorm.DB) error {
		if err := PostComment2DB(newComment); err != nil {
			return err
		}
		if err := AddCommentCount(videoId); err != nil {
			return err
		}
		return nil
	})
	getUser, err2 := GetUser(userId)
	videoAuthor, err3 := GetVideoAuthor(videoId)
	//3 响应处理
	if err1 != nil || err2 != nil || err3 != nil {
		return CommentActionResponse{
			Response: common.Response{
				StatusCode: http.StatusInternalServerError,
				StatusMsg:  "Post comment failed" + err1.Error() + err2.Error() + err3.Error(),
			},
		}, err1
	}
	response := CommentActionResponse{Response: common.Response{
		StatusCode: 0,
		StatusMsg:  "post the comment successfully",
	},
		Comment: CommentResponse{
			ID:         newComment.ID,
			Content:    newComment.Content,
			CreateDate: newComment.CreatedAt.Format("01-02"),
			User: UserResponse{
				ID:            getUser.ID,
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      IsFollowing(userId, videoAuthor),
			},
		}}
	return response, nil
}

// DeleteComment 删除评论
func DeleteComment(videoId uint, commentId uint) (common.Response, error) {
	//1 调用service层删除评论并改变评论数量，获取video作者信息
	err := dao.SqlSession.Transaction(func(db *gorm.DB) error {
		if err := DeleteComment2DB(commentId); err != nil {
			return err
		}
		if err := ReduceCommentCount(videoId); err != nil {
			return err
		}
		return nil
	})
	//2 响应处理
	if err != nil {

		return common.Response{
			StatusCode: http.StatusInternalServerError,
			StatusMsg:  "Post comment failed" + err.Error(),
		}, err
	}

	response := common.Response{
		StatusCode: 0,
		StatusMsg:  "delete the comment successfully",
	}

	return response, nil
}
