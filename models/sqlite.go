package models

import (
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var isInitialized bool

// InitSqlite 初始化 SQLite 数据库
func InitSqlite() error {
	// 检查并创建 db 目录
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		if err := os.Mkdir("./db", os.ModePerm); err != nil {
			return fmt.Errorf("创建数据库目录失败：%w", err)
		}
	}

	// 连接数据库
	db, err := gorm.Open(sqlite.Open("./db/sublink.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败：%w", err)
	}
	DB = db

	// 检查是否已经初始化
	if isInitialized {
		log.Println("数据库已经初始化，无需重复初始化")
		return nil
	}

	// 自动迁移数据表
	if err := db.AutoMigrate(&User{}, &Subcription{}, &SubLogs{}, &GroupNode{}, &Node{}); err != nil {
		return fmt.Errorf("数据表迁移失败：%w", err)
	}

	// 初始化默认管理员用户
	if err := db.First(&User{}).Error; err == gorm.ErrRecordNotFound {
		admin := &User{
			Username: "admin",
			Password: "123456",
			Role:     "admin",
			Nickname: "管理员",
		}
		if err := admin.Create(); err != nil {
			return fmt.Errorf("初始化管理员用户失败：%w", err)
		}
	}

	isInitialized = true
	log.Println("数据库初始化成功")
	return nil
}
