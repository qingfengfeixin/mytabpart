package model

import (
	"database/sql"
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
