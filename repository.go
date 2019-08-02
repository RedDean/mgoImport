package mgoImport

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Repository struct {
	JsonField  map[string]string
	DbName     string
	Collection string
	Properties []Model
}

type dataModel struct {
	inner_json_cnt  int
	inner_json_type []string
	_map            map[string]interface{}
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

	if len(cols_name) != len(cols_type) {
		return errors.New("length of cols_name doesn't equal length of cols_type")
	}

	r.Properties = make([]Model, len(cols_name))
	for i := 0; i < len(cols_name); i++ {
		r.Properties[i] = Model{
			FieldType: cols_type[i],
			FieldName: cols_name[i],
		}
	}

	return nil
}

func (r Repository) BuildModel(input []string) (map[string]interface{}, error) {

	defer func() {
		if rc := recover(); rc != nil {
			fmt.Printf("[DEBUG] recordId:%s input: %d, props: %d ,maybe data syntx has wrong type! \n", input[0], len(input), len(r.Properties))
		}
	}()

	var (
		err     error
		dataMap dataModel
	)

	dataMap._map = make(map[string]interface{})

	for i := range input {
		inputVal := input[i]
		props := r.Properties[i]
		if err = dataMap.setDataMapValue(props.FieldType, props.FieldName, inputVal); err != nil {
			return nil, err
		}
	}

	// reset json field type
	for fName, fType := range r.JsonField {
		if err = dataMap.setDataMapValue(fType, fName, decodeJsonInterface(dataMap._map[fName])); err != nil {
			return nil, err
		}
	}

	return dataMap._map, nil
}

func (dm *dataModel) setDataMapValue(fileType, fieldName, input string) error {
	var err error
	switch fileType {
	default:
		// string
		dm._map[fieldName] = input

	case "int":
		var val int
		if input == "" {
			dm._map[fieldName] = 0
			return nil
		}

		val, err = strconv.Atoi(input)
		if err != nil {
			return err
		}
		dm._map[fieldName] = val

	case "json":
		if input == "" {
			return nil
		}
		m := make(map[string]interface{})
		err = json.Unmarshal([]byte(input), &m)
		if err != nil {
			fmt.Printf("[DEBUG] input string:%s ,err:%v \n", input, err)
		}
		// merge json field into data map
		for k, v := range m {
			dm._map[k] = v
		}

	case "date":
		var t time.Time
		t, err = time.Parse("2006-01-02 15:04:05+00", input)
		dm._map[fieldName] = t

	case "bool":
		var b bool
		if input == "t" {
			b = true
		}
		dm._map[fieldName] = b

	case "nested-json":
		m := make(map[string]interface{})
		err = json.Unmarshal([]byte(input), &m)
		dm._map[fieldName] = m

	case "[]string":
		dm._map[fieldName] = splitStr(input)

	case "[]int":
		strArr := splitStr(input)
		if len(strArr) == 0 {
			dm._map[fieldName] = []int{}
			return nil
		}
		retArr := make([]int, len(strArr))
		for key := range strArr {
			v := strArr[key]
			if vv, err := strconv.Atoi(v); err != nil {
				return err
			} else {
				retArr[key] = vv
			}
		}
		dm._map[fieldName] = retArr
		return nil

	}

	return err
}

func (dm *dataModel) getInnerJsonType() string {
	var idx int
	if dm.inner_json_type == nil || len(dm.inner_json_type) == 0 {
		return ""
	}
	if dm.inner_json_cnt > len(dm.inner_json_type)-1 {
		return ""
	}
	idx = dm.inner_json_cnt
	dm.inner_json_cnt++

	return dm.inner_json_type[idx]
}

func splitStr(input string) []string {
	var ret []string
	if strings.TrimSpace(input) == "" { // psql 字段为null值
		return ret
	}
	input = input[1 : len(input)-1]
	return strings.Split(input, ",")
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
	default:
		ret = vv.String()
	}

	return ret
}
