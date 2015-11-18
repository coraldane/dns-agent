package cron

import (
	"fmt"
	"github.com/toolkits/net/httplib"
	"log"
	"strings"
	"time"
)

var lastIp string

func UpdateIpRecord() {
	strIp := getIp()
	if "" == strIp || lastIp == strIp {
		return
	}

	for _, domain := range g.Config().Domains {
		for _, record := range domain.Records {
			data := url.Values{}

			data.Add("login_email", g.Config().LoginEmail)
			data.Add("login_password", g.Config().LoginPassword)
			data.Add("format", "json")
			data.Add("domain_id", strconv.Itoa(domain.DomainId))
			data.Add("record_id", strconv.Itoa(record.RecordId))
			data.Add("sub_domain", record.SubDomain)
			data.Add("record_type", "A")
			data.Add("record_line", "默认")
			data.Add("value", strIp)
			strResponse, err := utils.Post("https://dnsapi.cn/Record.Modify", data)
			log.Println("RESPONSE:", strResponse, err)
			if "" != strResponse && strings.Contains(strResponse, `"code":"1"`) {
				lastIp = strIp
			}
		}
	}
}

func getIp() string {
	strUrl := fmt.Sprintf("https://cgi1.apnic.net/cgi-bin/my-ip.php?callback=jQuery%d_%d&_=%d", time.Now().Unix()*1234+1234, time.Now().Unix(), time.Now().Unix())
	log.Printf("REQUEST_URL:%s\n", strUrl)
	httpRequest := httplib.Post(strUrl).SetTimeout(time.Second*10, time.Minute)
	httpResponse, err := httpRequest.Bytes()
	if nil != err {
		log.Println("GET_IP error", err)
		return ""
	}
	strContent := string(httpResponse)
	if !strings.Contains(strContent, "{") {
		log.Printf("IP_DATA error====>%s\n", strContent)
		return ""
	}

	strIp := strContent[strings.Index(strContent, `"`)+1 : strings.LastIndex(strContent, `"`)]
	log.Println("RESPONSE_IP:", strIp)
	return strIp
}
