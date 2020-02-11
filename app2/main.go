/*

BoolVar  长参数形式
BoolVarP 长短参数形式

PersistentFlags() 和 Flags() 在简单应用场合并无差别
*/

package main

import (
    "fmt"
    "os"
    "os/signal"
	"syscall"

    "time"

    "github.com/spf13/cobra"
    
    app "github.com/latelee/myproject/app2/app"
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
        app.Demo(debug)
        GracefulShutdown()
	},
    }
    
    // 使用默认的输出方式，不用自定义的形式
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
    cmd.PersistentFlags().BoolVarP(&deamon, "daemon", "d", false, "deamon mode")
    
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 89, "port number")
    cmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "http request timeout")
	
}

// 最后调用
func GracefulShutdown() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	select {
	case s := <-c:
		fmt.Printf("Get os signal %v\n", s.String())
		// 此处处理销毁事务。。。
	}
}

func main() {
	cmd := NewCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
