package context

import "github.com/sirupsen/logrus"

func NewLogIDHook() *LogIDHook {
	return &LogIDHook{}
}

type LogIDHook struct {
}

func (hook *LogIDHook) Fire(entry *logrus.Entry) error {
	entry.Data["log_id"] = GetLogID()
	return nil
}

func (hook *LogIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
