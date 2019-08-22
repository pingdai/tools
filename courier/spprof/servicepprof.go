package spprof

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pingdai/tools/constants"
	"os"
	"runtime/pprof"
	"time"
)

func Init(routerGroup *gin.RouterGroup) {
	routerGroup.GET("spprof", ServicePprof)
}

type ServicePprofReq struct {
	// pprof 检测时间，单位秒
	DelayTime time.Duration `json:"delay_time" form:"delay" binding:"required"`
}

func ServicePprof(c *gin.Context) {
	var body struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Data struct{} `json:"data"`
	}
	var err error
	req := &ServicePprofReq{}
	defer func() {
		if err != nil {
			body.Code = -1
			body.Msg = err.Error()
		} else {
			body.Code = 0
			body.Msg = "success"
		}
		rsp, _ := json.Marshal(body)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(rsp)
	}()

	if err = c.Bind(req); err != nil {
		return
	}

	cpuProfile, err := os.Create(fmt.Sprintf("%s_cpu_profile", os.Getenv(constants.EnvVarKeyProjectName)))
	if err != nil {
		return
	}
	defer cpuProfile.Close()
	err = pprof.StartCPUProfile(cpuProfile)
	if err != nil {
		return
	}
	defer pprof.StopCPUProfile()

	memProfile, err := os.Create(fmt.Sprintf("%s_mem_profile", os.Getenv(constants.EnvVarKeyProjectName)))
	if err != nil {
		return
	}
	defer memProfile.Close()
	if err = pprof.WriteHeapProfile(memProfile); err != nil {
		return
	}

	time.Sleep(req.DelayTime * time.Second)

	return
}
