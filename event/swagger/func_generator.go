package swagger

import (
	"fmt"
	"os/exec"
	"strings"
)

/*
* 当前只是通过系统命令调用第三方工具来生成的
* 后续改成扫描生成，不用再依赖第三方的工具
* https://github.com/go-swagger/go-swagger
* go get -u -v github.com/go-swagger/go-swagger
* 拉下来然后go install
 */

type SwaggerGenerator struct {
}

func (g *SwaggerGenerator) Load(cwd string) {

}

func (g *SwaggerGenerator) Process() {}

func (g *SwaggerGenerator) Output() {
	cmd := exec.Command("swagger", "generate", "spec", "-o", "./swagger.json")
	if err := cmd.Run(); err != nil {
		fmt.Printf("generate swagger doc failed! err:%v\n", err)
		if strings.Contains(err.Error(), "not found") {
			fmt.Println("Uage: " +
				"	go get -u -v github.com/go-swagger/go-swagger\n" +
				"	cd $GOPATH/src/github.com/go-swagger/go-swagger/cmd/swagger\n" +
				"	go install")
		}
		return
	}

	fmt.Println("generate swagger doc success!")
	return
}
