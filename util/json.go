package util

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

// JSONSerialize : golang object Serialized to json string
func JSONSerialize(obj interface{}) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(&obj)

	if err != nil {
		return fmt.Sprintf("{ \"code\" : -1, \"message\" : \"json marshl error:%s\"}", err.Error())
	}

	return string(data)
}

// JSONUnSerialize : Reverse the json string [byte] into a golang object
func JSONUnSerialize(data []byte, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, v)
	return err
}

// JSONUnFormSerialize : Reverse the json string into a golang object
func JSONUnFormSerialize(data string, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.UnmarshalFromString(data, v)
	return err
}
