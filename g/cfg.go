package g

import (
	"encoding/json"
	"fmt"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type GlobalConfig struct {
	LogLevel      string         `json:"log_level"`
	Interval      int            `json:"interval"`
	LoginEmail    string         `json:"login_email"`
	LoginPassword string         `json:"login_pwd"`
	Domains       []DomainConfig `json:"domains"`
	Redis         RedisConfig    `json:"redis"`
}

type RedisConfig struct {
	Enabled      bool   `json:"enabled"`
	Dsn          string `json:"dsn"`
	MaxIdle      int    `json:"maxIdle"`
	ConnTimeout  int    `json:"connTimeout"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
	Passwd       string `json:"passwd"`
}

type DomainConfig struct {
	DomainName  string   `json:"domain_name"`
	RecordNames []string `json:"record_names"`
}

var (
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("config file %s is nonexistent", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read config file %s fail %s", cfg, err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse config file %s fail %s", cfg, err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
	return nil
}
