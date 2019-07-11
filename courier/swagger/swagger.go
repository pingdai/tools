package swagger

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pingdai/tools/constants"
	"io/ioutil"
	"os"
)

func getSwaggerJSON() []byte {
	data, err := ioutil.ReadFile("./swagger.json")
	if err != nil {
		return data
	}
	return data
}

func Init(routerGroup *gin.RouterGroup) {
	routerGroup.GET("", Swagger)
	routerGroup.GET("doc", SwaggerDoc)
}

func SwaggerDoc(c *gin.Context) {
	projectName := os.Getenv(constants.EnvVarKeyProjectName)
	html := &bytes.Buffer{}
	html.WriteString(`<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="//cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.17.4/swagger-ui.css" >
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>
    <script src="//cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.17.4/swagger-ui-bundle.js"> </script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.17.4/swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      // Build a system
      var ui = SwaggerUIBundle({
        url: "` + fmt.Sprintf("http://%s/%s", c.Request.Host, projectName) + `",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })

      window.ui = ui
    }
  </script>
  </body>
</html>`)

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Write(html.Bytes())
	return
}

func Swagger(c *gin.Context) {
	json := &bytes.Buffer{}
	json.Write(getSwaggerJSON())

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.Write(json.Bytes())
	return
}
