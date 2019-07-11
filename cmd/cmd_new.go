package cmd

import (
	"fmt"
	"github.com/pingdai/tools/codegen"
	"github.com/pingdai/tools/event/newservice"
	"github.com/spf13/cobra"
)

var cmdNewFlagName string

var cmdNew = &cobra.Command{
	Use:   "new",
	Short: "new service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			panic(fmt.Errorf("need service name"))
		}

		clientGenerator := newservice.ServiceGenerator{
			ServiceName:  args[0],
			DatabaseName: cmdNewFlagName,
		}

		codegen.Generate(&clientGenerator)
	},
}

func init() {
	cmdRoot.AddCommand(cmdNew)

	cmdNew.Flags().
		StringVarP(&cmdNewFlagName, "db-name", "", "", "with db name")
}
