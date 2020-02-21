
package main

import (
    "fmt"
    "os"
    "time"
    "com"

    "github.com/spf13/cobra"
    "github.com/kubeedge/beehive/pkg/core"
    "k8s.io/klog"

    "github.com/latelee/myproject/app4/cmd/server"
    "github.com/latelee/myproject/app4/cmd/update"
    "github.com/latelee/myproject/app4/cmd/gin1"
    "github.com/latelee/myproject/app4/cmd/udpp"
    "github.com/latelee/myproject/app4/cmd/tcpp"
    
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
        fmt.Println("test cobra")
        //fmt.Println("debug: ", debug, "deamon: ", deamon, "port:", port)
        // 执行业务程序，可用参数传递，或在内部读取配置文件
        doMain()
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
}

func registerModules() {
    server.Register()
    //update.Register()
    gin1.Register()
    udpp.Register()
    tcpp.Register()
}

func init() {
    //klog.InitFlags(nil)
}

func doMain() {
    err := os.Chdir("/vagrant/golang/src/vendor/github.com/latelee/myproject/app4")
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
        registerModules()
        core.Run()
    }
}

func main() {
    klog.Info("hello klog ...")
    klog.Info("111111111111111111111111111111111111111")
    //return
	cmd := NewCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
