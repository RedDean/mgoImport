package mgoImport

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
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

	ModifiedColumn string
}

type dataModel struct {
	inner_json_cnt  int
	inner_json_type []string
	_map            map[string]interface{}
}

func InitRepository(c *ConfigFile, mode int) *Repository {
	fmt.Println("[DEBUG] dbname:", c.Db.Name)
	repo := NewRepository(c.Db.Name, c.Db.Collection, c.JsonField)

	switch mode {
	default:
		fmt.Printf("[ERROR] can't match program's mode type: %d while initliaze repository component! \n", mode)
		panic(errors.New("wrong mode type"))
	case NORMAL, ITEM, ITEM_HIS:
		if err := repo.buildProperties(c.DataColumns, c.DataTypes); err != nil {
			panic(err)
		}
	case MODIFY:
		if c.ModifiedColumn != "" {
			repo.ModifiedColumn = c.ModifiedColumn
		} else {
			panic(fmt.Errorf("[ERROR] column name is required! plz check configure file"))
		}
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

func (r Repository) BuildItemModel(input []string) (map[string]interface{}, error) {

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

		//if "description" == props.FieldName {
		//	// correct quote in string
		//	if "" != inputVal {
		//		inputVal = inputVal[2:len(inputVal)-2]
		//		inputVal = strings.ReplaceAll(inputVal, `"`,`\"`)
		//	}
		//}

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

func (r Repository) BuildModifyModel(input string) map[string]interface{} {

	//data := input[3:len(input)-3]
	// todo

	return bson.M{
		"$set": bson.M{
			r.ModifiedColumn: input[3 : len(input)-3],
		},
	}
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

			// In case that field in json override filed that has same name in dataMap.
			if _, ok := dm._map[k]; !ok {
				dm._map[k] = v
			}
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
		if input != "" {
			err = json.Unmarshal([]byte(input), &m)
		}
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

var G_item_rebuild_func_map = map[string]func(map[string]interface{}) map[string]interface{}{

	"APP": func(i map[string]interface{}) map[string]interface{} {

		bin, ok := i["binaries"].(map[string]interface{})
		if ok {
			i["binaries"] = resetBinaries(bin)
		}

		delete(i, "consumable")
		delete(i, "masterItemSlug")
		return i
	},

	"IAP": func(i map[string]interface{}) map[string]interface{} {
		delete(i, "packageName")
		delete(i, "categories")
		delete(i, "binaries")
		return i
	},
}

func resetBinaries(bin map[string]interface{}) interface{} {
	ext := bin["extensions"].(map[string]interface{})
	for key, value := range ext {
		if strings.HasSuffix(key, "apks") {
			bin[key] = resetApks(value.(map[string]interface{}))
		} else {
			// move field in 'extensions' to 'binaries'
			bin[key] = value
		}
	}

	delete(bin, "extensions")
	return bin
}

func resetApks(apk map[string]interface{}) []interface{} {
	if apk == nil {
		return nil
	}

	var apkSlice []interface{}

	for key, value := range apk {
		v := value.(map[string]interface{})
		v["versionCode"] = key
		apkSlice = append(apkSlice, v)
	}

	return apkSlice
}

func resetChannels(channels map[string]interface{}) interface{} {
	//ext := channels["extensions"]
	if channels == nil {
		return channels
	}

	for ck, value := range channels {
		ch := value.(map[string]interface{})
		ext := ch["extensions"].(map[string]interface{})
		for ek, ev := range ext {
			ch[ek] = ev
		}
		delete(ch, "extensions")
		channels[ck] = ch
	}

	return channels
}
