package conf

import (
	"encoding/json"
	"log"
	"os"
)

// Cfg struct
type Cfg struct {
	ConfigMap map[string]string
}

// ServerCfg is global config for server
var ServerCfg Cfg

func init() {
	cfgFile, err := os.Open("conf/server.json")
	defer cfgFile.Close()
	if err != nil {
		log.Fatalln(err)
	}

	ServerCfg.ConfigMap = make(map[string]string)
	err = json.NewDecoder(cfgFile).Decode(&ServerCfg.ConfigMap)
	if err != nil {
		log.Fatalln(err)
	}

}

// Get function
func (cfg *Cfg) Get(key string) string {
	if val, ok := cfg.ConfigMap[key]; !ok {
		return val
	}
	log.Fatalln("No such config term: %s!", key)
	return ""
}
