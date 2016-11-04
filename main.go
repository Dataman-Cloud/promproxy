package main

import (
	"flag"
	"log"

	"github.com/Dataman-Cloud/promproxy/util"
	"github.com/Unknwon/goconfig"
)

func main() {
	flag.Parse()
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		log.Fatalf("Can't load the config file: %s", err)
	}

	conf := new(util.Conf)
	err = conf.Parse(cfg)
	if err != nil {
		log.Fatalf("Failed to parse the conf: %v", err)
	}
}
