package models

import (
	"database/sql"
	"fmt"
	"mytabpart/utils"
	"time"
)

type UFA_CFG_AUTO_MAN_TAB_PART_INFO struct {
	TABLE_NAME      string
	INTER_VAL       int
	PARTITION_FLAG  sql.NullString
	SPILT_DATE      sql.NullTime
	REAL_SPLID_DATE sql.NullTime
	TABLESPACE_NAME sql.NullString
	KEY_COLUMN      sql.NullString
	RETENTION_HOUR  int
}

func (this *UFA_CFG_AUTO_MAN_TAB_PART_INFO) Getall() []UFA_CFG_AUTO_MAN_TAB_PART_INFO {
	var m []UFA_CFG_AUTO_MAN_TAB_PART_INFO
	db := utils.NewDB()
	//defer db.Close()

	stmt, err := db.Prepare("select * from UFA_CFG_AUTO_MAN_TAB_PART_INFO")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	}

	rows, err := stmt.Query()
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&this.TABLE_NAME,
			&this.INTER_VAL,
			&this.PARTITION_FLAG,
			&this.SPILT_DATE,
			&this.REAL_SPLID_DATE,
			&this.TABLESPACE_NAME,
			&this.KEY_COLUMN,
			&this.RETENTION_HOUR,
		)
		if err != nil {
			fmt.Println("scan err =", err)
		}
		m = append(m, *this)
	}
	return m
}

func GetHi(tab string) (MaxPart sql.NullString) {
	db := utils.NewDB()
	//defer db.Close()

	row := db.QueryRow("select max(TRIM( BOTH '\\'' FROM TRIM(partition_description) ) ) from "+
		"INFORMATION_SCHEMA.PARTITIONS where table_schema='noap' AND lower(table_name) =?", tab)

	err := row.Scan(&MaxPart)
	if err != nil {
		fmt.Println(err)
	}
	return MaxPart

}

func TabAddPart(tab string, maxday string, hidate string, inter int) {
	v_hidate := hidate

	for {
		t, err := time.ParseInLocation("2006-01-02 15", v_hidate, time.Local)
		if err != nil {
			fmt.Println(err)
			return
		}

		v_hidate = t.Add(time.Hour * time.Duration(inter)).Format("2006-01-02 15")
		parname := "part" + "_" + t.Add(time.Hour*time.Duration(inter)).Format("2006010215")

		sqlstr := "alter table " + tab +
			" add partition (partition " +
			parname + " values less then('" +
			v_hidate + "'))"
		fmt.Println(sqlstr)

		if v_hidate > maxday {
			break
		}
	}

}

func TabParDrop(tab string, rent int) {

	var (
		partab  string
		datestr string
	)

	expire := time.Now().Add(-time.Hour * time.Duration(rent)).Format("20060102")

	sqlstr := "SELECT partition_name,SUBSTRING(partition_name, 6) FROM INFORMATION_SCHEMA.PARTITIONS " +
		"WHERE TABLE_SCHEMA = 'noap' AND table_name = ?"

	db := utils.NewDB()

	stmt, err := db.Prepare(sqlstr)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	}

	rows, err := stmt.Query(tab)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err := rows.Scan(&partab, &datestr)
		if err != nil {
			fmt.Println(err)
		}

		if datestr < expire {
			dropstr := "alter table " + tab + " drop partition " + partab
			fmt.Println(dropstr)
		}

	}

}
