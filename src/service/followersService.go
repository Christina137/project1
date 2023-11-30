package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"github.com/jinzhu/gorm"
)

func CreateFollower(HostId uint, GuestId uint) error {

	//1.Following数据模型准备
	newFollower := model.Followers{
		HostId:  HostId,
		GuestId: GuestId,
	}

	//2.新建following
	if err := dao.SqlSession.Model(&model.Followers{}).
		Create(&newFollower).Error; err != nil {
		return err
	}
	return nil
}

// IncreaseFollowerCount 增加HostId的粉丝数（Host_id 的 follow_count+1）
func IncreaseFollowerCount(HostId uint) error {
	if err := dao.SqlSession.Model(&model.User{}).
		Where("id=?", HostId).
		Update("follower_count", gorm.Expr("follower_count+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// DeleteFollower 删除粉丝
func DeleteFollower(HostId uint, GuestId uint) error {
	//1.Following数据模型准备
	newFollower := model.Followers{
		HostId:  HostId,
		GuestId: GuestId,
	}

	//2.删除following
	if err := dao.SqlSession.Model(&model.Followers{}).
		Where("host_id=? AND guest_id=?", HostId, GuestId).
		Delete(&newFollower).Error; err != nil {
		return err
	}

	return nil
}

// DecreaseFollowerCount 增加HostId的粉丝数（Host_id 的 follow_count-1）
func DecreaseFollowerCount(HostId uint) error {
	if err := dao.SqlSession.Model(&model.User{}).
		Where("id=?", HostId).
		Update("follower_count", gorm.Expr("follower_count-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// FollowerList  获取粉丝表
func FollowerList(HostId uint) ([]model.User, error) {
	//1.userList数据模型准备
	var userList []model.User
	//2.查HostId的关注表
	if err := dao.SqlSession.Model(&model.User{}).
		Joins("left join followers on "+"users.id = "+"followers.guest_id").
		Where("followers.host_id=? AND followers.deleted_at is null", HostId).
		Scan(&userList).Error; err != nil {
		return userList, nil
	}
	return userList, nil
}
