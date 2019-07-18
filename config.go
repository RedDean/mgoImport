package mgoImport

import (
	"encoding/json"
	"io"
	"os"
)

//var G_config *ConfigFile

type ConfigFile struct {
	DataColumns []string `json:"data_columns"`
	DataTypes   []string `json:"data_types"`
	Delimiter   string   `json:"delimiter"`
}

func InitConfig(dir string) *ConfigFile {
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
