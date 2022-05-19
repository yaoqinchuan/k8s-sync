package utils

import (
	"github.com/gogf/gf/v2/os/gcfg"
)

var ConfigData *gcfg.Config

func init() {
	ConfigData, _ = gcfg.New()
	// read path is /internal/manifest/config/config.yaml
	ConfigData.GetAdapter().(*gcfg.AdapterFile).SetPath("/internal/manifest/config/config.yaml")
}
