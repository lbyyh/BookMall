package logic

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"library-study/app/model"
	"net/http"
)

//// 获取用户借书记录
//func GetBookUserList(c *gin.Context) []model.BookUser{
//	var bookUser []model.BookUser
//	if err = model.MySQLDB.Table("book_user").Find(&bookUser).Error;
//	if err != nil {
//		fmt.Println("读取用户借书记录表失败！：%s",errors.Error())
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return surveys
//	}
//}

// GetBookUserList 获取用户借书记录
func GetBookUserList(c *gin.Context) {
	var BookUsers []model.BookUser
	if err := model.MySQLDB.Table("book_user").Find(&BookUsers).Error; err != nil {
		fmt.Printf("Retrieve book_user error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// 将处理后的数据返回到前端
	c.JSON(http.StatusOK, BookUsers)
}

// UpdateBookUser 修改用户借书记录
func UpdateBookUser(c *gin.Context) {
	var bookUser model.BookUser

	// 从请求体中解析数据
	if err := c.ShouldBindJSON(&bookUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建更新的过滤条件和更新数据
	filter := bson.M{"_id": bookUser.ID} // 确保这里使用正确的字段，如果ID是bson.ObjectID，可能需要做转换
	update := bson.M{"$set": bson.M{
		"user_id":      bookUser.UserId, // 确保字段名与数据库一致
		"book_id":      bookUser.BookId,
		"status":       bookUser.Status,
		"time":         bookUser.Time,
		"created_time": bookUser.CreatedTime,
		"updated_time": bookUser.CreatedTime,
	}}

	// 执行更新操作
	result, err := model.Coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Printf("Update book_user error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		// 如果没有文档被修改，可以认为是找不到对应的文档
		c.JSON(http.StatusNotFound, gin.H{"error": "No record updated"})
		return
	}

	// 更新成功，返回更新后的数据
	c.JSON(http.StatusOK, bookUser)
}
