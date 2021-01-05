package service

import "mytabpart/dao"

type Service struct {
	dao *dao.Dao
}

func NewService() *Service {
	s := &Service{
		dao: dao.NewDao(),
	}
	return s
}
