package hooks

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"sync"
)

// go get gopkg.in/natefinch/lumberjack.v2
func NewLogWriterHook(p string) *LogWriterHook {
	return &LogWriterHook{
		lock: new(sync.Mutex),
		levels: []logrus.Level{
			logrus.DebugLevel,
			logrus.InfoLevel,
		},
		writer: &lumberjack.Logger{
			Filename:   p,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     30, //days
		},
	}
}

func NewLogWriterForErrorHook(p string) *LogWriterHook {
	return &LogWriterHook{
		lock: new(sync.Mutex),
		levels: []logrus.Level{
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		},
		writer: &lumberjack.Logger{
			Filename:   p + ".wf",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     30, //days
		},
	}
}

type LogWriterHook struct {
	levels []logrus.Level
	writer io.Writer
	lock   *sync.Mutex
}

func (hook *LogWriterHook) Fire(entry *logrus.Entry) error {
	hook.lock.Lock()
	defer hook.lock.Unlock()

	msg, err := entry.String()
	if err != nil {
		return err
	} else {
		hook.writer.Write([]byte(msg))
	}

	return nil
}

func (hook *LogWriterHook) Levels() []logrus.Level {
	return hook.levels
}
