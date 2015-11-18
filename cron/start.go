package cron

import (
	"bytes"
	"github.com/axgle/mahonia"
	"github.com/coraldane/dns-agent/g"
	"github.com/toolkits/net/httplib"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
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

			strResponse, err := Post("https://dnsapi.cn/Record.Modify", data)

			// httpRequest := httplib.Post("https://dnsapi.cn/Record.Modify").SetTimeout(time.Second*10, time.Minute)
			// log.Println(GenerateQueryString(data))
			// httpRequest.Body(GenerateQueryString(data))
			// httpResponse, err := httpRequest.Bytes()
			if nil != err {
				log.Printf("RECORD_MODIFY ERROR", err)
			} else {
				// strResponse := string(httpResponse)
				log.Println("RECORD_MODIFY_RESPONSE:", strResponse, err)
				if "" != strResponse && strings.Contains(strResponse, `"code":"1"`) {
					lastIp = strIp
				}
			}
		}
	}
}

func getIp() string {
	strUrl := "http://1111.ip138.com/ic.asp"
	log.Printf("REQUEST_URL:%s\n", strUrl)
	httpRequest := httplib.Get(strUrl).SetTimeout(time.Second*10, time.Minute)
	httpResponse, err := httpRequest.Bytes()
	if nil != err {
		log.Println("GET_IP error", err)
		return ""
	}

	decoder := mahonia.NewDecoder("GBK")
	strContent := decoder.ConvertString(string(httpResponse))
	if !strings.Contains(strContent, "[") {
		log.Printf("IP_DATA error====>%s\n", strContent)
		return ""
	}

	strIp := strContent[strings.Index(strContent, "[")+1 : strings.Index(strContent, "]")]
	log.Println("RESPONSE_IP:", strIp)
	return strIp
}

func Post(strUrl string, data url.Values) (string, error) {
	response, err := http.PostForm(strUrl, data)

	if nil != err {
		return "post data error", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return "read body error", err
	}
	return string(body), err
}

func GenerateQueryString(v url.Values) string {
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}
