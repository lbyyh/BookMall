package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library-study/app/model"
	"library-study/app/tools"
	"net/http"
)

type Admin struct {
	Name         string `json:"name" form:"name"`
	Password     string `json:"password" form:"password"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}

func AdminLogin(c *gin.Context) {
	var admin Admin
	if err := c.ShouldBind(&admin); err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(), //这里有风险
		})
	}

	fmt.Printf("admin:%+v\n", admin)

	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: admin.CaptchaId,
		Data:      admin.CaptchaValue,
	}) {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10002,
			Message: "验证码校验失败！", //这里有风险
		})
		return
	}

	ret := model.GetAdmin(admin.Name)
	if ret.ID < 1 || ret.Password != tools.EncryptV1(admin.Password) {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: "帐号密码错误！",
		})
		return
	}

	_ = model.SetSession(c, admin.Name, ret.ID)
	c.JSON(http.StatusOK, tools.ECode{
		Message: "登录成功",
	})
}

func AdminLogout(c *gin.Context) {
	_ = model.FlushSession(c)
	c.JSON(http.StatusUnauthorized, tools.ECode{
		Code:    0,
		Message: "您已退出登录",
	})
}
func VisitorLogout(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, tools.ECode{
		Code:    0,
		Message: "您已退出登录",
	})
}

// AdminList 获取管理员列表
func AdminList(c *gin.Context) {
	var admin []model.Admin
	if err := model.MySQLDB.Table("admin").Find(&admin).Error; err != nil {
		fmt.Printf("Retrieve admin error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// 将处理后的数据返回到前端
	c.JSON(http.StatusOK, admin)
}
