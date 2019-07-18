package mgoImport

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Repository struct {
	Db    []map[string]interface{}
	Properties []Model
}

func NewRepository() *Repository {

	return &Repository{
		Properties: make([]Model,0),
		Db: make([]map[string]interface{},0),
	}

}

func (r *Repository) BuildProperties(cols_name []string, cols_type []string) error {

	if cols_name == nil || cols_type == nil {
		return errors.New("given slice is nil")
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
	return dataMap, nil
}

func setDataMapValue(fileType, fieldName, input string, data map[string]interface{}) (map[string]interface{}, error) {
	//fmt.Println("filetype: ",fileType, ",filename:",fieldName, ", input value:", input)
	var err error
	switch fileType {
	default:
		// string
		data[fieldName] = input
	case "int":
		val, _ := strconv.Atoi(input)
		data[fieldName] = val
	case "json":
		m := make(map[string]interface{})
		err = json.Unmarshal([]byte(input), &m)
		// merge json field into data map
		for k, v := range m {
			data[k] = v
		}
	}
	return data, err
}
