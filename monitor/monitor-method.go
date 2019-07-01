package monitor

import (
  //"magicNet/logger"
  "magicNet/util"
  "regexp"
  "strings"
  "net/http"
)

type MonitorMethod struct {
  suffixRegexp *regexp.Regexp
  methods map[string] MonitorHttpFunction
}

type MonitorHttpFunction func(arg1 http.ResponseWriter, arg2 *http.Request)

func NewMonitorMethod() *MonitorMethod {
  r := &MonitorMethod{}
  r.suffixRegexp, _ = regexp.Compile(`\.\w+.*`)
  return r
}

func (mm *MonitorMethod) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  f := mm.match(r.RequestURI)
  if f == nil {
    w.WriteHeader(404)
    return
  }

  f(w, r)
}

func (mm *MonitorMethod) match(requestURI string) MonitorHttpFunction {
  suffix := mm.suffixRegexp.FindStringSubmatch(requestURI)
  if suffix != nil && len(suffix) >= 1 {
    return nil
  }

  tmpURI := requestURI
  idx := strings.LastIndex(requestURI, "?")
  if idx > 0 {
    tmpURI = util.SubStr2(tmpURI, 0, idx)
  }

  r := mm.methods[tmpURI]
  return r
}
