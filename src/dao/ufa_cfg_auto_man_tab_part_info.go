package dao

import (
	"database/sql"
	"fmt"
	"mytabpart/model"
	"time"
)

func (d *Dao) Getallparttab() (m []model.UFA_CFG_AUTO_MAN_TAB_PART_INFO) {
	ufa_cfg_tab_part := &model.UFA_CFG_AUTO_MAN_TAB_PART_INFO{}

	stmt, err := d.db.Prepare("select * from UFA_CFG_AUTO_MAN_TAB_PART_INFO")
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
		err = rows.Scan(&ufa_cfg_tab_part.TABLE_NAME,
			&ufa_cfg_tab_part.INTER_VAL,
			&ufa_cfg_tab_part.PARTITION_FLAG,
			&ufa_cfg_tab_part.SPILT_DATE,
			&ufa_cfg_tab_part.REAL_SPLID_DATE,
			&ufa_cfg_tab_part.TABLESPACE_NAME,
			&ufa_cfg_tab_part.KEY_COLUMN,
			&ufa_cfg_tab_part.RETENTION_HOUR,
		)
		if err != nil {
			fmt.Println("scan err =", err)
		}
		m = append(m, *ufa_cfg_tab_part)
	}
	return m
}

func (d *Dao) GetTabHi(tab string) (high string) {

	row := d.db.QueryRow("select max(TRIM( BOTH '\\'' FROM TRIM(partition_description) ) ) from "+
		"INFORMATION_SCHEMA.PARTITIONS where table_schema='noap' AND lower(table_name) =?", tab)

	var MaxPart sql.NullString

	err := row.Scan(&MaxPart)
	if err != nil {
		fmt.Println(err)
	}

	if MaxPart.Valid {
		high = MaxPart.String
	} else {
		high = time.Now().Format("2006-01-02 15")
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
		partname := "part" + "_" + v_hidate.Format("2006010215")

		sqlstr := "alter table " + tab + " add partition (partition " +
			partname + " values less then('" + highvalue + "'))"
		fmt.Println(sqlstr)

		if highvalue > maxday {
			break
		}
	}

}

func (d *Dao) TabParDrop(tab string, rent int) {

	expire := time.Now().Add(-time.Hour * time.Duration(rent)).Format("20060102")

	sqlstr := "SELECT partition_name FROM INFORMATION_SCHEMA.PARTITIONS " +
		"WHERE TABLE_SCHEMA = 'noap' AND table_name = ? and SUBSTRING(partition_name, 6)< ?"

	stmt, err := d.db.Prepare(sqlstr)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	}

	rows, err := stmt.Query(tab, expire)
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
		fmt.Println(dropstr)

	}

}
