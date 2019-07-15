package mgoImport

import (
	"encoding/json"
	"io"
)

//var G_config *ConfigFile

type ConfigFile struct {
	DataColumns []string `json:"data_columns"`
	DataTypes   []string `json:"data_types"`
	Delimiter   string   `json:"delimiter"`
}

func (cf *ConfigFile) LoadJson(data io.Reader) error {
	return json.NewDecoder(data).Decode(cf)
}
