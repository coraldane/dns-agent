package main

import (
	"flag"
	"fmt"
	"github.com/coraldane/dns-agent/cron"
	"github.com/coraldane/dns-agent/g"
	"log"
	"os"
	// "time"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	log.Println(g.Config())

	g.InitRedisConnPool()

	go cron.SyncDomainRecord()
	go cron.Heartbeat()

	select {}
	// ticker1 := time.NewTicker(time.Duration(g.Config().Interval) * time.Second)
	// for {
	// 	select {
	// 	case <-ticker1.C:
	// 		go func() {
	// 			cron.UpdateIpRecord()
	// 		}()
	// 	}
	// }
}
