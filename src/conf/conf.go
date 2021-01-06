package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     = &Config{}
)

type Config struct {
	Title string
	Dc    mysql `toml:"mysql"`
}

type mysql struct {
	Driver string
	Dsn    string
	Dbname string
}

func init() {
	flag.StringVar(&confPath, "conf", "../conf/conf.toml", "default config path")

}

func Init() {
	if _, err := toml.DecodeFile(confPath, &Conf); err != nil {
		panic(err)
	}

}
