package main

import (
	"library-study/app"
	"library-study/app/model"
	"log"
	"os"
)

// @contact.name   library API
// @contact.email  hzs

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	if err := model.RootCmd.Execute(); err != nil {
		log.Fatalf("命令执行错误: %s\n", err)
		os.Exit(1)
	}

	app.Start()
}

//package main
//
//import (
//	_ "github.com/go-sql-driver/mysql"
//)
//
//func main() {
// 数据库连接配置
//my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", "book", "p8YwmmZrL7y7trjt", "192.168.30.23:3306", "book")
//Conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
//if err != nil {
//	fmt.Printf("err:%s\n", err)
//	panic(err)
//}
//Db = Conn

//	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", "book", "p8YwmmZrL7y7trjt", "192.168.30.133:3306", "book")
//	db, err := sql.Open("mysql", dsn)
//	if err != nil {
//		log.Fatalf("Failed to open database: %v", err)
//	}
//	defer db.Close()
//
//	snowflake.Epoch = time.Now().UnixNano()/1000000 - snowflake.Epoch
//
//	// 创建雪花算法节点
//	node, err := snowflake.NewNode(1)
//	if err != nil {
//		log.Fatalf("Failed to create snowflake node: %v", err)
//	}
//
//	// 查询表中所有记录的ID
//	rows, err := db.Query("SELECT id FROM book_info")
//	if err != nil {
//		log.Fatalf("Failed to execute query: %v", err)
//	}
//	defer rows.Close()
//
//	// 开启事务
//	tx, err := db.Begin()
//	if err != nil {
//		log.Fatalf("Failed to begin transaction: %v", err)
//	}
//	defer tx.Rollback()
//
//	for rows.Next() {
//		var id int
//		if err := rows.Scan(&id); err != nil {
//			log.Fatalf("Failed to scan row: %v", err)
//		}
//
//		// 为每个ID生成新的UID
//		newUID := node.Generate().Int64()
//
//		_, err = tx.Exec("UPDATE book_info SET uid = ? WHERE id = ?", newUID, id)
//		if err != nil {
//			log.Fatalf("Failed to update uid: %v", err)
//		}
//	}
//
//	// 提交事务
//	if err := tx.Commit(); err != nil {
//		log.Fatalf("Failed to commit transaction: %v", err)
//	}
//
//	log.Println("All UIDs have been successfully updated.")
//}
