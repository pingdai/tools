package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cmdRoot = &cobra.Command{
	Use:   "tools",
	Short: "ast tools",
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
