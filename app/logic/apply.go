package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"
	"log"
	"net/http"
)

var client *alipay.Client

const (
	kAppId = "9021000137656016"
	//私钥
	kPrivateKey = "MIIEpAIBAAKCAQEAnuavYwyatBBHrcsT+ADTZt+l1Ohmwj/u03UWgGye66xpElfZC+DPN+GFbHAecnQYwnDVm104JQF6kKq5xms5eUX2LwmnxWmTN3F7o+TIFeqbX5Jw+bBhT3S4gbU7EAQd2ca3IFXs6oeQ8EZ2dlgDkwVlHSlbFgNKb/6f2sJpGsAGGykaguZKwuqjYv09qC8fzMiEzZzvm5oT0UlhK6v0YZ2EsRjMEPg0WiwtzoEo9BModzT0z0OeW33NhOD706MRpq3WcpQscHCJqKOVmJQCnA3NCePgdxxtytGzFbj/74gQ/9X8w9K2UV1Yb9BcM0BgZ4PTw+H1522BQu58RKIF+wIDAQABAoIBABk+tSabXgi1fW3TEb0ZBH0XkxUcRxcdaSgXNhf5KdZvcdIEOut0L/fE0JnFxlCQuU5K9uTUDpNyhLJvLUykxGDMCKy4b/shJs5sLSSAuHki6MRqU6CXsR4agSW6UUPeI4/xzi5I+HbaSuChkTiECy1UchgL9fitVSot3d+3e1NQZn/AyD5gsBAUfL/lehLhKThLwEH0il6FxNJWX+7boiGqlD4LkR6uy/KpCi92o39Lc+3SICubIA+QpopdiipxskMcqIZfeOSCz9nXWdP4se+fW0hFl5FfsVufHmy3abjQEE5qBQ8bcuZ3bSwYTbbOkaVjCirYXvqTu/44HNhIgGECgYEA01zvzJNB4G/NDJCTccxkRLIOrMTsLIVpH7nreIDWq65Qqltnmtc6DXKHTsrg1rbkaXBHtInTSl0zqmpch5tCVKaYf9Vc98hso481SKcNMrva6S193KATIMSZrj+asUkydhe9LTebEhGg3+2xnbFTEzjHBvx59UOZFVEKdaN3dGcCgYEAwHV2xaV3BFGgeFMYnFkCIAK0hPx5xwssP7fj//IMPHkotlLxbyEoZsuXLhLZldrcQFLNCtNfUnOzHv9F5D0zSLuo3o4R8ygrhjGk5FFam1mgegwslWuojrzW5ifmFy0G/1uAFtKGEvEZUfe1RUQoTJfz+/CNO0gBG1mMBwNiBU0CgYBElBL6PY1SVPQi74XnlnmyEFPSmtJGX8MMGDbekm8UpSpnG+ExzEN5uX9NgWYSRKU30MZzPYTgy/zHflsnZKjQ7nzsfT3853rYVs7jE9CkdW9B2RDNVOLf7uouL1Tx0N4ekvU+hpw58J5SCb1nfPGHexSYn7KycYxp7jGGmdNYYQKBgQCNckc2f6N3Mx7DEB9YWTpsmFBgJMbDePyuX9Jb+2Lu1wUK6u6yhCYTVrHnlMcBkfap97Dmse6uxIXy1B5j3m7gl7tGxhd/JBjI6ZeMjhYPctG0oVnq/1LEhRlT0iMTCW7JIlCDdXpAVZ4MVgeNvsf3cv5IPcUuun7FwQxe4yeZSQKBgQDNk76NYtBgkeBIh7W19rjr7miIt1/62FbNC9kLMn147001RHuO9Dsgd1y6alEyGvbXLIxu9yjHAa+kQIz3M/0b8myuFcNgOsK0vnM7LXz7/e7cfCUlF7wf6mEdIeSBFKJy0J/9fJXC6NKN9eKidzUOhxOlxw+OIWenq5BJrpiGEQ=="
	// TODO 设置回调地址域名
	kServerDomain = ""
)

func init() {
	// 在程序启动初始化支付宝客户端
	var err error
	client, err = InitAliPayClient()
	if err != nil {
		log.Fatalf("初始化支付宝客户端失败：%v", err)
	}
}

// PayHandler 处理支付请求
func PayHandler(c *gin.Context) {
	// ...原有的pay逻辑...
	var tradeNo = fmt.Sprintf("%d", xid.Next())
	totalAmount := c.PostForm("totalAmount")

	var p = alipay.TradePagePay{}
	p.NotifyURL = kServerDomain + "/alipay/notify"
	p.ReturnURL = kServerDomain + "/alipay/callback"
	p.Subject = "支付测试:" + tradeNo
	p.OutTradeNo = tradeNo
	p.TotalAmount = totalAmount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	// 根据支付请求参数生成支付页URL
	url, err := client.TradePagePay(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成支付页URL失败"})
		return
	}

	// 前端需要重定向到这个URL完成支付
	c.JSON(http.StatusOK, gin.H{"url": url.String()})
}

// CallbackHandler 处理支付宝回调
func CallbackHandler(c *gin.Context) {
	// 解析请求参数
	if err := c.Request.ParseForm(); err != nil {
		log.Println("解析请求参数发生错误", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析请求参数发生错误"})
		return
	}

	// 获取通知参数
	outTradeNo := c.Request.Form.Get("out_trade_no")

	// 验证签名
	if err := client.VerifySign(c.Request.Form); err != nil {
		log.Println("回调验证签名发生错误", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "回调验证签名发生错误"})
		return
	}

	log.Println("回调验证签名通过")

	// 查询订单状态
	var p = alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo

	rsp, err := client.TradeQuery(c, p)
	if err != nil {
		log.Printf("验证订单 %s 信息发生错误: %s", outTradeNo, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("验证订单 %s 信息发生错误: %s", outTradeNo, err.Error())})
		return
	}

	if rsp.IsFailure() {
		log.Printf("验证订单 %s 信息发生错误: %s-%s", outTradeNo, rsp.Msg, rsp.SubMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("验证订单 %s 信息发生错误: %s-%s", outTradeNo, rsp.Msg, rsp.SubMsg)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("订单 %s 支付成功", outTradeNo)})
}

// NotifyHandler 处理支付宝通知
func NotifyHandler(c *gin.Context) {
	// 解析请求参数
	if err := c.Request.ParseForm(); err != nil {
		log.Println("解析请求参数发生错误", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析请求参数发生错误"})
		return
	}

	// 解析异步通知
	notification, err := client.DecodeNotification(c.Request.Form)
	if err != nil {
		log.Println("解析异步通知发生错误", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析异步通知发生错误"})
		return
	}

	log.Println("解析异步通知成功:", notification.NotifyId)

	// 查询订单状态
	var p = alipay.NewPayload("alipay.trade.query")
	p.AddBizField("out_trade_no", notification.OutTradeNo)

	var rsp *alipay.TradeQueryRsp
	if err := client.Request(c, p, &rsp); err != nil {
		log.Printf("异步通知验证订单 %s 信息发生错误: %s", notification.OutTradeNo, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("异步通知验证订单 %s 信息发生错误: %s", notification.OutTradeNo, err.Error())})
		return
	}

	if rsp.IsFailure() {
		log.Printf("异步通知验证订单 %s 信息发生错误: %s-%s", notification.OutTradeNo, rsp.Msg, rsp.SubMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("异步通知验证订单 %s 信息发生错误: %s-%s", notification.OutTradeNo, rsp.Msg, rsp.SubMsg)})
		return
	}

	log.Printf("订单 %s 支付成功", notification.OutTradeNo)

	client.ACKNotification(c.Writer)
}

// InitAliPayClient 初始化支付宝客户端并加载证书
func InitAliPayClient() (*alipay.Client, error) {
	c, err := alipay.New(kAppId, kPrivateKey, false) // 使用 false 表示在沙箱环境中
	if err != nil {
		return nil, fmt.Errorf("创建支付宝客户端失败: %w", err)
	}

	// 加载应用公钥证书
	if err = c.LoadAppCertPublicKeyFromFile("app/certificate/appPublicCert.crt"); err != nil {
		log.Println("加载证书发生错误", err)
		return nil, fmt.Errorf("加载应用公钥证书失败: %w", err)
	}
	// 加载支付宝根证书
	if err = c.LoadAliPayRootCertFromFile("app/certificate/alipayRootCert.crt"); err != nil {
		log.Println("加载证书发生错误", err)
		return nil, fmt.Errorf("加载支付宝根证书失败: %w", err)
	}
	// 加载支付宝公钥证书
	if err = c.LoadAlipayCertPublicKeyFromFile("app/certificate/alipayPublicCert.crt"); err != nil {
		log.Println("加载证书发生错误", err)
		return nil, fmt.Errorf("加载支付宝证书失败: %w", err)
	}
	return c, nil
}
