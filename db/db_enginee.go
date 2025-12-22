package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	HOST     = "localhost"
	DATABASE = "code_lab"
	PORT     = 5432
	USER     = "postgres"
	PASSWORD = "P@ssw0rd"
)

var (
	table = ""
)

var dbPool *sql.DB

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func DBconn() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, USER, PASSWORD, DATABASE)

	pool, err := sql.Open("postgres", connStr)
	CheckError(err)

	// 建議配置（根據實際需求調整）
	pool.SetMaxOpenConns(25) // 不要超過 PostgreSQL 的 max_connections
	pool.SetMaxIdleConns(10) // 通常設為 MaxOpenConns 的 25-50%

	dbPool = pool

	fmt.Println("DB connected success")
}

func CloseDBConn() {
	dbPool.Close()
	fmt.Println("DB connect closed")

}
