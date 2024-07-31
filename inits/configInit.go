package inits

import (
	"forum/pkg/utils"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// ConfigInit 初始化配置
func ConfigInit() {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&(utils.AppConfig)); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}
}
