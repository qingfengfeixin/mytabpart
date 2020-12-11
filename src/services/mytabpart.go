package services

import (
	"fmt"
	"mytabpart/models"
	"time"
)

type Mytabpart struct {
	UFA_CFG_AUTO_MAN_TAB_PART_INFOS []models.UFA_CFG_AUTO_MAN_TAB_PART_INFO
}

func (this *Mytabpart) PartTab() {
	u := models.UFA_CFG_AUTO_MAN_TAB_PART_INFO{}
	this.UFA_CFG_AUTO_MAN_TAB_PART_INFOS = u.Getall()
	for i, v := range this.UFA_CFG_AUTO_MAN_TAB_PART_INFOS {
		fmt.Println(i, v)
		this.Partadd(v.TABLE_NAME, v.INTER_VAL)
		this.Partdrop(v.TABLE_NAME, v.RETENTION_HOUR)
	}
}

func (this *Mytabpart) Partadd(tab string, inter int) {
	// 小时分区预建立7天
	// 天分区预建立7天
	// 周分区预建2周
	// 其他分区预建2月

	var (
		maxDay string
		high   string
	)

	switch inter {
	case 1, 24:
		maxDay = time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02")
	case 168:
		maxDay = time.Now().Add(time.Hour * 24 * 7 * 2).Format("2006-01-02")
	default:
		maxDay = time.Now().Add(time.Hour * 24 * 7 * 2).Format("2006-01-02")
	}

	fmt.Println("maxday=", maxDay)

	tabHi := models.GetHi(tab)

	if tabHi.Valid {
		high = tabHi.String
	} else {
		high = time.Now().Format("2006-01-02 15")
	}

	fmt.Println("high=", high)

	models.TabAddPart(tab, maxDay, high, inter)

}

func (this *Mytabpart) Partdrop(tab string, rent int) {

	models.TabParDrop(tab, rent)
}
