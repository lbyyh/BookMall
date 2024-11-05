package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// GetUserLoginStatusFromDBOrCache 检索用户的登录状态，
// 假设我们使用Redis来存储和查询状态
func GetUserLoginStatusFromDBOrCache(state string) (*UserLoginStatus, error) {
	var userLoginStatus UserLoginStatus

	// 假设我们有一个Redis客户端实例叫 redisClient
	// 状态值存储在Redis中，key为 "login_status_"+state
	status, err := RedisDB.Get(context.Background(), "login_status_"+state).Result()

	if err == redis.Nil {
		// 缓存中没有找到状态，可能是因为它还没有被设置，或者已经过期
		return nil, fmt.Errorf("no login status found for state: %s", state)
	} else if err != nil {
		// 处理其他可能出现的错误
		return nil, err
	}

	// 假设状态值是一个JSON字符串，我们需要将其反序列化为UserLoginStatus类型
	err = json.Unmarshal([]byte(status), &userLoginStatus)
	if err != nil {
		return nil, err
	}

	return &userLoginStatus, nil
}

// UserLoginStatus 用户登录状态的类型定义，根据实际情况可以增减字段
type UserLoginStatus struct {
	IsConfirmed bool   `json:"is_confirmed"`
	UserID      string `json:"user_id"` // 如果有必要，也可以存储用户的ID
	// ... 可以根据需要增加额外字段，比如登录时间戳等
}

// SetUserLoginConfirmed 更新用户的登录确认状态
func SetUserLoginConfirmed(state string) error {
	// 生成新的状态值
	newUserLoginStatus := UserLoginStatus{
		IsConfirmed: true,
		UserID:      "用户的唯一标识ID", // 实际场景中应当用查询数据库等方法获取到
		// ... 可以添加额外的信息，比如确认登录的时间戳等
	}

	// 将UserLoginStatus序列化为JSON字符串用于存储
	statusJSON, err := json.Marshal(newUserLoginStatus)
	if err != nil {
		return err
	}

	// 假设我们有一个Redis客户端实例叫 redisClient
	// 更新Redis缓存中的状态
	err = RedisDB.Set(context.Background(), "login_status_"+state, statusJSON, 0).Err() // 0表示不设置过期时间
	if err != nil {
		return err
	}

	return nil
}
