package util

// CMDLine :
type cmdLine struct {
	key map[string]interface{}
}

var (
	defaultCMDLine = cmdLine{make(map[string]interface{}, 8)}
)

// PushArgCmd ：
func PushArgCmd(k string, v interface{}) {
	defaultCMDLine.key[k] = v
}

// GetArgInt :
func GetArgInt(k string, def int) int {
	v := defaultCMDLine.key[k]
	if v == nil {
		return def
	}
	return v.(int)
}

// GetArgBool :
func GetArgBool(k string, def bool) bool {
	v := defaultCMDLine.key[k]
	if v == nil {
		return def
	}
	return v.(bool)
}

// GetArgString :
func GetArgString(k string, def string) string {
	v := defaultCMDLine.key[k]
	if v == nil {
		return def
	}
	return v.(string)
}
