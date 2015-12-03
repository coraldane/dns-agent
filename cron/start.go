package cron

import (
	"bytes"
	"encoding/json"
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

func ModifyRecord(domainId int, recordId, subDomain, strIp string) bool {
	data := url.Values{}

	data.Add("login_email", g.Config().LoginEmail)
	data.Add("login_password", g.Config().LoginPassword)
	data.Add("format", "json")
	data.Add("domain_id", strconv.Itoa(domainId))
	data.Add("record_id", recordId)
	data.Add("sub_domain", subDomain)
	data.Add("record_type", "A")
	data.Add("record_line", "默认")
	data.Add("value", strIp)

	strResponse, err := Post("https://dnsapi.cn/Record.Modify", data)
	if nil != err {
		log.Printf("RECORD_MODIFY ERROR", err)
	} else {
		log.Println("RECORD_MODIFY_RESPONSE:", strResponse, err)
		if "" != strResponse && strings.Contains(strResponse, `"code":"1"`) {
			return true
		} else {
			log.Printf("domainId:%d,recordId:%s,sub_domain:%s\n", domainId, recordId, subDomain)
		}
	}
	return false
}

func GetDomainList() []g.DomainResult {
	data := url.Values{}

	data.Add("login_email", g.Config().LoginEmail)
	data.Add("login_password", g.Config().LoginPassword)
	data.Add("format", "json")

	strResponse, err := Post("https://dnsapi.cn/Domain.List", data)
	if nil == err && "" != strResponse && strings.Contains(strResponse, `"code":"1"`) {
		var dlr g.DomainListResult
		err = json.Unmarshal(bytes.NewBufferString(strResponse).Bytes(), &dlr)
		if nil != err {
			log.Printf("decode DOMAIN_LIST response fail %v\n", err)
		} else {
			if g.Config().Debug {
				for _, domain := range dlr.Domains {
					log.Printf("domain_id:%d,name:%s,created:%v,updated:%v\n", domain.Id, domain.Name, domain.Created, domain.Updated)
				}
			}
			return dlr.Domains
		}
	} else {
		log.Printf("GET_DOMAIN_LIST RESPONSE<<<====%s, error: %v", strResponse, err)
	}
	return nil
}

func GetRecordList(domainId int) []g.RecordResult {
	data := url.Values{}

	data.Add("login_email", g.Config().LoginEmail)
	data.Add("login_password", g.Config().LoginPassword)
	data.Add("format", "json")
	data.Add("domain_id", strconv.Itoa(domainId))

	strResponse, err := Post("https://dnsapi.cn/Record.List", data)
	if nil == err && "" != strResponse && strings.Contains(strResponse, `"code":"1"`) {
		var rlr g.RecordListResult
		err = json.Unmarshal(bytes.NewBufferString(strResponse).Bytes(), &rlr)
		if nil != err {
			log.Printf("decode RECORD_LIST response fail %v\n", err)
		} else {
			var records []g.RecordResult
			for _, record := range rlr.Records {
				if "A" != record.Type {
					continue
				}
				records = append(records, record)
				if g.Config().Debug {
					log.Printf("domain_id:%d,name:%s,record_id:%s,name:%s,value:%s,status:%s\n",
						domainId, rlr.Domain.Name, record.Id, record.Name, record.Value, record.Status)
				}
			}

			return records
		}
	} else {
		log.Printf("GET_RECORD_LIST RESPONSE<<<====%s, error: %v", strResponse, err)
	}
	return nil
}

func getIp() string {
	strUrl := "http://ops.mixuan.org/api/ip"
	if g.Config().Debug {
		log.Printf("REQUEST_URL:%s\n", strUrl)
	}
	httpRequest := httplib.Get(strUrl).SetTimeout(3*time.Second, 10*time.Second)
	httpResponse, err := httpRequest.Bytes()
	if nil != err {
		log.Println("GET_IP error", err)
		return ""
	}

	strIp := ""
	var resp g.ServletResponse
	err = json.Unmarshal(httpResponse, &resp)
	if err != nil {
		log.Printf("decode GET_IP response fail %v\n", err)
	} else if false == resp.Success {
		log.Printf("GET_IP fail %s\n", resp.Message)
	} else {
		strIp = resp.Message
	}

	if g.Config().Debug {
		log.Println("RESPONSE_IP:", strIp)
	}
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
