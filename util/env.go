package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/yamakiller/magicNet/engine/files"
)

var (
	instEnv map[string]interface{}
)

// LoadEnv : 载入环境变量配置信息
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

// UnLoadEnv : 卸载载入的环境变量
func UnLoadEnv() {
	instEnv = nil
}

// GetEnvRoot : 获取环境变量根map对象
func GetEnvRoot() map[string]interface{} {
	return instEnv
}

// GetEnvMap : 获取当前环境变量 k -> 的map接点
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

// GetEnvArray : 获取当前环境变量 k -> 数组对象
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

// ToEnvMap : 获取当前环境变的 map 对象
func ToEnvMap(v interface{}) map[string]interface{} {
	var outv map[string]interface{}
	var inv interface{} = &outv
	reflect.ValueOf(inv).Elem().Set(reflect.ValueOf(v))
	return outv
}

// GetEnvBoolean : 获取当前环境变量 k-> Bool
func GetEnvBoolean(v map[string]interface{}, k string, defaultValue bool) bool {
	istr := getEnvValue(v, k, reflect.Bool)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).Bool()
}

// GetEnvString : 获取当前环境变量 k -> String
func GetEnvString(v map[string]interface{}, k string, defaultValue string) string {
	istr := getEnvValue(v, k, reflect.String)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).String()
}

// GetEnvInt : 获取当前环境变量 k -> Int
func GetEnvInt(v map[string]interface{}, k string, defaultValue int) int {
	istr := getEnvValue(v, k, reflect.Int)
	if istr == nil {
		return defaultValue
	}

	return int(reflect.ValueOf(istr).Int())
}

// GetEnvInt64 : 获取当前环境变量 k -> int64
func GetEnvInt64(v map[string]interface{}, k string, defaultValue int64) int64 {
	istr := getEnvValue(v, k, reflect.Int64)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).Int()
}

// GetEnvFloat : 获取当前环境变量 k -> float32
func GetEnvFloat(v map[string]interface{}, k string, defaultValue float32) float32 {
	istr := getEnvValue(v, k, reflect.Float32)
	if istr == nil {
		return defaultValue
	}

	return float32(reflect.ValueOf(istr).Float())
}

// GetEnvDouble : 获取当前环境变量 k -> float64
func GetEnvDouble(v map[string]interface{}, k string, defaultValue float64) float64 {
	istr := getEnvValue(v, k, reflect.Float64)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).Float()
}

// GetEnvValue : 获取当前环境变量 k -> interface
func getEnvValue(v map[string]interface{}, k string, c reflect.Kind) interface{} {
	ival := v[k]
	if ival == nil ||
		reflect.ValueOf(ival).Kind() != c {
		return nil
	}

	return ival
}
