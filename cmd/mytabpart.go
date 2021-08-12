package cmd

import (
	"flag"
	"mytabpart/internal/conf"
	"mytabpart/internal/service"
)

func MyTabPart() {

	flag.Parse()
	conf.Init()

	service.NewService(conf.Conf)

}
