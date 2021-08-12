package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"mytabpart/internal/conf"
)

type Dao struct {
	c      *conf.Config
	db     *sql.DB
	dbname string
}

var err error

func NewDao(c *conf.Config) *Dao {
	d := &Dao{
		c:      c,
		dbname: c.Dc.Dbname,
	}

	// mysql

	d.db, err = sql.Open(d.c.Dc.Driver, d.c.Dc.Dsn)
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

func (d *Dao) Close() {

	if d.db != nil {
		d.db.Close()
	}

}
