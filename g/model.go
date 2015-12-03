package g

import (
	"strings"
	"time"
)

type ServletResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type myTime struct {
	time.Time
}

func (t *myTime) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse("2006-01-02 15:04:05", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

type DomainResult struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Records string `json:"records"`
	Status  string `json:"status"`
	Created myTime `json:"created_on"`
	Updated myTime `json:"updated_on"`
}

type DomainListResult struct {
	Domains []DomainResult `json:"domains"`
}

type RecordListResult struct {
	Domain  DomainResult   `json:"domain"`
	Records []RecordResult `json:"records"`
}

type RecordResult struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Value   string `json:"value"`
	Status  string `json:"status"`
	Updated myTime `json:"updated_on"`
}
