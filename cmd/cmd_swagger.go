package cmd

import (
	"github.com/pingdai/tools/codegen"
	"github.com/pingdai/tools/event/swagger"
	"github.com/spf13/cobra"
)

var cmdSwagger = &cobra.Command{
	Use:   "swagger",
	Short: "scan and generate swagger.json",
	Run: func(cmd *cobra.Command, args []string) {
		swaggerGenerator := swagger.SwaggerGenerator{}
		codegen.Generate(&swaggerGenerator)
	},
}

func init() {
	cmdRoot.AddCommand(cmdSwagger)
}
