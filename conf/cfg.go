package conf

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/walk1ng/gin-photo-gallery-storage/utils"
	"go.uber.org/zap"
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
		utils.AppLogger.Fatal(err.Error(), zap.String("service", "init()"))

	}

	ServerCfg.ConfigMap = make(map[string]string)
	err = json.NewDecoder(cfgFile).Decode(&ServerCfg.ConfigMap)
	if err != nil {
		utils.AppLogger.Fatal(err.Error(), zap.String("service", "init()"))
	}

}

// Get function
func (cfg *Cfg) Get(key string) string {
	if val, ok := cfg.ConfigMap[key]; !ok {
		return val
	}
	utils.AppLogger.Fatal(fmt.Sprintf("No such config term: %s", key), zap.String("service", "init()"))
	return ""
}
