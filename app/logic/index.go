package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "library-study/docs"
	"net/http"
)

// Index godoc
// @Summary      主页服务
// @Description  显示应用程序主页
// @Tags         general
// @Accept       html
// @Produce      html
// @Success      200  {string}  string  "成功渲染主页"
// @Router       / [get]
func Index(context *gin.Context) {
	//fmt.Println("--------------------------------------")
	context.HTML(http.StatusOK, "index.tmpl", nil) //http.statusOK == 200
	fmt.Println("主页")
}

func Wxack(context *gin.Context) {
	//fmt.Println("--------------------------------------")
	context.HTML(http.StatusOK, "wxack.tmpl", nil) //http.statusOK == 200
	fmt.Println("")
}

func BorrowingRecord(context *gin.Context) {
	//fmt.Println("--------------------------------------")
	context.HTML(http.StatusOK, "BorrowingRecord.tmpl", nil) //http.statusOK == 200
	fmt.Println("")
}

// UserLogin godoc
// @Summary      执行用户登录
// @Description  执行用户登录，返回JWT令牌
// @Tags         login
// @Accept       json
// @Produce      json
// @Param        username   formData      string  true  "用户名"
// @Param        password   formData      string  true  "密码"
// @Success      200  {object}  tools.ECode
// @Failure      400  {object}  tools.ECode
// @Failure      500  {string} string       "Internal Server Error"
// @Router       /login [post]
func UserLogin(context *gin.Context) {
	//fmt.Println("--------------------------------------")
	context.HTML(http.StatusOK, "UserLogin.tmpl", nil) //http.statusOK == 200
	fmt.Println("用户主页")
}

// AdminLoginS godoc
// @Summary      管理员登录服务
// @Description  进行管理员登录并展示管理员首页
// @Tags         login
// @Accept       html
// @Produce      html
// @Success      200  {string}  string  "成功展示管理员登录页面"
// @Router       /admin/login [get]
func AdminLoginS(context *gin.Context) {
	//fmt.Println("--------------------------------------")
	context.HTML(http.StatusOK, "AdminLogin.tmpl", nil) //http.statusOK == 200
	fmt.Println("管理员主页")
}

// VisitorLoginS godoc
// @Summary      游客登录服务
// @Description  进行游客登录并展示游客首页
// @Tags         login
// @Accept       html
// @Produce      html
// @Success      200  {string}  string  "成功展示游客登录页面"
// @Router       /images/login [get]
func VisitorLoginS(context *gin.Context) {
	//fmt.Println("--------------------------------------")
	context.HTML(http.StatusOK, "VisitorLogin.tmpl", nil) //http.statusOK == 200
	fmt.Println("游客主页")
}

func RBACPermissionManagement(context *gin.Context) {
	//fmt.Println("--------------------------------------")
	context.HTML(http.StatusOK, "RBACPermissionManagement.tmpl", nil) //http.statusOK == 200
	fmt.Println("游客主页")
}

func BookList(context *gin.Context) {
	//fmt.Println("--------------------------------------")
	context.HTML(http.StatusOK, "bookList.tmpl", nil) //http.statusOK == 200
	fmt.Println("图书列表界面")
}

func BookBS(context *gin.Context) {
	//fmt.Println("--------------------------------------")
	context.HTML(http.StatusOK, "booksBS.tmpl", nil) //http.statusOK == 200
	fmt.Println("买书卖书界面")
}
func PictureUpload(context *gin.Context) {
	//fmt.Println("--------------------------------------")
	context.HTML(http.StatusOK, "PictureUpload.tmpl", nil) //http.statusOK == 200
	fmt.Println("图片上传界面")
}
