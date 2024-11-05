package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"library-study/app/model"
	"library-study/app/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./app")
	// 读取配置数据
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func Start() {
	InitConfig()

	// 从配置文件中读取数据库配置
	redisConf := model.RedisConfig{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetInt("redis.port"),
		Password: viper.GetString("redis.password"),
	}
	mysqlConf := model.MysqlConfig{
		Host:     viper.GetString("mysql.host"),
		Port:     viper.GetInt("mysql.port"),
		Username: viper.GetString("mysql.username"),
		Password: viper.GetString("mysql.password"),
		Database: viper.GetString("mysql.database"),
	}

	// 初始化数据库连接
	model.InitializeDatabases(mysqlConf, redisConf)
	defer model.CloseDatabases() // 确保在程序退出时关闭数据库连接

	// 初始化Gin引擎
	r := gin.Default()

	// 初始化路由
	router.New()

	// 设置 tools 包中的 Redis 客户端
	//tools.SetRedisClient(model.RedisDB)

	// 创建自定义的 HTTP 服务器
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 启动服务器并在新的 goroutine 中运行
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	fmt.Println("Service started successfully.")

	// 优雅停止服务的逻辑
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞直至接收到终止信号
	fmt.Println("Shutting down server...")

	// 创建一个超时时间为5秒的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
