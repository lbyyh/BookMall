package logic

import (
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"library-study/app/model"
	"sort"
	"time"
)

// TOKEN 假设您在Go代码中定义了一个名为TOKEN的常量，用于存储您的令牌值
const TOKEN = "111111"

// 配置公众号的token
func CheckSignature(c *gin.Context) {
	// 获取查询参数中的签名、时间戳和随机数
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	// 创建包含令牌、时间戳和随机数的字符串切片
	tmpArr := []string{TOKEN, timestamp, nonce}
	// 对切片进行字典排序
	sort.Strings(tmpArr)
	// 将排序后的元素拼接成单个字符串
	tmpStr := ""
	for _, v := range tmpArr {
		tmpStr += v
	}
	// 对字符串进行SHA-1哈希计算
	tmpHash := sha1.New()
	tmpHash.Write([]byte(tmpStr))
	tmpStr = fmt.Sprintf("%x", tmpHash.Sum(nil))
	fmt.Println(tmpStr)
	fmt.Println(signature)
	// 将计算得到的签名与请求中提供的签名进行比较，并根据结果发送相应的响应
	if tmpStr == signature {
		c.String(200, echostr)
		model.RedisDB.Set(c, "library:token", tmpStr, 7*24*time.Hour)
	} else {
		c.String(403, "签名验证失败 "+timestamp)
	}
}
