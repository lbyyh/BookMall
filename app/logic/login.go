package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library-study/app/model"
	"library-study/app/tools"
	_ "library-study/docs"
	"net/http"
)

type User struct {
	Name         string `json:"name" form:"name"`
	Password     string `json:"password" form:"password"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}

// Login godoc
// @Summary      执行用户登录
// @Description  执行用户登录
// @Tags         login
// @Accept       json
// @Produce      json
// @Param        name   body      User true	"login User"
// @Success      200  {object}  tools.ECode
// @Router       /login [post]
func Login(c *gin.Context) {
	fmt.Println("-------------------------")
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: "输入参数有误", // 更改为通用错误消息，避免敏感信息泄露
		})
		c.Abort() // 终止请求处理
		return
	}

	fmt.Printf("user:%+v\n", user)

	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: user.CaptchaId,
		Data:      user.CaptchaValue,
	}) {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10002,
			Message: "验证码校验失败", // 更改为通用错误消息，避免敏感信息泄露
		})
		c.Abort() // 终止请求处理
		return
	}

	ret := model.GetUser(user.Name)
	//fmt.Printf("tools.EncryptV1(user.Password):%v\n", tools.EncryptV1(user.Password))
	if ret.ID < 1 || ret.Password != tools.EncryptV1(user.Password) {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: "帐号密码错误",
		})
		c.Abort() // 终止请求处理
		return
	}

	// 生成TOKEN
	token, err := model.GetJwt(ret.ID, user.Name)
	c.SetCookie("token", token, 3600, "/user", "", true, false)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "登录失败，无法生成token",
		})
		c.Abort() // 终止请求处理
		return
	}

	// 将token发送给客户端
	c.JSON(http.StatusOK, tools.ECode{
		Code:    0, // 一般来说成功的响应代码是 0
		Message: "登录成功",
		Data:    token,
	})

	// 不需要再次调用 c.JSON，因为token已经在上面发送过了
	// c.JSON(http.StatusOK, gin.H{
	// 	"token": token,
	// })

	return
}

// Wxlogin godoc
// @Summary      执行用户登录
// @Description  执行用户登录
// @Tags         login
// @Accept       json
// @Produce      json
// @Param        name   body      User true	"login User"
// @Success      200  {object}  tools.ECode
// @Router       /login [post]
func Wxlogin(c *gin.Context) {
	//var user User
	//
	//if err := c.ShouldBind(&user); err != nil {
	//	c.JSON(http.StatusOK, tools.ECode{
	//		Code:    10001,
	//		Message: "输入参数有误", // 更改为通用错误消息，避免敏感信息泄露
	//	})
	//	c.Abort() // 终止请求处理
	//	return
	//}

	//// 从上下文中取出我们之前存储的nickname
	//nickname, _ := c.Get("Name")
	//
	//user.Name = nickname.(string)
	//fmt.Printf("-------------------------0--------------c:%+v\n", c)
	//fmt.Printf("-------------------------1--------------user:%+v\n", user)

	//if !tools.CaptchaVerify(tools.CaptchaData{
	//	CaptchaId: user.CaptchaId,
	//	Data:      user.CaptchaValue,
	//}) {
	//	c.JSON(http.StatusOK, tools.ECode{
	//		Code:    10002,
	//		Message: "验证码校验失败", // 更改为通用错误消息，避免敏感信息泄露
	//	})
	//	c.Abort() // 终止请求处理
	//	return
	//}

	//ret := model.GetUser(user.Name)
	//fmt.Printf("-------------------------1--------------ret:%+v\n", ret)
	//if ret.ID < 1 || ret.Password != tools.EncryptV1(user.Password) {
	//	c.JSON(http.StatusOK, tools.ECode{
	//		Code:    10001,
	//		Message: "帐号密码错误",
	//	})
	//	c.Abort() // 终止请求处理
	//	return
	//}

	// 生成TOKEN
	//token, err := model.GetJwt(ret.ID, user.Name)
	//c.SetCookie("token", token, 3600, "/user", "", true, false)
	//if err != nil {
	//	c.JSON(http.StatusOK, tools.ECode{
	//		Code:    10003,
	//		Message: "登录失败，无法生成token",
	//	})
	//	c.Abort() // 终止请求处理
	//	return
	//}
	//c.SetCookie("token", atoken, 3600, "/user", "", true, false)
	//c.SetCookie("token", atoken, 3600, "/user", "", true, false)
	// 将token发送给客户端
	c.JSON(http.StatusOK, tools.ECode{
		Code:    0, // 一般来说成功的响应代码是 0
		Message: "登录成功",
		Data:    atoken,
	})

	// 不需要再次调用 c.JSON，因为token已经在上面发送过了
	// c.JSON(http.StatusOK, gin.H{
	// 	"token": token,
	// })

	return
}

func Logout(c *gin.Context) {

}

func GetCaptcha(context *gin.Context) {
	captcha, err := tools.CaptchaGenerate()
	if err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, tools.ECode{
		Data: captcha,
	})
}
