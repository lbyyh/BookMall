package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 向数据库获取用户的函数
func GetUser(name string) *User {
	var ret User
	err := MySQLDB.Table("user").Where("name=?", name).First(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
		return &ret
	}
	return &ret
}

// 向数据库添加新用户的函数
func CreateUser(name, password string) (*User, error) {
	// 创建User实例
	newUser := User{
		Name:        name,
		Password:    password,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	// 插入新创建的用户记录到数据库
	if err := MySQLDB.Table("user").Create(&newUser).Error; err != nil {
		fmt.Printf("Create user error: %s", err.Error())
		return nil, err
	}
	return &newUser, nil
}

// UserList 获取管理员列表
func UserList(c *gin.Context) {
	var user []User
	if err := MySQLDB.Table("user").Find(&user).Error; err != nil {
		fmt.Printf("Retrieve user error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// 将处理后的数据返回到前端
	c.JSON(http.StatusOK, user)
}
