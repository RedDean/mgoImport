package mgoImport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

//var G_config *ConfigFile

const CONFIG_CONTENT_LENGTH = 4096

type ConfigFile struct {
	DataColumns []string          `json:"data_columns"`
	DataTypes   []string          `json:"data_types"`
	Delimiter   string            `json:"delimiter"`
	Db          DbConfig          `json:"db"`
	JsonField   map[string]string `json:"json_field"`

	ModifiedColumn string `json:"modified_column"`
	ID             IDconf `json:"id"`

	Enums []EnumNode `json:"enums"`
}

type DbConfig struct {
	Url        string `json:"url"`
	Name       string `json:"name"`
	Collection string `json:"collection"`
}

type IDconf struct {
	IdColumn      string   `json:"id_column"`
	RelatedColumn string   `json:"related_column"`
	ForeignColumn string   `json:"foreign_column"`
	Collections   []string `json:"collections"`
}

type EnumNode struct {
	CollectionName string `json:"collection_name,omitempty"`
	EnumColumn     string `json:"enum_column,omitempty"`
	OldValue       string `json:"old_value,omitempty"`
	NewValue       string `json:"new_value,omitempty"`
}

func InitConfig(normalDir, dbConfigDir string) *ConfigFile {
	fmt.Printf("[INFO] config dir is %s \n", normalDir)

	config := &ConfigFile{}
	if err := config.load(readFile(normalDir)); err != nil {
		panic(err)
	}

	if err := config.loadDBConf(readFile(dbConfigDir)); err != nil {
		panic(err)
	}

	return config
}

func readFile(dir string) *bytes.Buffer {
	file, err := os.OpenFile(dir, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	contentBytes := make([]byte, CONFIG_CONTENT_LENGTH)
	_, err = file.Read(contentBytes)
	if err != nil {
		panic(err)
	}

	return bytes.NewBuffer(contentBytes)
}

func (cf *ConfigFile) loadDBConf(data io.Reader) error {
	if err := json.NewDecoder(data).Decode(&cf.Db); err != nil {
		return err
	}

	fmt.Println("[debug] ", cf.Db)
	return nil
}

func (cf *ConfigFile) load(data io.Reader) error {
	return json.NewDecoder(data).Decode(cf)
}

func (cf *ConfigFile) GetIDConf() IDconf {
	return cf.ID
}

func (cf *ConfigFile) GetEnumNodeArray() []EnumNode {
	return cf.Enums
}
