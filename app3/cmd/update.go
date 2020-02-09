package cmd

import (
    "fmt"
	"github.com/spf13/cobra"
)

var file string

var theCmd = &cobra.Command{
	Use:   "update",
	Short: "update test",
	Long: `
    update test...`,
	Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("update test demo...")
	},
}

func init() {
	rootCmd.AddCommand(theCmd)
    
    // 本命令参数
    theCmd.Flags().StringVarP(&file, "file", "f", "foo.yaml", "config file")
}
