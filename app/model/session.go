package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
	"net/http"
	"time"
)

var store *redisstore.RedisStore
var sessionName = "session-name"

func GetSession(c *gin.Context) map[interface{}]interface{} {
	session, _ := store.Get(c.Request, sessionName)
	fmt.Printf("session:%+v\n", session.Values)
	return session.Values
}

func SetSession(c *gin.Context, name string, id int64) error {
	session, err := store.Get(c.Request, sessionName)
	if err != nil {
		// 错误处理: 如日志记录和/或错误传递
		return err
	}
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	session.Values["name"] = name
	session.Values["id"] = id
	return session.Save(c.Request, c.Writer)
}

func FlushSession(c *gin.Context) error {
	session, _ := store.Get(c.Request, sessionName)
	fmt.Printf("session : %+v\n", session.Values)
	session.Flashes()
	return session.Save(c.Request, c.Writer)
}

// GetAdmin 通过用户名查询管理员用户信息。
func GetAdmin(name string) *Admin {
	var ret Admin
	err := MySQLDB.Table("user").Where("name=?", name).First(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
		return &ret
	}
	return &ret
}

//// 该函数中的db是一个全局变量，用于表示数据库连接，它应该在程序初始化时被设置。
//var db *gorm.DB

func SetWithBooks(c *gin.Context, cacheKey string, books []BookInfo) error {
	// 将books序列化为JSON
	booksJson, err := json.Marshal(books)
	if err != nil {
		// 如果序列化失败，记录到日志并返回错误
		fmt.Printf("error marshalling books to JSON: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error marshalling books"})
		return err
	}

	// 将序列化后的JSON字符串存储到Redis中
	err = RedisDB.Set(context.Background(), cacheKey, booksJson, 5*time.Minute).Err()
	if err != nil {
		// 如果存储失败，记录到日志并返回错误
		fmt.Printf("error setting books cache in redis: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to set books cache"})
		return err
	}
	return nil
}
