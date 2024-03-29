package service

import (
	"fmt"
	"time"
)

func (s *Service) Mytabpart() {

	for i, v := range s.dao.Getallparttab() {
		fmt.Println(i, v)
		s.Partadd(v.TABLE_NAME, v.INTER_VAL)

		// 是否需要根据保留时长drop分区 如果为1 则执行drop操作
		// 暂时先注释，防止误操作

		/*if s.c.IsDrop==1{
			s.Partdrop(v.TABLE_NAME, v.RETENTION_HOUR)
		}*/

	}
	defer s.dao.Close()

}

func (s *Service) Partadd(tab string, inter int) {
	// 小时分区预建立7天
	// 天分区预建立7天
	// 周分区预建2周
	// 其他分区预建2月

	var maxDay string

	switch inter {
	case 1, 24:
		maxDay = time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02")
	case 168:
		maxDay = time.Now().Add(time.Hour * 24 * 7 * 2).Format("2006-01-02")
	default:
		maxDay = time.Now().Add(time.Hour * 24 * 7 * 2).Format("2006-01-02")
	}

	fmt.Println("maxday=", maxDay)

	high := s.dao.GetTabHi(tab)

	fmt.Println("high=", high)

	s.dao.TabAddPart(tab, maxDay, high, inter)

}

func (s *Service) Partdrop(tab string, rent int) {

	s.dao.TabParDrop(tab, rent)
}
