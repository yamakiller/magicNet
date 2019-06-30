package monitor

import "net/http"

type MonitorMethod struct {
}

func (*MonitorMethod) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//r.RequestURI()
}
