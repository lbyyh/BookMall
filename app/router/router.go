package router

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"library-study/app/logic"
	"library-study/app/middleware"
	"library-study/app/model"
	"library-study/app/tools"
	_ "library-study/docs"
	"net/http"
	"time"
)

func New() {
	r := gin.Default()

	// 配置CORS策略
	config := cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8087"}, // 仅允许这个来源的跨源请求
		AllowMethods:     []string{"GET", "POST"},           // 允许GET和POST请求方法
		AllowHeaders:     []string{"Origin"},                // 允许带这些请求头
		ExposeHeaders:    []string{"Content-Length"},        // 允许浏览器访问这些响应头
		AllowCredentials: true,                              // 因为涉及支付，可能需要cookies或授权验证
		MaxAge:           12 * time.Hour,                    // 预检请求的缓存持续时间
	}
	r.Use(cors.New(config))

	r.LoadHTMLGlob("app/view/*")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//初始登陆界面
	index := r.Group("")
	{
		index.GET("/index", logic.Index)                                       //静态页面
		index.GET("/wxack", logic.Wxack)                                       //静态页面
		index.GET("/BorrowingRecord", logic.BorrowingRecord)                   //静态页面
		index.GET("/RBACPermissionManagement", logic.RBACPermissionManagement) //静态页面
		index.GET("/BookList", logic.BookList)                                 //图书列表页面
		index.GET("/BookBS", logic.BookBS)                                     //图书列表页面
		index.GET("/PictureUpload", logic.PictureUpload)                       //图片上传页面
	}
	r.Static("/images", "app/images")
	// 上传接口
	r.POST("/user/uploadAvatar", logic.UploadAvatar)
	r.StaticFile("/favicon.ico", "./app/images/avatars/图书(1).png")

	{ //用户user
		user := r.Group("/user")
		user.POST("/login", logic.Login)
		user.POST("/wxlogin", logic.Wxlogin)
		user.POST("/update-login-status", logic.Updateloginstatus)
		user.GET("/GetToken", model.GetToken)
		user.POST("/SendEmailCaptcha", logic.SendEmailCaptcha)
		user.POST("/VerifyEmailCaptcha", logic.VerifyEmailCaptcha)
		user.POST("/SendSMSCaptcha", logic.SendSMSCaptcha)
		user.POST("/VerifySMSCaptcha", logic.VerifySMSCaptcha)
		user.GET("/wechat", logic.CheckSignature)
		user.GET("/wechat/login", logic.Redirect)
		user.GET("/wechat/Callback", logic.Callback)
		user.GET("/wechat/check_login", logic.CheckLogin)
		user.GET("/UserList", model.UserList) //管理界面

		user.Use(middleware.CheckUser)
		user.POST("/book_info/borrow", logic.BorrowBook)
		user.POST("/book_info/return", logic.ReturnBook)
		user.GET("/UserLogin", logic.UserLogin)  //用户界面
		user.POST("/PicUpload", logic.PicUpload) //用户界面
	}

	{ //管理员admin
		admin := r.Group("/admin")
		admin.POST("/login", logic.AdminLogin)
		admin.Use(middleware.CheckAdmin)
		admin.GET("/logout", logic.AdminLogout)
		admin.GET("/AdminLogin", logic.AdminLoginS) //管理员界面
		admin.GET("/AdminList", logic.AdminList)    //管理员界面

	}

	{ //游客visitor
		visitor := r.Group("/visitor")
		visitor.GET("/VisitorLogout", logic.VisitorLogout)
		visitor.GET("/VisitorLogin", logic.VisitorLoginS) //游客界面
	}
	{ //书籍
		book := r.Group("/book_info")
		book.GET("/GetBook", logic.GetBook)
		book.GET("/list", logic.GetPaginatedBooks)
		book.GET("/BooksBorrowingRecord", logic.BooksBorrowingRecord)
		book.POST("/AddBook", logic.AddBook)
		book.PUT("/SaveBook", logic.SaveBook)
		book.DELETE("/DelBook", logic.DelBook)
		book.GET("/GetPaginatedBooks", logic.GetPaginatedBooks)
	}
	//借阅模块
	{
		bookUser := r.Group("/book_user")
		bookUser.GET("/GetBookUserList", logic.GetBookUserList)
		bookUser.POST("/UpdateBookUser", logic.UpdateBookUser)
	}
	//权限模块
	{
		rights := r.Group("/rights")
		rights.GET("/GetRoles", logic.GetRoles)
		rights.GET("/GetPermissions", logic.GetPermissions)
		rights.GET("/GetRole_Pre", logic.GetRolePre)
		rights.GET("/GetUer_Roles", logic.GetUerRoles)
		rights.GET("/Role_Pre", logic.GetRolePre)
		rights.PUT("/UpdateRole_Pre", logic.UpdateRolePre)
		rights.POST("/UpdateUer_Roles", logic.AddRolePermissions)
	}
	{ //支付模块
		alipay := r.Group("/alipay")
		alipay.POST("/pay", logic.PayHandler)
		alipay.GET("/callback", logic.CallbackHandler)
		alipay.POST("/notify", logic.NotifyHandler)
	}
	{
		//验证码
		r.GET("/captcha", func(context *gin.Context) {

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
		})

		r.POST("/captcha/verify", func(context *gin.Context) {
			var param tools.CaptchaData
			if err := context.ShouldBind(&param); err != nil {
				context.JSON(http.StatusOK, tools.ParamErr)
				return
			}

			fmt.Printf("参数为：%+v", param)
			if !tools.CaptchaVerify(param) {
				context.JSON(http.StatusOK, tools.ECode{
					Code:    10008,
					Message: "验证失败",
				})
				return
			}
			context.JSON(http.StatusOK, tools.OK)
		})
	}

	if err := r.Run(":8087"); err != nil {
		panic(err)
	}
}
