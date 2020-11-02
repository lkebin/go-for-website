package helper

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type tagOptions string

func BuildForInsert(modelPtr interface{}) (string, string, []interface{}, error) {
	columnValueMap, err := fieldProcessor(modelPtr)

	if err != nil {
		return "", "", nil, err
	}

	var (
		columns []string
		params  []string
		values  []interface{}
	)

	if len(columnValueMap) <= 0 {
		return "", "", nil, errors.New("empty data is given")
	}

	for k, v := range columnValueMap {
		columns = append(columns, k)
		params = append(params, "?")
		values = append(values, v)
	}

	return fmt.Sprintf("`%s`", strings.Join(columns, "`,`")), strings.Join(params, ","), values, nil
}

func BuildForUpdate(modelPtr interface{}) (string, []interface{}, error) {
	columnValueMap, err := fieldProcessor(modelPtr)

	if err != nil {
		return "", nil, err
	}

	var (
		columns []string
		values  []interface{}
	)

	if len(columnValueMap) <= 0 {
		return "", nil, errors.New("empty data is given")
	}

	for k, v := range columnValueMap {
		columns = append(columns, fmt.Sprintf("`%s`=?", k))
		values = append(values, v)
	}

	return strings.Join(columns, ","), values, nil
}

func fieldProcessor(modelPtr interface{}) (map[string]interface{}, error) {
	value := reflect.ValueOf(modelPtr)

	if value.Kind() != reflect.Ptr {
		return nil, errors.New("must pass a pointer, not a value")
	}

	if value.IsNil() {
		return nil, errors.New("nil pointer passed to field processor")
	}

	columnValueMap := make(map[string]interface{})

	v := reflect.Indirect(value)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		columnName, columnOption := parseTag(t.Field(i).Tag.Get("db"))
		if columnName == "" || columnName == "-" {
			continue
		}

		if columnOption.Contains("unsafe") {
			continue
		}

		columnValueMap[columnName] = v.Field(i).Interface()
	}

	return columnValueMap, nil
}

func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
