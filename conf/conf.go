package conf

import (
	"github.com/pingdai/tools/reflectx"
	"reflect"
)

func UnmarshalConf(c interface{}) {
	rv := reflectx.Indirect(reflect.ValueOf(c))

	if !rv.CanSet() || rv.Type().Kind() != reflect.Struct {
		panic("UnmarshalConf need an variable which can set")
	}

	InitialRoot(rv)
}

type ICanInit interface {
	Init()
}

func InitialRoot(rv reflect.Value) {
	tpe := rv.Type()
	for i := 0; i < tpe.NumField(); i++ {
		value := rv.Field(i)
		if conf, ok := value.Interface().(ICanInit); ok {
			conf.Init()
		}
	}
}
