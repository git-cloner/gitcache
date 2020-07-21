package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB = nil

func InitDb() {
	connectionString := os.Getenv("MYSQL_DSN")
	dbConn, _ = sql.Open("mysql", connectionString)
	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(5)
	_, err := dbConn.Query("select now()")
	if err != nil {
		dbConn.Close()
		dbConn = nil
	}
	log.Printf("connect to db: %v , err: %v", err == nil, err)
}

func SaveRepsInfo(name string, path string, utime time.Time) {
	if dbConn == nil {
		log.Printf("db error : connection is nil")
		return
	}
	var count int64
	rows, err := dbConn.Query("select count(*) from gitcache_repos where path = ?", path)
	if err != nil {
		log.Printf("db error : %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}
	if count == 0 {
		_, err = dbConn.Exec("insert into gitcache_repos (name,path,ctime,utime) values (?,?,?,?)", name, path, utime, utime)
		if err != nil {
			log.Printf("db error : %v", err)
		}
	} else {
		_, err = dbConn.Exec("update gitcache_repos set utime = ? where path = ?", utime, path)
		if err != nil {
			log.Printf("db error : %v", err)
		}
	}
}
