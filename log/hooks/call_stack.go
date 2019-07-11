package hooks

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func NewCallStackHook() *CallStackHook {
	return &CallStackHook{}
}

type CallStackHook struct {
}

func (hook *CallStackHook) Fire(entry *logrus.Entry) error {
	_, file, line := findCaller()
	entry.Data["file"] = fmt.Sprintf("%s:%d", file, line)
	return nil
}

func (hook *CallStackHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func findCaller() (inFunc string, file string, line int) {
	pc := uintptr(0)
	skip := 5

	for i := 0; i < 10; i++ {
		pc, file, line = getCaller(skip + i)
		if !(strings.HasPrefix(file, "logrus")) {
			inFunc = runtime.FuncForPC(pc).Name()
			break
		}
	}

	return
}

func getGoPkgName(file string) string {
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n += 1
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file
}

func getCaller(skip int) (uintptr, string, int) {
	pc, file, line, ok := runtime.Caller(skip)

	if !ok {
		return 0, "", 0
	}

	return pc, getGoPkgName(file), line
}
