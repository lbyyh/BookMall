package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library-study/app/model"
	"library-study/app/tools"
	_ "library-study/docs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// BorrowBook godoc
// @Summary      借书服务
// @Description  用户通过提交图书ID来借阅图书
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        token  header      string  true  "用户验证Token"
// @Param        id     formData    string  true  "图书ID"
// @Success      200    {object}    tools.ECode
// @Failure      400    {object}    tools.ECode
// @Router       /book/borrow [post]
func BorrowBook(c *gin.Context) {
	// 获取用户信息
	cookie, err := c.Request.Cookie("token")
	if err != nil || cookie.Value == "" {
		c.JSON(http.StatusUnauthorized, tools.ECode{
			Code:    10000,
			Message: "无法获取token或token为空",
		})
		return
	}
	jwt := cookie.Value
	userData, err := model.CheckJwt(jwt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, tools.ECode{
			Code:    10001,
			Message: "无效的token",
		})
		return
	}
	uid := userData.Id
	fmt.Printf("-----------------------------")
	fmt.Printf("userData:%v\n", userData)
	// 检查是否超出限流条件
	limited, err := tools.IsRateLimited(strconv.FormatInt(uid, 10), "borrow", 1*time.Hour, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误，请稍候再试"})
		return
	}
	if limited {
		c.JSON(http.StatusTooManyRequests, tools.ECode{
			Code:    10003,
			Message: "您的借书次数过多",
		})
		return
	}

	// 获取图书ID
	idStr := c.PostForm("Id") // 保持表单的名称与 ReturnBook 一致
	fmt.Printf("idstr:%v\n", idStr)
	if idStr == "" {
		c.JSON(http.StatusBadRequest, tools.ECode{
			Code:    10004,
			Message: "缺少图书ID",
		})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id == 0 { // 添加对 id 为 0 的检查
		c.JSON(http.StatusBadRequest, tools.ECode{
			Code:    10005,
			Message: "无效的图书ID",
		})
		return
	}

	// 执行借书逻辑
	err = model.BorrowBook(uid, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, tools.ECode{
			Code:    10006,
			Message: err.Error(), // 添加一个 Detail 字段用于提供详细的错误信息
		})
		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, tools.OK)
}

// ReturnBook godoc
// @Summary      还书服务
// @Description  用户通过提交图书ID来归还图书w
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        token  header      string  true  "用户验证Token"
// @Param        id     formData    string  true  "图书ID"
// @Success      200    {object}    tools.ECode
// @Failure      400    {object}    tools.ECode
// @Router       /book/return [post]
func ReturnBook(c *gin.Context) {
	//获取用户信息
	cookie, err := c.Request.Cookie("token")
	jwt := cookie.Value
	d, err := model.CheckJwt(jwt)
	uid := d.Id
	//获取图书ID
	idStr := c.PostForm("Id")
	if idStr == "" || idStr == "0" {
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	//执行借书逻辑
	err = model.ReturnBook(uid, id)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10002,
			Message: err.Error(),
		})
		return
	}
	//返回成功
	c.JSON(http.StatusOK, tools.OK)
}

func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "获取上传文件失败",
		})
		return
	}

	// 确保avatars目录存在
	avatarDirectory := "app/images/avatars"
	if _, err := os.Stat(avatarDirectory); os.IsNotExist(err) {
		if err := os.MkdirAll(avatarDirectory, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    1,
				"message": "创建目录失败",
			})
			return
		}
	}

	// 创建唯一文件名
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%d-%s", timestamp, filepath.Base(file.Filename))
	filePath := filepath.Join(avatarDirectory, fileName)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    2,
			"message": "保存文件失败",
		})
		return
	}

	// 返回文件路径
	c.JSON(http.StatusOK, gin.H{
		"code":     0,
		"message":  "文件上传成功",
		"filePath": fmt.Sprintf("/images/avatars/%s", fileName),
	})
}

func PicUpload(c *gin.Context) {
	// 获取多个文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "获取上传文件失败",
		})
		return
	}
	files := form.File["avatar"] // 前端通过 `avatar` 字段上传文件

	// 确保avatars目录存在
	avatarDirectory := "app/images/avatars"
	if _, err := os.Stat(avatarDirectory); os.IsNotExist(err) {
		if err := os.MkdirAll(avatarDirectory, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    1,
				"message": "创建目录失败",
			})
			return
		}
	}

	// 循环处理每个文件
	for _, file := range files {
		timestamp := time.Now().UnixNano() // 使用纳秒来确保文件名的唯一性
		fileName := fmt.Sprintf("%d-%s", timestamp, filepath.Base(file.Filename))
		filePath := filepath.Join(avatarDirectory, fileName)

		// 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    2,
				"message": "保存文件失败",
			})
			continue // 处理下一个文件
		}

		// 可以在这里添加逻辑，例如更新数据库记录
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "所有文件上传成功",
	})
}
