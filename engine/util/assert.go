package util

func Assert(isAs bool, errMsg string) {
  if !isAs {
    panic(errMsg)
  }
}

func AssertEmpty(isNull interface{}, errMsg string) {
  if isNull == nil {
    panic(errMsg)
  }
}
