package util

import (
  "fmt"
  "github.com/json-iterator/go"
)

func JsonSerialize(obj interface{}) string {
   var json = jsoniter.ConfigCompatibleWithStandardLibrary
   data, err := json.Marshal(&obj)
   if err != nil {
     return fmt.Sprintf("{ \"code\" : -1, \"message\" : \"json marshl error:%s\"}", err.Error())
   }

   return string(data)
}

func JsonUnSerialize(data []byte, v interface{}) error {
  var json = jsoniter.ConfigCompatibleWithStandardLibrary
  err := json.Unmarshal(data, v)
  return err
}
