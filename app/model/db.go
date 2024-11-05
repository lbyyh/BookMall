package model

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// Globals to hold the database connections
var (
	MySQLDB *gorm.DB
	RedisDB *redis.Client
)

// 假设已经有了一个MongoDB集合的全局变量
var Coll *mongo.Collection

// Initialize both MySQL and Redis connections
func InitializeDatabases(mysqlConfig MysqlConfig, redisConfig RedisConfig) {
	var err error
	MySQLDB, err = newMySQLDB(mysqlConfig)
	if err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}
	RedisDB, store = initRedis(redisConfig)
}

func newMySQLDB(config MysqlConfig) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %s", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB from GORM: %s", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %s", err)
	}
	fmt.Println("Connected to MySQL successfully.")
	return db, nil
}

func initRedis(c RedisConfig) (*redis.Client, *redisstore.RedisStore) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password,
		DB:       0,
	})
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	store, err := redisstore.NewRedisStore(ctx, client)
	if err != nil {
		log.Fatalf("Failed to create redis store: %v", err)
	}
	return client, store
}

// CloseDatabases to disconnect from MySQL and Redis
func CloseDatabases() {
	if MySQLDB != nil {
		sqlDB, err := MySQLDB.DB()
		if err != nil {
			log.Fatalf("Failed to get SQL DB from GORM: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			log.Fatalf("Failed to close MySQL connection: %v", err)
		}
	}
	if RedisDB != nil {
		if err := RedisDB.Close(); err != nil {
			log.Fatalf("Failed to close Redis connection: %v", err)
		}
	}
}
