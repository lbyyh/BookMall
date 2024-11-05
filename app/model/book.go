package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func CreateBook(b *BookInfo) error {
	return MySQLDB.Create(b).Error
}

func GetBook(id int64) (*BookInfo, error) {
	var ret BookInfo
	MySQLDB.Where("id = ?", id).First(&ret)
	return &ret, nil
}

func SaveBook(data *BookInfo) error {
	return MySQLDB.Save(data).Error
}

func DeleteBook(id int64) error {
	return MySQLDB.Where("id = ?", id).Delete(&BookInfo{}).Error
}

func BorrowBook(uid, id int64) error {
	tx := MySQLDB.Begin()
	//查询用户是否存在
	var user User
	tx.Where("id = ?", uid).First(&user)
	if user.ID == 0 {
		tx.Rollback()
		return errors.New("用户信息不存在")
	}

	//查询图书是否存在，是否正常
	var book BookInfo
	tx.Where("id = ?", id).First(&book)
	if book.Id == 0 || book.Num <= 0 {
		tx.Rollback()
		return errors.New("图书信息不存在或库存不足")
	}
	//创建借阅记录
	now := time.Now()
	bu := BookUser{
		UserId:      uid,
		BookId:      id,
		Status:      1,
		Time:        1,
		CreatedTime: now,
		UpdatedTime: now,
	}
	if tx.Create(&bu).Error != nil {
		tx.Rollback()
		return errors.New("创建一个借阅记录")
	}
	//扣减图书库存
	book.Num = book.Num - 1
	if tx.Save(&book).Error != nil {
		tx.Rollback()
		return errors.New("扣减图书库存")
	}

	tx.Commit()
	return nil
}

func ReturnBook(uid, id int64) error {
	tx := MySQLDB.Begin()
	//查询用户是否存在
	var user User
	tx.Where("id = ?", uid).First(&user)
	if user.ID == 0 {
		tx.Rollback()
		return errors.New("用户信息不存在")
	}

	//查询图书是否存在，是否正常
	var book BookInfo
	tx.Where("id = ?", id).First(&book)
	if book.Id == 0 {
		tx.Rollback()
		return errors.New("图书信息不存在")
	}

	//查询借书记录是否存在
	var bu BookUser
	tx.Where("user_id = ? and book_id = ? and status = ?", uid, id, 1).First(&bu)
	if bu.ID <= 0 {
		tx.Rollback()
		return errors.New("借阅记录不存在")
	}

	//更新借阅状态
	bu.Status = 0
	if err := tx.Save(&bu).Error; err != nil {
		tx.Rollback()
		return errors.New(fmt.Sprintf("修改借阅记录失败：%s", err.Error()))
	}

	//更新图书库存
	book.Num = book.Num + 1
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		return errors.New(fmt.Sprintf("增加库存失败：%s", err.Error()))
	}
	tx.Commit()
	return nil
}

func GetBooks(limit, offset int) ([]BookInfo, error) {
	var ret []BookInfo
	err := MySQLDB.Offset(offset).Limit(limit).Table("book_info").Find(&ret).Error
	if err != nil {
		fmt.Printf("查询书籍列表失败: %s", err.Error())
		return nil, err
	}
	return ret, nil
}

func GetBooksBorrowingRecord() ([]BookUser, error) {
	var ret []BookUser
	err := MySQLDB.Table("book_user").Find(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret, err
}

// 获取BookInfo，缓存totalPages
func GetPaginatedBooksData(page int, perPage int) ([]BookInfo, error) {
	var books []BookInfo
	// 通过页码和每页数量构建缓存键
	cacheKey := fmt.Sprintf("totalpage")
	temp, err := RedisDB.Get(context.TODO(), cacheKey).Result()
	totalPages, _ := strconv.ParseInt(temp, 10, 64)
	if err == redis.Nil {
		MySQLDB.Model(&BookInfo{}).Count(&totalPages)
	}

	// 缓存总页数信息
	RedisDB.Set(context.Background(), cacheKey, strconv.Itoa(int(totalPages)), 5*time.Minute)

	offset := (page - 1) * perPage

	if err := MySQLDB.Offset(offset).Limit(perPage).Find(&books).Error; err != nil {
		return nil, err
	}
	fmt.Printf("-------------books:%v\n", books)
	return books, nil
}

//// 获取书籍总数
//func GetTotalBooks() (int, error) {
//	var totalBooks int64
//	// 假设 `db` 是预先配置好的数据库连接
//	err := MySQLDB.Table("book_info").Count(&totalBooks).Error
//	if err != nil {
//		// 处理错误: 这里可以自定义错误或者直接返回错误依情况而定
//		return 0, err
//	}
//	return int(totalBooks), nil
//}
