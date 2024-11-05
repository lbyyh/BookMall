package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"library-study/app/model"
	"library-study/app/tools"
	_ "library-study/docs"
	"net/http"
	"strconv"
)

type Book struct {
	ID         int
	Title      string
	Author     string
	IsBorrowed bool
}

// GetBook godoc
// @Summary      获取图书信息
// @Description  通过图书ID获取单个图书的详细信息
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        id  query  string  true  "图书ID"
// @Success      200  {object}  tools.ECode
// @Failure      400  {object}  tools.ECode
// @Router       /book/get [get]
func GetBook(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" || idStr == "0" {
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	ret, err := model.GetBook(id)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tools.ECode{
		Data: ret,
	})
	return

}

// GetBooks godoc
// @Summary      获取所有图书
// @Description  获取所有图书的列表
// @Tags         book
// @Accept       json
// @Produce      json
// @Success      200  {object}  tools.ECode
// @Router       /books/get [get]
func GetBooks(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid limit"})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid offset"})
		return
	}

	// Redis缓存的键，包括页码和每页数量
	// 尝试从Redis获取缓存的图书列表
	var books []model.BookInfo
	cacheKey := fmt.Sprintf("book_page_%s_%s", offset, limit)
	val, err := model.RedisDB.Get(context.Background(), cacheKey).Result()

	if err == nil {
		// 缓存中有数据，解析JSON并返回
		err = json.Unmarshal([]byte(val), &books)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"data": books})
			return
		}
		// 如果缓存的数据解析失败，继续向下执行以查询数据库
	}

	// 缓存中没有数据或者无法解析，需要查询数据库
	books, err = model.GetBooks(limitInt, offsetInt)
	// 这里假设 GetBooksFromDB 是在 model 包中定义的，用于获取书籍列表
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	// 将书籍列表信息序列化并存入Redis缓存
	if err := model.SetWithBooks(c, cacheKey, books); err != nil {
		// 失败处理已经在SetSessionWithBooks函数内部完成，这里不需再次处理
		return
	}

	// 返回数据库查询结果
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// AddBook godoc
// @Summary      添加新图书
// @Description  添加一个新的图书记录到库存
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book  body  model.Book  true  "图书信息"
// @Success      200  {object}  tools.ECode
// @Failure      400  {object}  tools.ECode
// @Router       /book/add [post]
func AddBook(c *gin.Context) {
	var data model.BookInfo
	if err := c.ShouldBind(&data); err != nil {
		fmt.Printf("data:%v\n", data)
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	data.Uid = tools.Snowflake()
	//TODO:增加参数校验
	//id, _ := strconv.ParseInt(idStr, 10, 64)
	fmt.Printf("data:%v\n", data)
	err := model.CreateBook(&data)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tools.OK)
	return
}

// DelBook godoc
// @Summary      删除图书
// @Description  通过图书ID删除一个图书记录
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        id  query  string  true  "图书ID"
// @Success      200  {object}  tools.ECode
// @Router       /book/delete [delete]
func DelBook(c *gin.Context) {
	idStr := c.Query("id")
	fmt.Printf("id:%v\n", idStr)
	if idStr == "" || idStr == "0" {
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	err := model.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tools.OK)
	return
}

// SaveBook godoc
// @Summary      保存图书信息
// @Description  保存或更新一个图书记录的信息
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        id  query  string  false  "图书ID"  // 此处根据是否为新增或更新决定id是否必须
// @Param        book  body  model.Book  true  "图书信息"
// @Success      200  {object}  tools.ECode
// @Router       /book/save [put]
func SaveBook(c *gin.Context) {
	var data model.BookInfo
	if err := c.ShouldBind(&data); err != nil {
		fmt.Printf("data:%v\n", data)
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	data.Uid = tools.Snowflake()
	//TODO:增加参数校验
	err := model.SaveBook(&data)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tools.ECode{
		Code: 0,
	})
	return
}

// BooksBorrowingRecord godoc
// @Summary      借书记录
// @Description  获取所有当前借出的图书记录
// @Tags         book
// @Accept       json
// @Produce      json
// @Success      200  {object}  tools.ECode
// @Router       /books/borrowing-record [get]
func BooksBorrowingRecord(c *gin.Context) {
	ret, err := model.GetBooksBorrowingRecord()
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tools.ECode{
		Code: 0,
		Data: ret,
	})
	return
}

// GetPaginatedBooks godoc
// @Summary      获取分页图书列表
// @Description  获取分页后的图书列表和总页数
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        page       query     string  false  "页码"
// @Param        perPage    query     string  false  "每页显示数量"
// @Success      200        {object}  tools.ECode
// @Failure      400        {object}  tools.ECode
// @Router       /book_info/list [get]
func GetPaginatedBooks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")        // 默认为第一页
	perPageStr := c.DefaultQuery("perPage", "10") // 默认每页显示10条数据

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
		return
	}
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid perPage parameter"})
		return
	}

	// 通过页码和每页数量构建缓存键
	cacheKey := fmt.Sprintf("book_page_%d_%d", page, perPage)

	var books []model.BookInfo
	var totalPages int

	// 尝试从Redis获取缓存的图书列表
	val, err := model.RedisDB.Get(context.Background(), cacheKey).Result()
	if err == nil {
		// 缓存中有数据，解析JSON并返回
		err = json.Unmarshal([]byte(val), &books)
		fmt.Println("-----------------------------?----------")

	} else {
		// 缓存中没有数据或解析失败，需要查询数据库
		books, err = model.GetPaginatedBooksData(page, perPage)
		for i := range books {
			books[i].ImgUrl = "/images/" + books[i].ImgUrl
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching paginated books"})
			return
		}
	}

	//// 解析总页数（如果缓存了的话）
	totalPagesVal, totalPagesErr := model.RedisDB.Get(context.Background(), "totalpage").Result()
	if totalPagesErr == nil {
		totalPages, _ = strconv.Atoi(totalPagesVal)
	}

	// 将书籍列表信息序列化并存入Redis缓存
	json.Marshal(books)
	if err := model.SetWithBooks(c, cacheKey, books); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set books cache"})
		return
	}

	c.JSON(http.StatusOK, tools.ECode{
		Code:    0,
		Message: "Success",
		Data: gin.H{
			"books":      books,
			"totalPages": totalPages,
		},
	})
}
