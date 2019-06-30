package util

import (
	"sync"

	"github.com/tidwall/gjson"
)

var once sync.Once

type Env struct {
	_v map[string]gjson.Result
}

var instance *Env

// GetEnvInstance 单例模式获取全局 Env 对象
func GetEnvInstance() *Env {
	once.Do(func() {
		instance = new(Env)
	})
	return instance
}

// Put : 插入数据
func (E *Env) Put(v map[string]gjson.Result) {
	E._v = v
}

// Put 插入数据
/*func (E *Env) Put(k string, v interface{}) {
	E._v[k] = v
}*/

// GetInt : 获取k对映的整型值
func (E *Env) GetInt(k string, defaultValue int) int {
	v := E._v[k]
	if !v.Exists() {
		return defaultValue
	}
	return int(v.Int())
}

// GetString : 获取k对映的字符串
func (E *Env) GetString(k string, defaultValue string) string {
	v := E._v[k]
	if !v.Exists() {
		return defaultValue
	}

	return v.String()
}

// GetBoolean : 获取k对映的布尔值
func (E *Env) GetBoolean(k string, defaultValue bool) bool {
	v := E._v[k]
	if !v.Exists() {
		return defaultValue
	}
	return v.Bool()
}

// GetFloat : 获取k对映的 32位浮点数
func (E *Env) GetFloat(k string, defaultValue float32) float32 {
	v := E._v[k]
	if !v.Exists() {
		return defaultValue
	}

	return float32(v.Float())
}

// GetDouble : 获取k对映的64位浮点数
func (E *Env) GetDouble(k string, defaultValue float64) float64 {
	v := E._v[k]
	if !v.Exists() {
		return defaultValue
	}
	return v.Float()
}

// GetArray : 获取数组
func (E *Env) GetArray(k string) []gjson.Result {
	v := E._v[k]
	if !v.Exists() {
		return nil
	}
	return v.Array()
}

// GetMap : 获取子 MAP
func (E *Env) GetMap(k string) map[string]gjson.Result {
	v := E._v[k]
	if !v.Exists() {
		return nil
	}
	return v.Map()
}
