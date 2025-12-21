package utils

import (
	"bytes"
	"strings"

	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(in interface{}) (str string, err error) {
	var (
		buf []byte
	)
	buf, err = json.Marshal(in)
	if err != nil {
		return
	}
	str = string(buf)
	return
}

func ObjToByte(in interface{}) (buf []byte, err error) {
	buf, err = json.Marshal(in)
	return
}

func MustObjToByte(in interface{}) []byte {
	buf, err := json.Marshal(in)
	if err != nil {
		xlog.Errorf("MustObjToByte error: %v", err)
	}
	return buf
}

func ByteToObj(buf []byte, out interface{}) (err error) {
	dc := json.NewDecoder(bytes.NewReader(buf))
	dc.UseNumber()
	return dc.Decode(out)
}

func MustByteToMap(buf []byte) map[string]any {
	toMap, err := ByteToMap(buf)
	if err != nil {
		xlog.Errorf("MustByteToMap error: %v", err)
	}
	return toMap
}

func MustByteToMapSlice(buf []byte) []map[string]any {
	toMapSlice, err := ByteToMapSlice(buf)
	if err != nil {
		xlog.Errorf("MustByteToMapSlice error: %v", err)
		return []map[string]any{}
	}
	return toMapSlice
}

func ByteToMapSlice(buf []byte) ([]map[string]any, error) {
	var toMapSlice []map[string]any
	err := json.Unmarshal(buf, &toMapSlice)
	if err != nil {
		return nil, err
	}
	return toMapSlice, nil
}

func ByteToMap(buf []byte) (map[string]any, error) {
	var (
		maps map[string]any
		err  error
	)
	dc := json.NewDecoder(bytes.NewReader(buf))
	dc.UseNumber()
	if err = dc.Decode(&maps); err != nil {
		//fmt.Println(err)
	} else {
		for k, v := range maps {
			maps[k] = v
		}
	}
	return maps, err
}

func Unmarshal(in string, out interface{}) error {
	//return json.Unmarshal([]byte(in), out)
	dc := json.NewDecoder(strings.NewReader(in))
	dc.UseNumber()
	return dc.Decode(out)
}

func ObjToMap(in interface{}) map[string]interface{} {
	var (
		maps map[string]interface{}
		buf  []byte
		err  error
	)
	if buf, err = json.Marshal(in); err != nil {
		//fmt.Println(err)
	} else {
		d := json.NewDecoder(bytes.NewReader(buf))
		d.UseNumber()
		if err = d.Decode(&maps); err != nil {
			//fmt.Println(err)
		} else {
			for k, v := range maps {
				maps[k] = v
			}
		}
	}
	return maps
}

func MapToObj(maps map[string]interface{}, out interface{}) error {
	buf, err := json.Marshal(maps)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, out)
	if err != nil {
		return err
	}

	return nil
}

func MustMapToObj[T any](maps map[string]any) *T {
	buf, err := json.Marshal(maps)
	if err != nil {
		xlog.Errorf("MustMapToObj:%v", err)
		return nil
	}
	t := new(T)
	err = json.Unmarshal(buf, t)
	if err != nil {
		xlog.Errorf("MustMapToObj:%v", err)
		return nil
	}
	return t
}

func ObjToJsonStr(obj interface{}) (str string) {
	str = ""
	var err error
	if obj == nil {
		xlog.Errorf("ObjToJsonStr ERROR, obj is nil")
		return
	}
	str, err = Marshal(obj)
	if err != nil {
		xlog.Errorf("ObjToJsonStr:%v", err)
		return
	}
	return
}

func JsonStrToObj[T any](str string) (val *T, err error) {
	t := new(T)
	err = json.Unmarshal([]byte(str), t)
	return t, err
}
