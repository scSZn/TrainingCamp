package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Id   int
	Namg string
	Age  int
}

var db *sql.DB

func main() {
	err := initDB()
	if err != nil {
		log.Fatalf("connect to database failed, err: %v", err)
		return
	}

	name, err := service(2)
	// 这里统一处理error
	if errors.Is(err, sql.ErrNoRows) {
		// 在这里返回一些特定的响应信息
		log.Printf("not found, err: %v", err)
		return
	}
	if err != nil {
		log.Fatalf("unexpected error, err: %v", err)
		return
	}
	fmt.Println(name)
}

func initDB() error {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/db2019")
	if err != nil {
		return fmt.Errorf("database connection string have an syntax error, err: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("ping database failed, err: %w", err)
	}
	return nil
}

func service(id int) (string, error) {
	return dao(id)
}

func dao(id int) (string, error) {
	var name string
	err := db.QueryRow("select name from person where id = ?", id).Scan(&name)
	/*
		我认为应该包装后抛给上层，或者直接抛给上层，不应该在这里进行处理
		因为如果在dao层处理了，那么返回的结果应该就是类型的零值
		这里以string为例，返回为空字符串""有可能是数据库中存储的就是空字符串，也可能是查询出错了，返回了一个字符串类型的零值
		如果在dao层处理了error，按照错误只处理一次的原则，就不应该往上层抛出错误，这样就不能分清楚上面所说的两种情况
		因此，我认为应该直接抛给上层，可以封装一下自定义信息在里面
		请老师指正
	*/
	if err != nil {
		return "", fmt.Errorf("query failed, error: %w", err)
	}
	return name, nil
}
