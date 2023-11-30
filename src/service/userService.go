package service

import (
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/middleware"
	"douyin/src/model"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

const (
	MaxUsernameLength = 20
	MaxPasswordLength = 20
	MinPasswordLength = 6
)

type UserIdTokenResponse struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

func CreateRegisterUser(userName string, password string) (model.User, error) {
	pwdHash, _ := PasswordEncoder(password)

	newUser := model.User{
		Name:     userName,
		Password: pwdHash,
	}

	dao.SqlSession.AutoMigrate(&model.User{})

	if IsUserExistByName(userName) {
		return newUser, common.ErrorUserExit
	} else {
		if err := dao.SqlSession.Model(&model.User{}).Create(&newUser).Error; err != nil {
			panic(err)
		}
	}
	return newUser, nil
}

func IsUserExistByName(userName string) bool {
	var user = &model.User{}
	if err := dao.SqlSession.Model(&model.User{}).Where("name=?", userName).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}

	return true
}

func IsUserExist(userName string, password string, login *model.User) error {
	if login == nil {
		return common.ErrorNullPointer
	}
	dao.SqlSession.Where("name=?", userName).First(login)
	if !PasswordVerify(login.Password, password) {
		return common.ErrorPasswordFalse
	}
	if login.Model.ID == 0 {
		return common.ErrorFullPossibility
	}
	return nil
}

func PasswordEncoder(password string) (pwdHash string, err error) {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return
	}
	pwdHash = string(hash)
	return
}

func PasswordVerify(password string, pwdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(pwdHash))
	if err != nil {
		return false
	}
	return true
}

func GetUser(userId uint) (model.User, error) {
	var user = &model.User{}
	if err := dao.SqlSession.Model(&model.User{}).Where("id=?", userId).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return *user, err
	}
	return *user, nil

}
func GetUserById(userId uint, user *model.User) error {
	if user == nil {
		return common.ErrorNullPointer
	}
	dao.SqlSession.Where("id=?", userId).First(user)
	return nil
}

// IsUserLegal 用户名和密码合法性检验
func IsUserLegal(userName string, passWord string) error {
	//1.用户名检验
	if userName == "" {
		return common.ErrorUserNameNull
	}
	if len(userName) > MaxUsernameLength {
		return common.ErrorUserNameExtend
	}
	//2.密码检验
	if passWord == "" {
		return common.ErrorPasswordNull
	}
	if len(passWord) > MaxPasswordLength || len(passWord) < MinPasswordLength {
		return common.ErrorPasswordLength
	}
	return nil
}

func UserRegister(userName string, passWord string) (UserIdTokenResponse, error) {
	var response = UserIdTokenResponse{}
	//1.用户名和密码合法性检验
	err := IsUserLegal(userName, passWord)
	if err != nil {
		return response, err
	}
	//2.用户名是否已经存在
	if IsUserExistByName(userName) {
		return response, common.ErrorUserExit
	}
	//3.创建用户
	user, err := CreateRegisterUser(userName, passWord)
	if err != nil {
		return response, err
	}
	//4.创建token
	token, err := middleware.CreateToken(user.ID, user.Name)
	if err != nil {
		return response, err
	}
	response.UserId = user.ID
	response.Token = token
	return response, nil
}

func UserLoginService(userName string, passWord string) (UserIdTokenResponse, error) {
	var response = UserIdTokenResponse{}
	//1.用户名和密码合法性检验
	err := IsUserLegal(userName, passWord)
	if err != nil {
		return response, err
	}
	//2.用户名是否已经存在
	var user = &model.User{}
	err = IsUserExist(userName, passWord, user)
	if err != nil {
		return response, err
	}
	//3.创建token
	token, err := middleware.CreateToken(user.ID, user.Name)
	if err != nil {
		return response, err
	}
	response.UserId = user.ID
	response.Token = token
	return response, nil
}

type UserInfoQueryResponse struct {
	UserId         uint   `json:"user_id"`
	UserName       string `json:"name"`
	FollowCount    uint   `json:"follow_count"`
	FollowerCount  uint   `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	TotalFavorited uint   `json:"total_favorited"`
	FavoriteCount  uint   `json:"favorite_count"`
}

func UserInfoService(userId string) (UserInfoQueryResponse, error) {
	var response = UserInfoQueryResponse{}
	//1.用户id是否合法
	userIdUint, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return response, err
	}
	//2.用户是否存在
	user := &model.User{}
	err = GetUserById(uint(userIdUint), user)
	if err != nil {
		return response, err
	}

	//3.获取用户信息
	response.UserId = user.ID
	response.UserName = user.Name
	response.FollowCount = user.FollowCount
	response.FollowerCount = user.FollowerCount
	response.TotalFavorited = user.TotalFavorited
	response.FavoriteCount = user.FavoriteCount
	response.IsFollow = false
	return response, nil
}

// CheckIsFollow 检验已登录用户是否关注目标用户
func CheckIsFollow(targetId string, userid uint) bool {
	//1.修改targetId数据类型
	hostId, err := strconv.ParseUint(targetId, 10, 64)
	if err != nil {
		return false
	}
	//如果是自己查自己，那就是没有关注
	if uint(hostId) == userid {
		return false
	}
	//2.自己是否关注目标userId
	//TODO
	return true
}
