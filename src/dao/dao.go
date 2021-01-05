package dao

import (
	"database/sql"
	"fmt"
	"github.com/Unknwon/goconfig"
	_ "github.com/go-sql-driver/mysql"
)

type Dao struct {
	db *sql.DB
}

func NewDao() *Dao {
	// mysql
	dsn, err := getconfig().GetValue("mysql", "dsn")
	if err != nil {
		fmt.Println("get mysql cfg value err=", err)
		return nil
	}
	d := new(Dao)
	d.db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("connect mysql err=", err)
		return nil
	}
	if err = d.db.Ping(); err != nil {
		fmt.Println("mysql db ping err", err)
		return nil
	}
	return d
}

func(d *Dao) close(){
	d.db.Close()
}

func getconfig() *goconfig.ConfigFile {
	cfg, err := goconfig.LoadConfigFile("../conf/db.conf")
	if err != nil {
		fmt.Println("get cfg file err =", err)
		return nil
	}
	return cfg
}
