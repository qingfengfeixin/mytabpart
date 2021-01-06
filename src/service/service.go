package service

import (
	"mytabpart/conf"
	"mytabpart/dao"
)

type Service struct {
	c *conf.Config
	dao *dao.Dao
}

func NewService(c *conf.Config) *Service {
	s := &Service{
		c :c,
		dao: dao.NewDao(c),
	}

	go s.Mytabpart()

	return s
}
