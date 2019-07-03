package util

import (
	"reflect"
	"os"
	"io/ioutil"
	"magicNet/engine/logger"
)


var (
	instEnv map[string]interface{}
)

func LoadEnv(filename string) int {
	f, err := os.Open(filename)
	if err != nil {
		logger.Error(0, "open env config fail:%s", filename)
		return -1
	}
	defer f.Close()
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Error(0, "read env config fail:%s", err.Error())
		return -1
	}

	instEnv = make(map[string]interface{})
	err = JsonUnSerialize(contents, &instEnv)
	if err != nil {
		instEnv = nil
		logger.Error(0, "env unserialize fail:%s", err.Error())
		return -1
	}

	return 0
}

func UnLoadEnv() {
	instEnv = nil
}

func GetEnvRoot() map[string]interface{} {
	return instEnv
}

func GetEnvMap(v map[string]interface{}, k string) map[string]interface{} {
	if (v[k] == nil ) {
			return nil
	}

	elem := reflect.ValueOf(v[k])
	if elem.IsNil() ||
		 elem.Kind() != reflect.Map{
		return nil
	}

	var outv map[string]interface{}
	var inv interface{} = &outv
	reflect.ValueOf(inv).Elem().Set(elem)

	return outv
}

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

func ToEnvMap(v interface{}) map[string]interface{} {
	var outv map[string]interface{}
	var inv interface{} = &outv
	reflect.ValueOf(inv).Elem().Set(reflect.ValueOf(v))
	return outv
}

func GetEnvBoolean(v map[string]interface{}, k string, defaultValue bool) bool {
	istr := getEnvValue(v, k, reflect.Bool)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).Bool()
}

func GetEnvString(v map[string]interface{}, k string, defaultValue string) string {
	istr := getEnvValue(v, k, reflect.String)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).String()
}

func GetEnvInt(v map[string]interface{}, k string, defaultValue int) int {
	istr := getEnvValue(v, k, reflect.Int)
	if istr == nil {
		return defaultValue
	}

	return int(reflect.ValueOf(istr).Int())
}

func GetEnvInt64(v map[string]interface{}, k string, defaultValue int64) int64 {
	istr := getEnvValue(v, k, reflect.Int64)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).Int()
}

func GetEnvFloat(v map[string]interface{}, k string, defaultValue float32) float32 {
	istr := getEnvValue(v, k, reflect.Float32)
	if istr == nil {
		return defaultValue
	}

	return float32(reflect.ValueOf(istr).Float())
}

func GetEnvDouble(v map[string]interface{}, k string, defaultValue float64) float64 {
	istr := getEnvValue(v, k, reflect.Float64)
	if istr == nil {
		return defaultValue
	}

	return reflect.ValueOf(istr).Float()
}


func getEnvValue(v map[string]interface{}, k string, c reflect.Kind) interface {} {
	ival := v[k]
	if ival == nil ||
		 reflect.ValueOf(ival).Kind() != c {
		return nil
	}

	return ival
}
