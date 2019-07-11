package checkhealth

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func Init(routerGroup *gin.RouterGroup) {
	routerGroup.GET("checkhealth", CheckHealth)
}

func CheckHealth(c *gin.Context) {
	var body = struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Data struct{} `json:"data"`
	}{
		Code: 0,
		Msg:  "succ",
	}

	rsp, _ := json.Marshal(body)

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(rsp)
	return
}
