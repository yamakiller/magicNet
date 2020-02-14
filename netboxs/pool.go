package netboxs

//Pool connection pool interface
type Pool interface {
	Get() Connect
	Put(Connect)
}
