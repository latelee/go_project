
package main

import (
    //"fmt"
    "os"
    "time"
    //"webdemo/pkg/com"
    //"flag"
    
    "github.com/spf13/cobra"
    "github.com/kubeedge/beehive/pkg/core"
    "k8s.io/klog"

    "webdemo/app/cmd/gin"
    "webdemo/app/cmd/tcpp"
    "webdemo/app/conf"
    
)

var debug, deamon bool
var port int
var timeout time.Duration
var mode string

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
        // 执行业务程序，可用参数传递，或在内部读取配置文件
        // 可以移到函数中
        opts := conf.Config()
        conf.PrintDefaultAndExit()
        conf.PrintConfigAndExit(opts)

        doMain(opts)
	},
    }

    initFlags(rootCmd)
    
    return rootCmd
}

func initFlags(cmd *cobra.Command) {
    cmd.PersistentFlags().StringVarP(&mode, "mode", "m", " ", "run mode: upgrade|normaltest")
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
    cmd.PersistentFlags().BoolVarP(&deamon, "daemon", "d", false, "deamon mode")
    
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 89, "port number")
    cmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "http request timeout")
    
    conf.AddFlag(cmd)
}

func registerModules(opts *conf.AppCoreConfig) {
    gin.Register(opts.Modules.Gin)
    tcpp.Register(opts.Modules.TcpServer)
}

// 似乎写不到文件
// 写文件必须将logtostderr设置为false
func init() {
    // klog.InitFlags(nil)
    // flag.Set("logtostderr", "false")
	// flag.Set("log_file", "myfile.log")
	// flag.Parse()
}

func doInit() {  
    //conf.PrintDefaultAndExit()
    //klog.Printf("opt: %V\n", opt)
}

func doMain(opts *conf.AppCoreConfig) {
    doInit()
    registerModules(opts)
    core.Run()
}

func main() {
    klog.Info("server run...")
	cmd := NewCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
