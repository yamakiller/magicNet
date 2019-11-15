package logger

const (
	EnvKey = "log"
)

//LogDeploy desc
//@struct LogDeploy desc logger deploy informat [json format]
type LogDeploy struct {
	LogPath  string `json:"log-path"`
	LogLevel int    `json:"log-level"`
	LogSize  int    `json:"log-size"`
}

//NewDefault desc
//@method NewDefault desc: create default logger deploy informat
func NewDefault() *LogDeploy {
	return &LogDeploy{"", int(TRACELEVEL), 1024}
}
