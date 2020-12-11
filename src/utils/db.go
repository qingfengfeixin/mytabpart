package utils

import (
	"database/sql"
	"github.com/Unknwon/goconfig"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var MyDB *sql.DB

func NewDB() *sql.DB {
	if MyDB == nil {
		MyDB = connectdb()
	}
	return MyDB
}

func connectdb() *sql.DB {
	var (
		cfg *goconfig.ConfigFile
		err error
		dsn string
	)

	if cfg, err = goconfig.LoadConfigFile("../conf/db.conf"); err != nil {
		log.Fatalf("无法加载配置文件：%s", err)
		return nil
	}
	if dsn, err = cfg.GetValue("mysql", "dsn"); err != nil {
		log.Fatalf("无法获取数据库配置信息: %s", err)
		return nil
	}
	if MyDB, err = sql.Open("mysql", dsn); err != nil {
		log.Fatalf("无法打开数据库连接:%s", err)
		return nil
	}
	if err = MyDB.Ping(); err != nil {
		log.Fatalf("数据库无法登录:%s", err)
		return nil
	}
	return MyDB
}
