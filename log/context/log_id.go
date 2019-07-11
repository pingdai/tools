package context

import (
	"runtime"
	"sync"
)

var (
	logIDs = make(map[int64]string)
	locker = sync.RWMutex{}
)

func getGoID() int64 {
	return runtime.GoID()
}

func SetLogID(id string) {
	locker.Lock()
	defer locker.Unlock()

	logIDs[getGoID()] = id
}

func GetLogID() string {
	locker.RLock()
	defer locker.RUnlock()

	goID := getGoID()

	if logID, ok := logIDs[goID]; ok {
		return logID
	}

	return ""
}
