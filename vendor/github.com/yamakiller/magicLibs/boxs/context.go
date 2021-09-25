package boxs

import "github.com/yamakiller/magicLibs/actors"

//Context box 上下文
type Context struct {
	*actors.Context
	_funs []Method
	_idx  int
}

//Next 执行下一个关联函数
func (slf *Context) Next() {
	slf._idx++
	if slf._idx >= len(slf._funs) {
		return
	}

	slf._funs[slf._idx](slf)
}
