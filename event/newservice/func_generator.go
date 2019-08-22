package newservice

import (
	"fmt"
	"github.com/pingdai/tools/codegen"
	"go/build"
	"path"
)

type ServiceGenerator struct {
	ServiceName  string
	DatabaseName string
	Root         string
	cwd          string
	JsonTag      JsonTag
}

type JsonTag struct {
	LogJsonTag   string
	GinxJsonTag  string
	GormxJsonTag string
}

func (g *ServiceGenerator) Load(cwd string) {
	g.cwd = cwd
	g.JsonTag.LogJsonTag = fmt.Sprintf("`json:\"log\"`")
	g.JsonTag.GinxJsonTag = fmt.Sprintf("`json:\"ginx\"`")
	g.JsonTag.GormxJsonTag = fmt.Sprintf("`json:\"gormx\"`")
}

func (g *ServiceGenerator) Process() {}

func (g *ServiceGenerator) Output() {
	outputs := codegen.Outputs{}

	// main.go
	codegen.NewGenFile("main", path.Join(g.ServiceName, "doc.go")).
		WithData(g).
		OutputTo(outputs)

	outputs.WriteFiles()

	pkg, _ := build.ImportDir(path.Join(g.cwd, g.ServiceName), build.ImportComment)
	g.Root = pkg.ImportPath

	// database
	if g.DatabaseName != "" {
		codegen.NewGenFile("database", path.Join(g.ServiceName, "database/doc.go")).
			WithData(g).
			OutputTo(outputs)

		outputs.WriteFiles()
	}

	// global.config.go
	codegen.NewGenFile("global", path.Join(g.ServiceName, "global/config.go")).
		WithData(g).
		Block(`
const (
	PROJECT_NAME = "{{ .Data.ServiceName }}"
)

var Config Cfg

func init() {
	{{ .PureUse "github.com/pingdai/tools/servicex" }}.SetServiceName(PROJECT_NAME)
	{{ .PureUse "github.com/pingdai/tools/servicex" }}.ConfP(&Config)
}

type Cfg struct {
	Log      		*{{ ( .PureUse "github.com/pingdai/tools/log" ) }}.Log		{{ .Data.JsonTag.LogJsonTag }}
	Ginx     		*{{ ( .PureUse "github.com/pingdai/tools/ginx" ) }}.Ginx	{{ .Data.JsonTag.GinxJsonTag }}
{{ if .Data.DatabaseName }}
	Gormx    *{{ .PureUse "github.com/pingdai/tools/gormx" }}.Gormx		{{ .Data.JsonTag.GormxJsonTag }}
{{ end }}
}
`,
		).OutputTo(outputs)

	outputs.WriteFiles()

	// routes
	codegen.NewGenFile("routes", path.Join(g.ServiceName, "routes/root.go")).
		WithData(g).
		Block(`
func RootRouter(engine *{{ .PureUse "github.com/gin-gonic/gin" }}.Engine) {
	root := engine.Group({{ .PureUse .Data.Root "global" }}.PROJECT_NAME)
	{
		{{ .PureUse "github.com/pingdai/tools/courier/swagger" }}.Init(root)
		{{ .PureUse "github.com/pingdai/tools/courier/checkhealth" }}.Init(root)
		{{ .PureUse "github.com/pingdai/tools/courier/spprof" }}.Init(root)
	}

	// Register your router
	// e.g.
	// apiRouter := root.Group("v1")
	// {
	//
	// }
	// TODO
}
`,
		).OutputTo(outputs)

	outputs.WriteFiles()

	// types
	codegen.NewGenFile("types", path.Join(g.ServiceName, "constants/types/doc.go")).WithData(g).Block(`
// Defined enum types here
	`).OutputTo(outputs)

	codegen.NewGenFile("errors", path.Join(g.ServiceName, "constants/errors/doc.go")).
		WithData(g).
		Block(`
// Defined error types here
		`).OutputTo(outputs)

	// modules
	codegen.NewGenFile("modules", path.Join(g.ServiceName, "modules/doc.go")).WithData(g).Block(`
// Defined sub modules here
	`).OutputTo(outputs)

	// main
	codegen.NewGenFile("main", path.Join(g.ServiceName, "main.go")).
		WithData(g).
		Block(`
	func main() {
		{{( .PureUse .Data.Root "routes" )}}.RootRouter({{ .PureUse .Data.Root "global" }}.Config.Ginx.Engine)
		{{( .PureUse .Data.Root "global" )}}.Config.Ginx.Run()
	}
	`,
		).OutputTo(outputs)

	outputs.WriteFiles()

	// config/dev.conf
	codegen.NewGenFile("config", path.Join(g.ServiceName, "config/dev.conf")).
		WithData(g).
		Block(
			`{
    "log": {
        "level": "DEBUG"
    },
    "ginx": {
		"gin_mode": "debug",	// 调试开发阶段设置为debug，上线需设置为release
        "listen_port": 8080
    }{{ if .Data.DatabaseName }},
    "gormx": {
        "max_open_conns": 20
    }{{ end }}
}
	`,
		).OutputTo(outputs)

	outputs.WriteFiles()

	return
}
