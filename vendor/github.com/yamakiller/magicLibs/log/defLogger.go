package log

import (
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

//SpawnFileLogrus create an default logrus
//logPath log write file path [path]
//logName log write file base name
func SpawnFileLogrus(level logrus.Level,
	logPath string,
	logName string) (*logrus.Logger, error) {
	hlog := logrus.New()
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "white+h",
		TimestampStyle: "black+h"})
	hlog.SetFormatter(formatter)
	hlog.SetOutput(os.Stderr)
	hlog.SetLevel(level)

	if logPath == "" {
		return hlog, nil
	}

	if logName == "" {
		logName = "log"
	}

	baseLogPath := path.Join(logPath, logName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		return nil, err
	}

	fileformatter := new(prefixed.TextFormatter)
	fileformatter.FullTimestamp = true
	fileformatter.TimestampFormat = "2006-01-02 15:04:05"
	fileformatter.DisableColors = true

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, fileformatter)
	hlog.AddHook(lfHook)

	return hlog, nil
}
