package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"github.com/toolkits/net"
	"log"
	"sync"
)

type GlobalConfig struct {
	Debug    bool     `json:"debug"`
	LocalIp  string   `json:"localIp"`
	Servers  []string `json:"servers"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
	Docker   string   `json:"docker"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	if c.LocalIp == "" {
		// detect local ip
		localIps, err := net.IntranetIP()
		if err != nil {
			log.Fatalln("get intranet ip fail:", err)
		}

		if len(localIps) == 0 {
			log.Fatalln("no intranet ip found")
		}

		c.LocalIp = localIps[0]
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	if config.Debug {
		log.Println("read config file:", cfg, "successfully")
	}
}
