package servicex

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pingdai/tools/conf"
	"github.com/pingdai/tools/constants"
	"io/ioutil"
	"os"
	"reflect"
)

var ServiceName string

func GetZkConfPath() string {
	return fmt.Sprintf("/entry/config/service/%s", ServiceName)
}

func SetServiceName(serviceName string) {
	ServiceName = serviceName
}

// 解析结构体
func ConfP(in interface{}) {

	tpe := reflect.TypeOf(in)
	if tpe.Kind() != reflect.Ptr {
		panic(fmt.Errorf("ConfP pass ptr for setting value"))
	}

	bts := getConfigContent()

	printCfg := os.Getenv(constants.ENV_PRINT_CONFIG)
	if printCfg != "" {
		fmt.Printf("获取配置文件信息:%s\n", string(bts))
	} else {
		fmt.Printf("获取配置文件信息完成\n")
	}

	if err := json.Unmarshal(bts, in); err != nil {
		panic(fmt.Sprintf("json.Unmarshal conf err:%v", err))
	}

	// 进行解析
	os.Setenv(constants.EnvVarKeyProjectName, ServiceName)

	conf.UnmarshalConf(in)
}

func getConfigContent() []byte {
	var localCfgPath string
	//var zkHost = remoteCfgUrl
	var bts []byte
	var err error

	flag.StringVar(&localCfgPath, "c", "", "local config file path")
	flag.Parse()

	if localCfgPath != "" {
		// 进行本地文件解析
		bts, err = ioutil.ReadFile(localCfgPath)
		if err != nil {
			panic(fmt.Sprintf("Read local file[%s] err:%v", localCfgPath, err))
		}
	} else {
		// 从远端download文件进行解析
		// todo from zk or etcd or other
	}

	return bts
}
