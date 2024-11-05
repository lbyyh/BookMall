package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
	"library-study/app/model"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// var state, _ = RandString(5)
var atoken = ""

// Redirect 微信扫码登录
// @Summary 用户登录接口3
// @Description 通过微信扫码登录，手机进行登录验证
// @Tags 公开
// @Accept json
// @Produce application/json
// @Param Url query string true "内网穿透地址"
// @Router /user/wechat/login [get]
func Redirect(c *gin.Context) {
	path := c.Query("Url")
	state, err := RandString(5) // 这应该在每次请求时独立生成
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部错误",
		})
		return
	}
	//防止跨站请求伪造攻击 增加安全性
	redirectURL := url.QueryEscape("http://" + path + "/user/wechat/Callback") //userinfo,
	wechatLoginURL := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&state=%s&scope=snsapi_userinfo#wechat_redirect", "wx55b91f835d27d393", redirectURL, state)

	loginStatusKey := "login_status_" + state
	model.RedisDB.Set(c, loginStatusKey, 0, 3600*time.Second)
	wechatLoginURL, _ = url.QueryUnescape(wechatLoginURL)
	// 生成二维码
	qrCode, err := qrcode.Encode(wechatLoginURL, qrcode.Medium, 256)
	if err != nil {
		// 错误处理
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}
	// 将二维码图片作为响应返回给用户
	c.Header("Content-Type", "image/png")
	c.Header("X-WeChat-State", state)
	fmt.Printf("state------%v\n", state)

	c.Writer.Write(qrCode)
}

type ResponseData struct {
	Data    interface{}
	Message string
	Code    interface{}
}

var user1 model.User

func Callback(c *gin.Context) {
	// 获取微信返回的授权码
	code := c.Query("code")
	// 向微信服务器发送请求，获取access_token和openid
	tokenResp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", "wx55b91f835d27d393", "79e5eb42d91cdaf04dfcd73f5705de16", code))
	CodeServerBusy := "111111"
	if err != nil {
		fmt.Println(err)
		resp := &ResponseData{
			Data:    nil,
			Message: "error,获取token失败",
			Code:    CodeServerBusy,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	// 解析响应中的access_token和openid
	var tokenData struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenID       string `json:"openid"`
		Scope        string `json:"scope"`
	}
	if err1 := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err1 != nil {
		resp := &ResponseData{
			Data:    nil,
			Message: "error,获取token失败",
			Code:    CodeServerBusy,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	userInfoURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", tokenData.AccessToken, tokenData.OpenID)
	userInfoResp, err := http.Get(userInfoURL)
	if err != nil {
		// 错误处理
		fmt.Printf("获取失败--%v\n", err)
		return
	}
	defer userInfoResp.Body.Close()

	//------------------------------------
	var userData struct {
		OpenID   string `json:"openid"`
		Nickname string `json:"nickname"`
		LoggedIn bool
	}
	if err1 := json.NewDecoder(userInfoResp.Body).Decode(&userData); err1 != nil {
		// 错误处理
		fmt.Printf("获取用户信息失败--%v\n", err)
		return
	}

	//用户的名字
	nickname := userData.Nickname // 从userData结构体获取用户昵称

	c.Set("Name", nickname)            // 在上下文c中保存用户的名字，可能用于后续的中间件或处理函数中
	fmt.Printf("nickname%v", nickname) // 打印用户的昵称到控制台，用于调试

	// 在数据库中查找是否有name字段等于nickname的记录
	if err2 := model.MySQLDB.Where("name=?", nickname).First(&user1).Error; err2 != nil {
		// 如果查询结果是记录未找到
		if errors.Is(err2, gorm.ErrRecordNotFound) {
			user1.Name = nickname // 设置user1的Name字段为nickname
			//user1.ID = tools.Snowflake() // 这行代码被注释掉了，可能是之前用来生成ID的方法

			// 在数据库中创建一个新的用户记录
			_, err := model.CreateUser(user1.Name, "")
			if err != nil {
				// 如果创建过程出错，打印错误信息并返回
				fmt.Printf("保存登录信息过程中出错--%v\n", err)
				return
			}
		} else {
			// 如果查询数据库过程中出错（并且错误不是记录未找到）
			fmt.Printf("验证登录信息过程中出错--%v\n", err) // 打印错误信息，注意这里的变量应该是err2
			return                                          // 提前终止函数
		}
	}
	// 添加jwt验证
	if atoken, err := model.GetJwt(user1.ID, user1.Name); err == nil {
		state := c.Query("state")                                                                     // 从请求中提取state参数
		redirectURL := "/wxack?token=" + url.QueryEscape(atoken) + "&state=" + url.QueryEscape(state) // 重定向URL，附加token和state
		c.SetCookie("token", atoken, 3600, "/", "你的域名", false, true)
		c.Redirect(http.StatusFound, redirectURL)
		fmt.Printf("重定向URL：%v\n", redirectURL)
		return
	} else {
		// JWT生成失败的错误处理...
		fmt.Printf("生成JWT令牌失败--%v\n", err)
	}
	// 设置cookie，将token作为cookie传递给前端
	//c.SetCookie("token", atoken, 3600, "/user", "", true, false)

	c.Redirect(http.StatusFound, "/wxack")
	//c.JSON(http.StatusOK, gin.H{
	//	"code":    0,
	//	"message": "登陆成功",
	//})
	//c.Header("Authorization", atoken)
	c.SetCookie("token", atoken, 3600, "/", "", true, false)
	fmt.Printf("登录成功--%v\n", err)
	return
}

func RandString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes), nil
}
func CheckLogin(c *gin.Context) {
	state := c.Query("state")
	// 省略参数检查和错误处理代码

	// 获取login_status值
	loginStatusKey := "login_status_" + state
	loginStatusVal, _ := model.RedisDB.Get(c, loginStatusKey).Result()
	// 省略错误处理代码

	// 对比login_status的值是否为"1"
	if loginStatusVal == "1" {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
		// 注意这里不再调用任何其他处理函数
	} else {
		c.JSON(http.StatusAccepted, gin.H{
			"code": 202,
		})
	}
}

func Updateloginstatus(c *gin.Context) {
	// 使用c.Request.Body获取请求体中的state
	type RequestBody struct {
		State string `json:"state"`
	}
	var requestBody RequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	state := requestBody.State
	// 提取请求头中的Authorization
	authHeader := c.GetHeader("Authorization")
	// 通常，Authorization头的格式为"Bearer <token>"
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}
	// 这里设置了token cookie后，也将token存到Redis，以便于其他接口调用
	if token != "" {
		err := model.RedisDB.Set(c, "token_"+state, token, 3600*time.Second).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "无法保存token",
			})
			return
		}
	}

	//// SetCookie参数根据实际情况调整
	//c.SetCookie("token", token, 3600, "/user", "", true, false)

	// 通过userId获取login_status的键
	loginStatusKey := "login_status_" + state
	// 设置login_status为1
	err := model.RedisDB.Set(c, loginStatusKey, 1, 3600*time.Second).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
		return
	}

	// 这里只需要发送一条响应，因此只保留一个c.JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
