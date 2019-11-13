package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/yamakiller/magicNet/engine/files"
)

var (
	instEnv map[string]interface{}
)

//LoadMapEnv desc
//@method LoadMapEnv desc : load json file map to env
//@param (string) json file path
//@param (interface{}) map to env[struct]
//@return (error) a error message
func LoadMapEnv(filename string, out interface{}) error {
	fullpath := files.GetFullPathForFilename(filename)
	f, err := os.Open(fullpath)
	if err != nil {
		return fmt.Errorf("open env map file fail:%s", fullpath)
	}

	defer f.Close()
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read env map file fail:%s", err.Error())
	}

	err = json.Unmarshal(contents, out)
	if err != nil {
		return fmt.Errorf("env map fail:%s", err.Error())
	}
	return nil
}

//LoadEnv desc
//@method LoadEnv desc: load json file to global env
//@param (string) json file path
//@return (error) a error message
func LoadEnv(filename string) error {
	fullpath := files.GetFullPathForFilename(filename)
	f, err := os.Open(fullpath)
	if err != nil {
		return fmt.Errorf("open env config fail:%s", fullpath)
	}
	defer f.Close()
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read env config fail:%s", err.Error())
	}

	instEnv = make(map[string]interface{})
	err = JSONUnSerialize(contents, &instEnv)
	if err != nil {
		instEnv = nil
		return fmt.Errorf("env unserialize fail:%s", err.Error())
	}

	return nil
}

//UnLoadEnv desc
//@method UnLoadEnv desc: unload global env
func UnLoadEnv() {
	instEnv = nil
}

//GetEnvRoot desc
//@method GetEnvRoot desc: return global env root
//@return (map[string]interface{})
func GetEnvRoot() map[string]interface{} {
	return instEnv
}

//GetEnvMap desc
//@method GetEnvMap desc: return key=>value
//@param (map[string]interface{}) source map
//@param (string) key
//@return (map[string]interface{})
func GetEnvMap(v map[string]interface{}, k string) map[string]interface{} {
	if v[k] == nil {
		return nil
	}

	elem := reflect.ValueOf(v[k])
	if elem.IsNil() ||
		elem.Kind() != reflect.Map {
		return nil
	}

	var outv map[string]interface{}
	var inv interface{} = &outv
	reflect.ValueOf(inv).Elem().Set(elem)

	return outv
}

//GetEnvArray desc
//@method GetEnvArray desc: return a array
//@param (map[string]interface{}) source map
//@param (string) key
//@return a arrays
func GetEnvArray(v map[string]interface{}, k string) []interface{} {
	if v[k] == nil {
		return nil
	}

	elem := reflect.ValueOf(v[k])

	if elem.Kind() != reflect.Slice &&
		elem.Kind() != reflect.Array {
		return nil
	}

	var outv []interface{}
	var inv interface{} = &outv

	reflect.ValueOf(inv).Elem().Set(elem)

	return outv
}

//ToEnvMap desc
//@method ToEnvMap desc: return key=>value
//@param (interface{}) source object
//@return (map[string]interface{})
func ToEnvMap(v interface{}) map[string]interface{} {
	var outv map[string]interface{}
	var inv interface{} = &outv
	reflect.ValueOf(inv).Elem().Set(reflect.ValueOf(v))
	return outv
}

//GetEnvBoolean desc
//@method GetEnvBoolean desc: return key=>boolean
//@param (map[string]interface{}) source map
//@param (string) key
//@param (bool) default value
//@return (bool) a boolean value
func GetEnvBoolean(v map[string]interface{}, k string, defaultValue bool) bool {
	istr := getEnvValue(v, k, reflect.Bool)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).Bool()
}

//GetEnvString desc
//@method GetEnvString desc: return key=>string
//@param (map[string]interface{}) source map
//@param (string) key
//@param (string) default value
//@return (string) return a string
func GetEnvString(v map[string]interface{}, k string, defaultValue string) string {
	istr := getEnvValue(v, k, reflect.String)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).String()
}

//GetEnvInt desc
//@method  GetEnvInt desc: return key=>int
//@param   (map[string]interface{}) source map
//@param   (string) key
//@param   (int) default value
//@return  (int) return a int
func GetEnvInt(v map[string]interface{}, k string, defaultValue int) int {
	istr := getEnvValue(v, k, reflect.Int)
	if istr == nil {
		return defaultValue
	}

	return int(reflect.ValueOf(istr).Int())
}

//GetEnvInt64 desc
//@method GetEnvInt64 desc: return key=>int64
//@param  (map[string]interface{}) source map
//@param  (string) key
//@param  (int64) default value
//@return (int64) return a int64
func GetEnvInt64(v map[string]interface{}, k string, defaultValue int64) int64 {
	istr := getEnvValue(v, k, reflect.Int64)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).Int()
}

//GetEnvFloat desc
//@method GetEnvFloat desc: return key=>float32
//@param  (map[string]interface{}) source map
//@param  (string) key
//@param  (float32) default value
//@return (float32) return a float32
func GetEnvFloat(v map[string]interface{}, k string, defaultValue float32) float32 {
	istr := getEnvValue(v, k, reflect.Float32)
	if istr == nil {
		return defaultValue
	}

	return float32(reflect.ValueOf(istr).Float())
}

//GetEnvDouble desc
//@method GetEnvDouble desc: return key=>float64
//@param  (map[string]interface{}) source map
//@param  (string) key
//@param  (float64) default value
//@return (float64) return a float64
func GetEnvDouble(v map[string]interface{}, k string, defaultValue float64) float64 {
	istr := getEnvValue(v, k, reflect.Float64)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).Float()
}

func getEnvValue(v map[string]interface{}, k string, c reflect.Kind) interface{} {
	ival := v[k]
	if ival == nil ||
		reflect.ValueOf(ival).Kind() != c {
		return nil
	}

	return ival
}
