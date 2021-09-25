package util

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

//JSONSerialize doc
//@Method JSONSerialize @Summary golang object Serialized to json string
//@Param  (interface{}) json object
//@Return (string) json string
func JSONSerialize(obj interface{}) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(&obj)

	if err != nil {
		return fmt.Sprintf("{ \"code\" : -1, \"message\" : \"json marshl error:%s\"}", err.Error())
	}

	return string(data)
}

//JSONUnSerialize doc
//@Method JSONUnSerialize @Summary Reverse the json string [byte] into a golang object
//@Param  ([]byte) json []byte
//@Param  (interface{}) out json object
//@Return (error)
func JSONUnSerialize(data []byte, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, v)
	return err
}

//JSONUnFormSerialize doc
//@Method JSONUnFormSerialize @Summary Reverse the json string into a golang object
//@Param  (string) json string
//@Param  (interface{}) out json object
func JSONUnFormSerialize(data string, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.UnmarshalFromString(data, v)
	return err
}
