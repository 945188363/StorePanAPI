package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init()  {
	db,_ = sql.Open("mysql","root:123456@tcp(localhost:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("Fail to init DB,error:%s",err.Error())
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return db
}

