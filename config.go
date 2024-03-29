package mgoImport

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

//var G_config *ConfigFile

type ConfigFile struct {
	DataColumns []string          `json:"data_columns"`
	DataTypes   []string          `json:"data_types"`
	Delimiter   string            `json:"delimiter"`
	Db          DbConfig          `json:"db"`
	JsonField   map[string]string `json:"json_field"`
}

type DbConfig struct {
	Url        string `json:"url"`
	Name       string `json:"name"`
	Collection string `json:"collection"`
}

func InitConfig(dir string) *ConfigFile {
	fmt.Printf("[INFO] config dir is %s \n", dir)
	file, err := os.OpenFile(dir, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	config := &ConfigFile{}
	if err := config.LoadJson(file); err != nil {
		panic(err)
	}
	return config
}

func (cf *ConfigFile) LoadJson(data io.Reader) error {
	return json.NewDecoder(data).Decode(cf)
}
