package engine

import (
	"io/ioutil"
	"magicNet/logger"
	"magicNet/util"
	"os"

	"github.com/tidwall/gjson"
)

// LoadEnv read json config inforamt
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

	js := gjson.Parse(string(contents))
	util.GetEnvInstance().Put(js.Map())
	return 0
}
