package logger

import (
  "os"
  "time"
  "path"
  "github.com/pkg/errors"
  "github.com/lestrrat-go/file-rotatelogs"
  "github.com/rifflock/lfshook"
  "github.com/sirupsen/logrus"
  "gopkg.in/olivere/elastic.v5"
  "gopkg.in/sohlich/elogrus.v2"
  "github.com/vladoatanasov/logrus_amqp"
)

func InitLogger() {
  logrus.SetOutput(os.Stdout)
}

func ConfigLocalLogger(logPath string, logFileName string) {
  baseLogPath := path.Join(logPath, logFileName)
   writer, err := rotatelogs.New(
       baseLogPath + ".%Y%m%d%H%M",
       rotatelogs.WithLinkName(baseLogPath),
       rotatelogs.WithRotationCount(31),
       rotatelogs.WithRotationTime(time.Hour * 24))
   if err != nil {
     logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
     return
   }

   lfHook := lfshook.NewHook(lfshook.WriterMap{
        logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
        logrus.InfoLevel:  writer,
        logrus.WarnLevel:  writer,
        logrus.ErrorLevel: writer,
        logrus.FatalLevel: writer,
        logrus.PanicLevel: writer,
    }, &logrus.TextFormatter{DisableColors: true})
    logrus.AddHook(lfHook)
}

func ConfigRemoteAmqpLogger(server, username, password, exchange, exchangeType, virtualHost, routingKey string) {
  hook := logrus_amqp.NewAMQPHookWithType(server, username, password, exchange, exchangeType, virtualHost, routingKey)
  logrus.AddHook(hook)
}

func ConfigRemoteESLogger(esUrl string, esHOst string, index string) {
  client, err := elastic.NewClient(elastic.SetURL(esUrl))
    if err != nil {
        logrus.Errorf("config remote es logger error. %+v", errors.WithStack(err))
    }
    esHook, err := elogrus.NewElasticHook(client, esHOst, logrus.DebugLevel, index)
    if err != nil {
        logrus.Errorf("config remote es logger error. %+v", errors.WithStack(err))
    }
    logrus.AddHook(esHook)
}
