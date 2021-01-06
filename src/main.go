package main

import (
	"flag"
	"fmt"
	"mytabpart/conf"
	"mytabpart/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	flag.Parse()
	conf.Init()

	service.NewService(conf.Conf)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			fmt.Println("Program Mytabpart Exit...", s)
		case syscall.SIGQUIT:
			fmt.Println("Program Mytabpart Quit", s)
		default:
			fmt.Println("other signal", s)
		}
	}

}
