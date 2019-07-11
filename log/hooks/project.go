package hooks

import "github.com/sirupsen/logrus"

func NewProjectHook(name string) *ProjectHook {
	return &ProjectHook{
		Name: name,
	}
}

type ProjectHook struct {
	Name string
}

func (hook *ProjectHook) Fire(entry *logrus.Entry) error {
	entry.Data["project"] = hook.Name
	return nil
}

func (hook *ProjectHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
