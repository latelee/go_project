
package main

import (
    //"fmt"
    "os"
    "time"
    "com"
    //"flag"
    
    "github.com/spf13/cobra"
    "github.com/kubeedge/beehive/pkg/core"
    "k8s.io/klog"

    "github.com/latelee/myproject/app/cmd/devServer"
    "github.com/latelee/myproject/app/cmd/gin"
    "github.com/latelee/myproject/app/cmd/udpp"
    "github.com/latelee/myproject/app/cmd/tcpp"
    "github.com/latelee/myproject/app/conf"

    "github.com/latelee/myproject/app/pkg/update"
    
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
    udpp.Register()
    tcpp.Register()
    devServer.Register(opts.Modules.DevServer)
}

// 似乎写不到文件
// 写文件必须将logtostderr设置为false
func init() {
    klog.InitFlags(nil)
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

    err := os.Chdir("/vagrant/golang/src/vendor/github.com/latelee/myproject/app")
    if err != nil {
        klog.Printf("cant change dir.\n")
        //return
    }
    if mode == "upgrade" { // 升级切换功能
        klog.Printf("upgrade mode run and test\n")
        update.ProcessUpgrade()
    } else if mode == "normaltest"{ // 这是模拟升级后的，仅测试
        com.Sleep(10000)
        klog.Printf("normal test exit\n")
        os.Exit(43); // 如果不指定，默认返回0， 这里只是测试
    } else { // 业务程序
        registerModules(opts)
        core.Run()
    }
}

func main() {
    //return
	cmd := NewCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
