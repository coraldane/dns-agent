package cron

import (
	"github.com/coraldane/dns-agent/g"
	"log"
	"time"
)

func Heartbeat() {
	// SleepRandomDuration()
	for {
		heartbeat()
		d := time.Duration(g.Config().Interval) * time.Second
		time.Sleep(d)
	}
}

var (
	lastIp string
)

func heartbeat() {
	strIp := getIp()
	if "" == strIp || lastIp == strIp {
		return
	}

	domainResults := GetDomainList()
	if nil == domainResults {
		return
	}

	var modifyResult bool
	for _, domain := range g.Config().Domains {
		for _, domainResult := range domainResults {
			if domain.DomainName == domainResult.Name {
				recordResults := GetRecordList(domainResult.Id)
				for _, recordName := range domain.RecordNames {
					for _, recordResult := range recordResults {
						if recordName == recordResult.Name && strIp != recordResult.Value {
							modifyResult = ModifyRecord(domainResult.Id, recordResult.Id, recordName, strIp)
						}
					}
				}
			}
		}
	}

	if modifyResult {
		lastIp = strIp
		log.Println("last ip have changed into ", strIp)
	}
}
