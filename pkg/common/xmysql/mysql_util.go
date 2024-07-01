package xmysql

import (
	"database/sql/driver"
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
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	data, err := utils.ByteToMapSlice(bytes)
	if err != nil {
		return err
	}
	*j = data
	return nil
}
