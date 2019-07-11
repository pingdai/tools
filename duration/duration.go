package duration

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type Duration struct {
	startTime time.Time
}

func (t *Duration) Reset() {
	t.startTime = time.Now()
}

func (t *Duration) Get() string {
	now := time.Now()
	duration := now.Sub(t.startTime)
	return fmt.Sprintf("%0.3f", float64(duration/time.Millisecond))
}

func (t *Duration) GetAndReset() string {
	defer t.Reset()
	return t.Get()
}

func (t *Duration) ToLogger() *logrus.Entry {
	return logrus.WithField("cost", t.Get())
}

func NewDuration() *Duration {
	return &Duration{
		startTime: time.Now(),
	}
}
