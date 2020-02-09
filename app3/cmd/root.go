package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var debug, deamon bool
var port, timeout int

var (
	rootLong = `
    long description.
    A demo app using cobra command

`
	rootExample = `
    example message.
    demo -d -c -p 10086
`
)

var rootCmd = &cobra.Command{
	Use:   "demo",
	Short: "A demo app",
	Long: rootLong,
    Example: rootExample,
    Version: "1.0",
	Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("rootcmd...")
		
	},
}
func Execute() error {
    return rootCmd.Execute()
}

func init() {
    // 全局参数
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
    rootCmd.PersistentFlags().BoolVarP(&deamon, "daemon", "d", false, "deamon mode")
    
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 89, "port number")
    rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 10, "http request timeout")
}
