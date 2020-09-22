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
		_, err = dbConn.Exec("insert into gitcache_repos (name,path,ctime,utime,hitcount) values (?,?,?,?,?)", name, path, utime, utime, 0)
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

func AddHitCount(path string) {
	go Stats("cachehit")
	if dbConn == nil {
		log.Printf("db error : connection is nil")
		return
	}
	_, err := dbConn.Exec("update gitcache_repos set hitcount = hitcount + 1 where path = ?", path)
	if err != nil {
		log.Printf("db error : %v", err)
	}
}

func CacheExists(path string) bool {
	if dbConn == nil {
		log.Printf("db error : connection is nil")
		return false
	}
	var count int64
	rows, err := dbConn.Query("select count(*) from gitcache_repos where path = ?", path)
	if err != nil {
		log.Printf("db error : %v", err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}
	if count == 0 {
		return false
	} else {
		return true
	}
}

func Stats(stat_class string) {
	log.Printf("Stats : %v", stat_class)
	if dbConn == nil {
		log.Printf("db error : connection is nil")
		return
	}
	now := time.Now()
	ns := now.Format("2006-01-02")
	var count int64
	rows, err := dbConn.Query("select count(*) from gitcache_stats where stime = ?", ns)
	if err != nil {
		log.Printf("db error : %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}
	if count == 0 {
		_, err = dbConn.Exec("insert into gitcache_stats (stime,cachehit,redirect,visit,vipvisit,search,imagetest) "+
			"values (?,?,?,?,?,?,?)", ns, 0, 0, 0, 0, 0, 0)
		if err != nil {
			log.Printf("db error : %v", err)
		}
	}
	_, err = dbConn.Exec("update gitcache_stats "+
		" set "+stat_class+" = "+stat_class+"  + 1 where stime = ?", ns)
	if err != nil {
		log.Printf("db error : %v", err)
	}
}

func UpdateReposDetail(path string, star int64, lang string, desc string, upt time.Time) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("process recover: %s\n", err)
		}
	}()
	if dbConn == nil {
		log.Printf("db error : connection is nil")
		return
	}
	_, err := dbConn.Exec("update gitcache_repos set starcount = ?,language = ?,description = ?,updated_at = ? where path = ?", star, lang, desc, upt, path)
	if err != nil {
		log.Printf("db error : %v", err)
	}
}
