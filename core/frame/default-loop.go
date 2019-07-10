package frame

type defaultLoop struct {
}

func (dp *defaultLoop) Wait() int {

	return -1
}
