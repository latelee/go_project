/*
TODO 短参数形式
*/

package main

import (
    "fmt"
    "os"
    "time"

    "github.com/spf13/cobra"    
)

var debug, deamon bool
var port int
var timeout time.Duration

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


func NewCommand() *cobra.Command {
    rootCmd := &cobra.Command{
	Use:   "demo",
	Short: "A demo app",
	Long: rootLong,
    Example: rootExample,
    Version: "1.0",
	Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("test cobra")
        fmt.Println("debug: ", debug, "deamon: ", deamon, "port:", port)
    
	},
    }
    
    // 使用默认的输出方式
    /*
    rootCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), "usage: %s\n\n", cmd.UseLine())
		return nil
	})
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "help: %s\n\n%s\n\n", cmd.Long, cmd.UseLine())
	})
    */
    
    InitFlags(rootCmd)
    
    return rootCmd
}

func InitFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
    cmd.PersistentFlags().BoolVar(&deamon, "d", false, "deamon mode")
    
	cmd.PersistentFlags().IntVar(&port, "p", 89, "port number")
    cmd.PersistentFlags().DurationVar(&timeout, "timeout", 10*time.Second, "http request timeout")
	
}

func main() {
	cmd := NewCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
