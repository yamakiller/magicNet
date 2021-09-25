package actors

//Actor 对象
type Actor interface {
	Receive(*Context)
}
