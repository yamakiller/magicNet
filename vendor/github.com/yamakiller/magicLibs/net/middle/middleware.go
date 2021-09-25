package middle

//Exception 中间异常处理接口
type Exception interface {
	Error(error)
	Debug(error)
}
