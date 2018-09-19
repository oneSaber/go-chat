// 用户和数据库相关的内容

package data

import (
	"crypto/md5"
	"errors"
	"strings"
	"time"
)

// 用户注册, 成功返回true, 失败返回false
// /post/register 用户注册
func (user *User) Register() bool {
	// 查询数据库中是否存在该用户
	var isExitUser User
	database.Where("email = ?", user.Email).First(&isExitUser)
	if isExitUser.Id != 0 {
		return false
	}
	database.Create(&user)
	return true
}

//从数据库中从查询用户信息，返回User类型，如果数据库中不存在的话logUser.Id = 0
func QueryUser(queryInfo string) User {
	// 先用email查询
	var user User
	database.Where("email = ?", queryInfo).First(&user)
	if user.Id == 0 {
		// queryInfo 不是email，换用account查询
		database.Where("account = ?", queryInfo).First(&user)
	}
	return user
}

// 登陆创建session, 并向redis中写入session
func (user *User) CreateSession() (session string, err error) {
	sessStr := "Id" + string(user.Id) + "Email" + user.Email + "passwd" + user.Password
	h := md5.New()
	sessHash := h.Sum([]byte(sessStr))
	// session 在redis中保存七天
	err = redisClient.Set(string(user.Id), sessHash, 7*24*60*60*time.Second).Err()
	if err != nil {
		return string(sessHash), err
	}
	return string(sessHash), nil
}

// checklogin， 验证是否登陆
func CheckLogin(userId, cookie string) (valid bool, err error) {
	// 从redis中查询是否有该id的kv并和生成的checkSess进行对比
	val, err := redisClient.Get(userId).Result()
	if err != nil {
		return false, err
	}
	if strings.Compare(cookie, val) != 0 {
		return false, errors.New("session wrong")
	}
	return true, nil
}

// 退出登陆，将session从redis中删除
func Logout(userId string) (err error) {
	err = redisClient.Del(userId).Err()
	if err != nil {
		return err
	}
	return nil
}
