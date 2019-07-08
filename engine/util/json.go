package util

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

// JSONSerialize : golang object 序列化为 json字符串
func JSONSerialize(obj interface{}) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(&obj)
	if err != nil {
		return fmt.Sprintf("{ \"code\" : -1, \"message\" : \"json marshl error:%s\"}", err.Error())
	}

	return string(data)
}

// JSONUnSerialize : 把json字符串反列化为 golang object
func JSONUnSerialize(data []byte, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, v)
	return err
}
