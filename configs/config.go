package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// LoadMysqlConfig mysql json配置转换为 map
func LoadMysqlConfig(filename string) (map[string]interface{}, error) {
	file, err := os.Open("configs/" + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
