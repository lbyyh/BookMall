package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library-study/app/model"
	"library-study/app/tools"
	"net/http"
)

func CheckUser(c *gin.Context) {
	var name string
	var id int64

	cookie, err := c.Request.Cookie("token")
	if err != nil {
		fmt.Println("--------------------------------")
		c.JSON(http.StatusUnauthorized, tools.ECode{
			Code:    10004,
			Message: "Cookie中token缺失",
		})
		c.Abort()
		return
	}
	fmt.Printf("cookie: %v\n", cookie)

	jwt := cookie.Value
	fmt.Printf("jwt:%v\n", jwt)

	d, err := model.CheckJwt(jwt)
	if err != nil || d == nil { // 确保也检查了`d`是不是nil
		fmt.Printf("error: %v\n", err) // 这里将错误打印出来，以便调试
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "token校验失败",
		})
		c.Abort()
		return // 包含 return，确保在错误情况下函数结束执行
	}

	fmt.Printf("d:%v\n", d)
	name = d.Name
	id = d.Id

	if id <= 0 || name == "" {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "用户信息错误",
		})
		c.Abort()
		return
	}

	if model.GetJWTMap(name) {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "用户状态异常，请重新登录！",
		})
		c.Abort()
		return
	}

	c.Next()
}

// 管理员
func CheckAdmin(c *gin.Context) {
	session := model.GetSession(c)
	fmt.Printf("session:%v\n", session)
	adminID := session["id"]
	adminName := session["name"]

	if adminID == nil || adminName == nil {
		c.JSON(http.StatusUnauthorized, tools.ECode{
			Code:    10004,
			Message: "未授权的访问",
		})
		c.Abort()
		return
	}

	id, ok1 := adminID.(int64)
	name, ok2 := adminName.(string)

	if !ok1 || !ok2 || id <= 0 || name == "" {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "管理员信息错误",
		})
		c.Abort()
		return
	}

	if model.GetJWTMap(name) {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "管理员状态异常，请重新登录！",
		})
		c.Abort()
		return
	}

	c.Next()
}
