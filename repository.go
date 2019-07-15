package mgoImport

import (
	"errors"
)

func buildModelPropertiesMap(cols_name []string, cols_type []string) (map[string]string, error) {

	if cols_name == nil || cols_type == nil {
		return nil, errors.New("given slice is nil")
	}

	var (
		iName, jType int
		modelMap     map[string]string
	)

	if len(cols_name) != len(cols_type) {
		return nil, errors.New("length of cols_name doesn't equal length of cols_type")
	}

	modelMap = make(map[string]string)

	for iName < len(cols_name) && jType < len(cols_type) {
		modelMap[cols_name[iName]] = cols_type[jType]
		iName++
		jType++
	}

	return modelMap, nil
}
