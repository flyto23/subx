package api

import (
	"log"
	"sublink/models"
	"sublink/utils"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
}

// 新增用户
func UserAdd(c *gin.Context) {
	user := &models.User{
		Username: "test",
		Password: "test",
	}
	err := user.Create()
	if err != nil {
		log.Println("创建用户失败")
		utils.Fail(c, "创建用户失败")
		return
	}
	utils.Success(c, nil, "创建用户成功")
}

// 获取用户信息
func UserMe(c *gin.Context) {
	// 获取 jwt 中的 username
	// 返回用户信息
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	err := user.Find()
	if err != nil {
		utils.Error(c, "40001", err.Error())
		return
	}
	utils.Success(c, gin.H{
		"avatar":   "static/avatar.gif",
		"nickname": user.Nickname,
		"userId":   user.ID,
		"username": user.Username,
		"roles":    []string{"ADMIN"},
	}, "获取用户信息成功")
}

// 获取所有用户
func UserPages(c *gin.Context) {
	// 获取 jwt 中的 username
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	users, err := user.All()
	if err != nil {
		log.Println("获取用户信息失败")
		utils.Fail(c, "获取用户信息失败")
		return
	}
	list := []*User{}
	for i := range users {
		list = append(list, &User{
			ID:       users[i].ID,
			Username: users[i].Username,
			Nickname: users[i].Nickname,
			Avatar:   "static/avatar.gif",
		})
	}
	utils.Success(c, gin.H{
		"list": list,
	}, "获取用户信息成功")
}

// 更新用户信息
func UserSet(c *gin.Context) {
	NewUsername := c.PostForm("username")
	NewPassword := c.PostForm("password")
	log.Println(NewUsername, NewPassword)
	if NewUsername == "" || NewPassword == "" {
		utils.Error(c, "40001", "用户名或密码不能为空")
		return
	}
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	err := user.Set(&models.User{
		Username: NewUsername,
		Password: NewPassword,
	})
	if err != nil {
		log.Println(err)
		utils.Error(c, "50001", err.Error())
		return
	}
	// 修改成功
	utils.Success(c, nil, "修改成功")
}
