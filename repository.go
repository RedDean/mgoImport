package mgoImport

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Repository struct {
	JsonField  map[string]string
	DbName     string
	Collection string
	Properties []Model
}

func InitRepository(c *ConfigFile) *Repository {
	repo := NewRepository(c.Db.Name, c.Db.Collection, c.JsonField)
	if err := repo.buildProperties(c.DataColumns, c.DataTypes); err != nil {
		panic(err)
	}
	return repo
}

func NewRepository(db string, collection string, jsonField map[string]string) *Repository {
	return &Repository{
		Properties: make([]Model, 0),
		DbName:     db,
		Collection: collection,
		JsonField:  jsonField,
	}
}

func (r *Repository) buildProperties(cols_name []string, cols_type []string) error {

	if cols_name == nil || cols_type == nil {
		return errors.New("given colName || colType slice is nil")
	}

	var i int // index

	if len(cols_name) != len(cols_type) {
		return errors.New("length of cols_name doesn't equal length of cols_type")
	}

	r.Properties = make([]Model, len(cols_name))
	for ; i < len(cols_name); i++ {
		r.Properties[i] = Model{
			FieldType: cols_type[i],
			FieldName: cols_name[i],
		}
	}

	return nil
}

func (r Repository) BuildModel(input []string) (map[string]interface{}, error) {

	var err error
	dataMap := make(map[string]interface{})

	for i := range input {
		inputVal := input[i]
		props := r.Properties[i]
		dataMap, err = setDataMapValue(props.FieldType, props.FieldName, inputVal, dataMap)
		if err != nil {
			return nil, err
		}
	}

	// reset json field type
	for fName, fType := range r.JsonField {
		dataMap, err = setDataMapValue(fType, fName, decodeJsonInterface(dataMap[fName]), dataMap)
		if err != nil {
			return nil, err
		}
	}

	return dataMap, nil
}

func setDataMapValue(fileType, fieldName, input string, data map[string]interface{}) (map[string]interface{}, error) {
	var err error
	switch fileType {
	default:
		// string
		data[fieldName] = input

	case "int":
		var val int
		if input == "" {
			data[fieldName] = 0
			return data, nil
		}

		val, err = strconv.Atoi(input)
		if err != nil {
			return data, err
		}
		data[fieldName] = val

	case "json":
		m := make(map[string]interface{})
		err = json.Unmarshal([]byte(input), &m)
		// merge json field into data map
		for k, v := range m {
			data[k] = v
		}

	case "date":
		var t time.Time
		t, err = time.Parse("2006-01-02 15:04:05+00", input)
		data[fieldName] = t

	case "bool":
		var b bool
		if input == "t" {
			b = true
		}
		data[fieldName] = b
	}

	return data, err
}

func decodeJsonInterface(i interface{}) string {
	var ret string
	vv := reflect.ValueOf(i)
	if !vv.IsValid() {
		return ""
	}

	switch vv.Type().String() {
	case "float64": // json中数值类型默认转成float，当指定json类型时一般是将数值类型转成int，因此此处将其先转成int类型
		f := vv.Float()
		ret = fmt.Sprintf("%d", int(f))
	}

	return ret
}
