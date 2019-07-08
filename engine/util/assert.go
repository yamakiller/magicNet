package util

// Assert :  断言Bool并输出错误信息
func Assert(isAs bool, errMsg string) {
	if !isAs {
		panic(errMsg)
	}
}

// AssertEmpty : 断言Nil并输出错误信息
func AssertEmpty(isNull interface{}, errMsg string) {
	if isNull == nil {
		panic(errMsg)
	}
}
