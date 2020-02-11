
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/spf13/cobra"
    "github.com/kubeedge/beehive/pkg/core"
    "k8s.io/klog"
    
    _ "github.com/latelee/myproject/app4/cmd"
    "github.com/latelee/myproject/app4/cmd/server"
    "github.com/latelee/myproject/app4/cmd/update"
    "github.com/latelee/myproject/app4/cmd/gin1"
    "github.com/latelee/myproject/app4/cmd/udpp"
    "github.com/latelee/myproject/app4/cmd/tcpp"
    
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
        //fmt.Println("debug: ", debug, "deamon: ", deamon, "port:", port)
        // 执行业务程序，可用参数传递，或在内部读取配置文件
        registerModules()
        core.Run()

	},
    }

    initFlags(rootCmd)
    
    return rootCmd
}

func initFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
    cmd.PersistentFlags().BoolVarP(&deamon, "daemon", "d", false, "deamon mode")
    
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 89, "port number")
    cmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "http request timeout")	
}

func registerModules() {
    server.Register()
    update.Register()
    gin1.Register()
    udpp.Register()
    tcpp.Register()
}

func init() {
    //klog.InitFlags(nil)
}

func main() {
    klog.Info("hello klog...")
    //return
	cmd := NewCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
