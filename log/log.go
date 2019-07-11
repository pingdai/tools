package log

import (
	"github.com/pingdai/tools/constants"
	"github.com/pingdai/tools/log/context"
	"github.com/pingdai/tools/log/hooks"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

type Log struct {
	Name   string
	Path   string
	Level  string `json:"level"`
	Format string // json or text
	init   bool
}

func (log *Log) MarshalDefaults() {

	if log.Name == "" {
		log.Name = os.Getenv(constants.EnvVarKeyProjectName)
	}

	if log.Level == "" {
		log.Level = constants.LOG_LEVEL_DEBUG
	}

	if log.Format == "" {
		log.Format = constants.LOG_FORMAT_TEXT
	}

	// test环境以上强制打印json格式
	envFlag := os.Getenv(constants.EnvVarKeyEnvFlag)
	if strings.ToLower(envFlag) == constants.ENV_TEST ||
		strings.ToLower(envFlag) == constants.ENV_PRE {
		log.Format = constants.LOG_FORMAT_JSON
	}
}

func (log *Log) Init() {
	if !log.init {
		log.MarshalDefaults()
		log.Create()
		log.init = true
	}
}

func (log *Log) Create() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(getLogLevel(log.Level))
	if log.Format == constants.LOG_FORMAT_JSON {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// log-id
	logrus.AddHook(context.NewLogIDHook())
	// file
	logrus.AddHook(hooks.NewCallStackHook())
	// project
	logrus.AddHook(hooks.NewProjectHook(log.Name))

	logrus.SetOutput(ioutil.Discard)

	if log.Path != "" {
		logrus.AddHook(hooks.NewLogWriterHook(log.Path))
		logrus.AddHook(hooks.NewLogWriterForErrorHook(log.Path))
	} else {
		logrus.AddHook(hooks.NewLogPrinterHook())
		logrus.AddHook(hooks.NewLogPrinterForErrorHook())
	}

}

func getLogLevel(l string) logrus.Level {
	level, err := logrus.ParseLevel(strings.ToLower(l))
	if err == nil {
		return level
	}
	return logrus.InfoLevel
}
