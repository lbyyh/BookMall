package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type BooksResponse struct {
	Code    int        `json:"code"`
	Data    []BookInfo `json:"data"`
	Message string     `json:"message"`
}

// 列出图书的命令
var listBooksCmd = &cobra.Command{
	Use:   "list",
	Short: "列出库存中的所有图书",
	Run: func(cmd *cobra.Command, args []string) {
		books, err := getBooks()
		if err != nil {
			log.Fatalf("无法获取图书列表: %v", err)
		}

		fmt.Println("图书列表：")
		for _, book := range books {
			fmt.Printf("ID: %d, 名称: %s, 作者: %s, 出版社: %s, 译者: %s, 数量: %d, 出版日期: %s, 页数: %d, ISBN: %s, 价格: %s, 内容简介: %s, 作者简介: %s, 封面图片: %s\n",
				book.Id, book.BookName, book.Author, book.PublishingHouse, book.Translator, book.Num, book.PublishDate, book.Pages, book.ISBN, book.Price, book.BriefIntroduction, book.AuthorIntroduction, book.ImgUrl)
		}
	},
}

// getBooks 从 HTTP API 获取图书列表
func getBooks() ([]BookInfo, error) {
	resp, err := http.Get("http://127.0.0.1:8087/book_info/list")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var booksResponse BooksResponse
	err = json.Unmarshal(body, &booksResponse)
	if err != nil {
		return nil, err
	}

	if booksResponse.Code != 0 {
		return nil, fmt.Errorf("server returned non-zero code: %d, message: %s", booksResponse.Code, booksResponse.Message)
	}

	return booksResponse.Data, nil
}

var (
	// 用于标记。
	cfgFile string

	// RootCmd 表示没有调用子命令时的基础命令。
	RootCmd = &cobra.Command{
		Use:   "library",
		Short: "这是一个图书管理系统",
		Long:  `这个图书管理系统基于Cobra构建。你可以通过这个应用添加、列出和删除图书记录。`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// 在这里为addBookCmd添加参数
	addBookCmd.Flags().StringP("title", "t", "", "书籍的标题")
	addBookCmd.Flags().StringP("author", "a", "", "书籍的作者")
	// 这里可以根据需要，持久化添加或修改flag。
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件的路径 (默认是 $HOME/.library.yaml)")

	// 这里注册图书管理系统需要的命令
	RootCmd.AddCommand(addBookCmd)
	RootCmd.AddCommand(listBooksCmd)
	RootCmd.AddCommand(removeBookCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Printf("寻找主目录失败: %v\n", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("使用配置文件:", viper.ConfigFileUsed())
	} else {
		fmt.Println("无法读取配置文件，只使用环境变量和默认值")
	}
}

// 以下是图书管理的命令结构体和初始化函数示例
// 添加图书的命令
var addBookCmd = &cobra.Command{
	Use:   "add",
	Short: "添加一本新书到库存",
	Run: func(cmd *cobra.Command, args []string) {
		// 假设你需要书名和作者作为参数来添加一本书
		title, err := cmd.Flags().GetString("title")
		if err != nil {
			log.Fatalf("获取标题参数错误: %s\n", err)
		}
		author, err := cmd.Flags().GetString("author")
		if err != nil {
			log.Fatalf("获取作者参数错误: %s\n", err)
		}

		// 此处添加你的逻辑，例如添加到数据库等
		fmt.Printf("已添加书籍: %s 作者: %s\n", title, author)
	},
}

// 删除图书的命令
var removeBookCmd = &cobra.Command{
	Use:   "remove",
	Short: "从库存中删除一本书",
	Run: func(cmd *cobra.Command, args []string) {
		// 删除图书的逻辑
	},
}
