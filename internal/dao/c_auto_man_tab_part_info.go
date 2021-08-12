package dao

import (
	"database/sql"
	"fmt"
	"mytabpart/internal/model"
	"time"
)

func (d *Dao) Getallparttab() (m []model.C_AUTO_MAN_TAB_PART_INFO) {
	cfg_tab_part := &model.C_AUTO_MAN_TAB_PART_INFO{}

	stmt, err := d.db.Prepare("select TABLE_NAME ,INTER_VAL ,RETENTION_HOUR,IS_PART from C_AUTO_MAN_TAB_PART_INFO where IS_PART= 1 ")
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
		err = rows.Scan(&cfg_tab_part.TABLE_NAME,
			&cfg_tab_part.INTER_VAL,
			&cfg_tab_part.RETENTION_HOUR,
			&cfg_tab_part.IS_PART,
		)
		if err != nil {
			fmt.Println("scan err =", err)
		}

		// check 是否是分区表
		if d.checkPart(cfg_tab_part.TABLE_NAME) {
			m = append(m, *cfg_tab_part)
		} else {
			// 不是分区表，但是配置表里标识的是分区表，则更新配置信息为非分区表
			// todo
		}

	}
	return m
}

// 通过 INFORMATION_SCHEMA.PARTITIONS 的partition_name 判断是否是分区表 如果存在分区 则是分区表
func (d *Dao) checkPart(tab string) bool {
	row := d.db.QueryRow("select max(partition_name) from "+
		"INFORMATION_SCHEMA.PARTITIONS where table_schema= ? AND lower(table_name) =?", d.dbname, tab)

	var part sql.NullString

	err := row.Scan(&part)
	if err != nil {
		fmt.Println(err)
	}

	if part.Valid {
		return true
	} else {
		return false
	}

}

func (d *Dao) GetTabHi(tab string) (high string) {

	row := d.db.QueryRow("select max(TRIM( BOTH '\\'' FROM TRIM(partition_description) ) ) from "+
		"INFORMATION_SCHEMA.PARTITIONS where table_schema= ? AND lower(table_name) =?", d.dbname, tab)

	var MaxPart sql.NullString

	err := row.Scan(&MaxPart)
	if err != nil {
		fmt.Println(err)
	}

	// 获取最近一个月零点的时间
	h_date, err := time.Parse("2006-01-02", time.Now().Add(-time.Hour*24*30).Format("2006-01-02"))

	if err != nil {
		fmt.Println(err)
	}

	if MaxPart.Valid {
		high = MaxPart.String
	} else {
		// 如果未取到分区时间，则设置最近一个月的零点为hidate
		high = h_date.Format("2006-01-02 15")
	}

	// 如果配置文件中设置了如果hidate超过一个月，则设置hidate为最近一个月
	if d.c.IsMon == 1 && high < h_date.Format("2006-01-02 15") {
		high = h_date.Format("2006-01-02 15")
	}

	return high

}

func (d *Dao) TabAddPart(tab string, maxday string, hidate string, inter int) {

	v_hidate, err := time.ParseInLocation("2006-01-02 15", hidate, time.Local)

	if err != nil {
		fmt.Println("hidate parse err=", err)
		return
	}

	for {

		v_hidate = v_hidate.Add(time.Hour * time.Duration(inter))

		highvalue := v_hidate.Format("2006-01-02 15")

		if highvalue > maxday {
			break
		}

		partname := "part" + "_" + v_hidate.Format("2006010215")

		sqlstr := "alter table " + tab + " add partition (partition " +
			partname + " values less than('" + highvalue + "'))"

		if _, err = d.db.Exec(sqlstr); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(sqlstr)
		}

	}

}

func (d *Dao) TabParDrop(tab string, rent int) {

	expire := time.Now().Add(-time.Hour * time.Duration(rent)).Format("2006010215")

	sqlstr := "SELECT partition_name FROM INFORMATION_SCHEMA.PARTITIONS " +
		"WHERE TABLE_SCHEMA = ? AND table_name = ? and SUBSTRING(partition_name, 10)< ?"

	stmt, err := d.db.Prepare(sqlstr)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	}

	rows, err := stmt.Query(d.dbname, tab, expire)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {

		var partab string
		err := rows.Scan(&partab)
		if err != nil {
			fmt.Println(err)
		}

		dropstr := "alter table " + tab + " drop partition " + partab

		if _, err = d.db.Exec(dropstr); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(dropstr)
		}

	}

}
