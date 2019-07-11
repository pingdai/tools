package ginx

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pingdai/tools/constants"
	"os"
	"time"
)

type Ginx struct {
	Engine     *gin.Engine
	ListenPort int    `json:"listen_port"`
	GinMode    string `json:"gin_mode"`
	init       bool
}

func (ginx *Ginx) Init() {
	if !ginx.init {
		ginx.New()
		ginx.init = true
	}
}

func (ginx *Ginx) New() {
	ginx.Engine = gin.New()

	var ginMode = gin.DebugMode
	if os.Getenv(constants.I_EnvVarKeyGinMode) != "" {
		ginMode = os.Getenv(constants.I_EnvVarKeyGinMode)
	}
	if ginx.GinMode != "" {
		ginMode = ginx.GinMode
	}

	gin.SetMode(ginMode)
	os.Setenv(constants.I_EnvVarKeyGinMode, ginMode)

	// 添加一些中间件
	ginx.Engine.Use(gin.Recovery())
	ginx.Engine.Use(GinLogger())

	// TODO
}

// addr = [:port] e.g. :1500
func (ginx *Ginx) Run() error {
	if !ginx.init {
		panic("gin engine not init.")
	}
	if ginx.ListenPort == 0 {
		return errors.New("server listen port error")
	}

	// return ginx.Engine.Run(fmt.Sprintf(":%d", ginx.ListenPort))

	// 采用优雅退出的方式
	// 收到终止信号后强制退出秒数，15s
	DefaultHammerTime = 15 * time.Second
	return ListenAndServe(fmt.Sprintf(":%d", ginx.ListenPort), ginx.Engine)
}
