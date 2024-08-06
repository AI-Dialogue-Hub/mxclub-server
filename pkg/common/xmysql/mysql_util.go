package xmysql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"mxclub/pkg/utils"
)

type JSONArray []map[string]any

// Value 这里不要改成 *JSONArray
// @see:https://gorm.io/zh_CN/docs/data_types.html
func (j JSONArray) Value() (driver.Value, error) {
	return utils.ObjToJsonStr(j), nil
}

func (j *JSONArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok || len(bytes) == 0 {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONArray value:", value))
	}
	data, err := utils.ByteToMapSlice(bytes)
	if err != nil {
		return err
	}
	*j = data
	return nil
}

// StringArray type to handle StringArray encoding/decoding
type StringArray []string

func (j *StringArray) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), j)
}

func (j StringArray) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// JSON type to handle StringArray encoding/decoding
type JSON map[string]any

func (j *JSON) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), j)
}

func (j JSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}
