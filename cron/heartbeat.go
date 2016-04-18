package cron

import (
	"github.com/coraldane/dns-agent/g"
	"log"
	"time"
)

func SyncDomainRecord() {
	for {
		syncRecordList()
		d := time.Duration(5) * time.Minute
		time.Sleep(d)
	}
}

var (
	lastIp    string
	domainMap map[string]int
	recordMap map[string][]g.RecordResult
)

func Heartbeat() {
	time.Sleep(time.Duration(5) * time.Second)
	for {
		heartbeat()
		d := time.Duration(g.Config().Interval) * time.Second
		time.Sleep(d)
	}
}

func syncRecordList() {
	domainResults := GetDomainList()
	if nil == domainResults {
		return
	}

	domainMap = make(map[string]int)
	recordMap = make(map[string][]g.RecordResult)
	for _, domainResult := range domainResults {
		domainMap[domainResult.Name] = domainResult.Id
		recordResults := GetRecordList(domainResult.Id)
		if nil != recordResults {
			recordMap[domainResult.Name] = recordResults
		}
	}
}

func heartbeat() {
	strIp := getIp()
	if "" == strIp || lastIp == strIp {
		return
	}

	if 0 == len(domainMap) {
		return
	}

	var modifyResult bool

	for _, domain := range g.Config().Domains {
		if _, ok := domainMap[domain.DomainName]; ok {
			if recordResults, exists := recordMap[domain.DomainName]; exists {
				for _, recordResult := range recordResults {
					if existsRecordName(domain.DomainName, recordResult.Name) && strIp != recordResult.Value {
						modifyResult = ModifyRecord(domainMap[domain.DomainName], recordResult.Id, recordResult.Name, strIp)
					}
				}
			}
		}
	}

	if modifyResult {
		lastIp = strIp
		log.Println("last ip have changed into ", strIp)

		// send to redis
		rc := g.RedisConnPool.Get()
		defer rc.Close()
		rc.Do("LPUSH", "COMMAND_udai", "service nginx restart")
	}
}

func existsRecordName(domainName, recordName string) bool {
	for _, domain := range g.Config().Domains {
		if domain.DomainName == domainName {
			for _, record := range domain.RecordNames {
				if record == recordName {
					return true
				}
			}
		}
	}
	return false
}
