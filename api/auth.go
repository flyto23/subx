package api

import (
	"log"
	"net/http"
	"sublink/middlewares"
	"sublink/models"
	"sublink/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 获取 token
func GetToken(username string) (string, error) {
	// 过期时间天
	ExpireDays := models.ReadConfig().ExpireDays
	c := &middlewares.JwtClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(ExpireDays)).Unix(), // 设置过期时间
			IssuedAt:  time.Now().Unix(),                                                 // 签发时间
			Subject:   username,                                                          // 用户
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(middlewares.Secret)
}

// 获取 captcha 图形验证码
func GetCaptcha(c *gin.Context) {
	id, bs4, err := utils.GetCaptcha()
	if err != nil {
		log.Println("获取验证码失败")
		utils.ErrorWithStatus(c, http.StatusBadRequest, "40001", "获取验证码失败")
		return
	}
	utils.Success(c, gin.H{
		"captchaKey":    id,
		"captchaBase64": bs4,
	}, "获取验证码成功")
}

// 用户登录
func UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	captchaCode := c.PostForm("captchaCode")
	captchaKey := c.PostForm("captchaKey")
	// 验证验证码
	if !utils.VerifyCaptcha(captchaKey, captchaCode) {
		log.Println("验证码错误")
		utils.Error(c, "40001", "验证码错误")
		return
	}
	user := &models.User{Username: username, Password: password}
	err := user.Verify()
	if err != nil {
		log.Println("账号或者密码错误")
		utils.Error(c, "40002", "账号或者密码错误")
		return
	}
	// 生成 token
	token, err := GetToken(username)
	if err != nil {
		log.Println("获取 token 失败", err)
		utils.Error(c, "50001", "获取 token 失败")
		return
	}
	// 登录成功返回 token
	utils.Success(c, gin.H{
		"accessToken":  token,
		"tokenType":    "Bearer",
		"refreshToken": nil,
		"expires":      nil,
	}, "登录成功")
}

func UserOut(c *gin.Context) {
	// 拿到 jwt 中的 username
	if _, Is := c.Get("username"); Is {
		utils.Success(c, nil, "退出成功")
	}
}
