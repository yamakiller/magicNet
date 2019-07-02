package logger

/*const (
	logDefaultMode    int = 0
	logLocalFileMode  int = 1
	logRemoteAmqpMode int = 2
	logRemoteEsMode   int = 3
)*/

type RemoteLogger struct {
}

// Redirect is Config log mode and level
/*func Redirect() {
	Info(0, "redirect log level %s", util.GetEnvInstance().GetString("log-level", "panic"))
	configLevel(util.GetEnvInstance().GetString("log-level", "panic"))
	Info(0, "redirect log mode %s", util.GetEnvInstance().GetString("log-mode", "default"))

	switch parselMode(util.GetEnvInstance().GetString("log-mode", "default")) {
	case logLocalFileMode:
		logFilePath := util.GetEnvInstance().GetString("log-file-path", "./")
		logFileName := util.GetEnvInstance().GetString("log-file-name", "")
		Info(0, "redirect log local file mode, path:%s file:%s", logFilePath, logFileName)
		configLocalLogger(logFilePath,
			logFileName)
		break
	case logRemoteAmqpMode:
		logAmqpServer := util.GetEnvInstance().GetString("log-amqp-server", "")
		logAmqpAccount := util.GetEnvInstance().GetString("log-amqp-account", "")
		logAmqpPassword := util.GetEnvInstance().GetString("log-amqp-password", "")
		logAmqpExchange := util.GetEnvInstance().GetString("log-amqp-exchange", "")
		logAmqpExchangeType := util.GetEnvInstance().GetString("log-amqp-exchange-type", "")
		logAmqpVirtualHost := util.GetEnvInstance().GetString("log-amqp-virtual-host", "")
		logAmqpRoutingKey := util.GetEnvInstance().GetString("log-amqp-routing-key", "")
		Info(0, "redirect log remote amqp mode, server:%s, account:%s, password:%s, exchange:%s, exchange-type:%s, virtual-host:%s, routing-key:%s",
			logAmqpServer,
			logAmqpAccount,
			logAmqpPassword,
			logAmqpExchange,
			logAmqpExchangeType,
			logAmqpVirtualHost,
			logAmqpRoutingKey)

		configRemoteAmqpLogger(logAmqpServer,
			logAmqpAccount,
			logAmqpPassword,
			logAmqpExchange,
			logAmqpExchangeType,
			logAmqpVirtualHost,
			logAmqpRoutingKey)
		break
	case logRemoteEsMode:
		logEsURL := util.GetEnvInstance().GetString("log-es-url", "")
		logEsHost := util.GetEnvInstance().GetString("log-es-host", "")
		logEsIndex := util.GetEnvInstance().GetString("log-es-index", "")
		Info(0, "redirect log remote es mode, es-url:%s, es-host:%s, es-index:%s",
			logEsURL,
			logEsHost,
			logEsIndex)

		configRemoteESLogger(logEsURL,
			logEsHost,
			logEsIndex)
		break
	default:
		break
	}
}

func parselMode(mode string) int {
	switch mode {
	case "local":
		return logLocalFileMode
	case "amqp":
		return logRemoteAmqpMode
	case "es":
		return logRemoteEsMode
	default:
		return logDefaultMode
	}
}*/

/*func configLocalLogger(logPath string, logFileName string) {
	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),
		rotatelogs.WithRotationCount(31),
		rotatelogs.WithRotationTime(time.Hour*24))
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
		logrus.TraceLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true})
	logrus.AddHook(lfHook)
}

func configRemoteAmqpLogger(server, username, password, exchange, exchangeType, virtualHost, routingKey string) {
	hook := logrus_amqp.NewAMQPHookWithType(server, username, password, exchange, exchangeType, virtualHost, routingKey)
	logrus.AddHook(hook)
}

func configRemoteESLogger(esURL string, esHost string, index string) {
	client, err := elastic.NewClient(elastic.SetURL(esURL))
	if err != nil {
		logrus.Errorf("config remote es logger error. %+v", errors.WithStack(err))
	}
	esHook, err := elogrus.NewElasticHook(client, esHost, logrus.DebugLevel, index)
	if err != nil {
		logrus.Errorf("config remote es logger error. %+v", errors.WithStack(err))
	}
	logrus.AddHook(esHook)
}*/
