package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	macaron "gopkg.in/macaron.v1"

	"github.com/Dataman-Cloud/promproxy/api"
	"github.com/Dataman-Cloud/promproxy/util"
	"github.com/Unknwon/goconfig"
)

func main() {
	flag.Parse()
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		log.Fatalf("Can't load the config file: %s", err)
		os.Exit(1)
	}

	conf := new(util.Conf)
	err = conf.Parse(cfg)
	if err != nil {
		log.Fatalf("Failed to parse the conf: %v", err)
		os.Exit(1)
	}

	client := api.NewClient(*conf)

	m := macaron.Classic()
	m.Get("/query", client.Query)
	m.Get("/query_range", client.QueryRange)
	m.Get("/graph", client.GetGraph)

	log.Printf("Server is running on %s ...", conf.Addr)
	log.Println(http.ListenAndServe(conf.Addr, m))

}
