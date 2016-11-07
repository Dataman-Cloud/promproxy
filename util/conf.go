package util

import (
	"flag"
	"fmt"
	"log"

	"github.com/Unknwon/goconfig"
)

var (
	promflag    = flag.String("prometheus", "", "URL of the Prometheus Server <http://ip:port>")
	grafanaflag = flag.String("grafana", "", "URL of the Grafana Server <http://ip:port>")
	addrflag    = flag.String("addr", "", "API Server address <ip:port>")
)

// Conf contains the link of service Prometheus and Grafan
type Conf struct {
	PromServer string
	GrafServer string
	Addr       string
}

// Parse function will set the vaule of Conf fields. The command line will cover the vaule.
func (c *Conf) Parse(cfg *goconfig.ConfigFile) error {
	var err error
	c.PromServer, err = cfg.GetValue(goconfig.DEFAULT_SECTION, "promtheus")
	if err != nil {
		log.Fatalf("Can't get the value(%s): %s", "promtheus", err)
	}
	c.GrafServer, err = cfg.GetValue(goconfig.DEFAULT_SECTION, "grafana")
	if err != nil {
		log.Fatalf("Can't get the value(%s): %s", "grafana", err)
	}
	c.Addr, err = cfg.GetValue(goconfig.DEFAULT_SECTION, "address")
	if err != nil {
		log.Fatalf("Can't get the value(%s): %s", "grafana", err)
	}

	if *promflag != "" {
		c.PromServer = *promflag
	}

	if *grafanaflag != "" {
		c.GrafServer = *grafanaflag
	}

	if *addrflag != "" {
		c.Addr = *addrflag
	}

	fmt.Println(c.PromServer, c.GrafServer)
	return nil
}
