package service

import (
	"mytabpart/internal/conf"
	"mytabpart/internal/dao"
)

type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

func NewService(c *conf.Config) *Service {
	s := &Service{
		c:   c,
		dao: dao.NewDao(c),
	}

	s.Mytabpart()

	return s
}

func (s *Service) write(str string) {

}
