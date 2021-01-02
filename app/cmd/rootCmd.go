package cmd

import (
    "fmt"
	"os"
	// "io/ioutil"
	// "bytes"
    "path/filepath"
	//"golang.org/x/net/context"
	
	"k8s.io/klog"

	"github.com/spf13/cobra"
	// "github.com/spf13/afero"
	"github.com/spf13/viper"
	// "github.com/fsnotify/fsnotify"    
	"github.com/kubeedge/beehive/pkg/core"

    "webdemo/app/cmd/gin"
    "webdemo/app/cmd/tcpp"
	conf "webdemo/app/conf"
)


var (
	cfgFile string
	debug bool
	deamon bool
	mode string

	BuildTime string
	Version string
	shortDescription = `web demo`
    longDescription = `  A web demo
`
    example = `  comming soon...
`
)

func getVersion() string {
	return fmt.Sprintf("  version: %v build: %v\n", Version, BuildTime)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: shortDescription,
	Long: getVersion() + longDescription,
	Example: example,
	Version: getVersion(),
	Run: func(cmd *cobra.Command, args []string) {
		klog.Println("server run...")
		registerModules()
		core.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	// no child cmd...
    // rootCmd.AddCommand(test.NewCmdTest())
	return rootCmd.Execute()
}

// 命令行参数
func init() {	
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "specify the config file name")
	rootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", " ", "run mode: upgrade|normaltest")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
    rootCmd.PersistentFlags().BoolVarP(&deamon, "daemon", "d", false, "deamon mode")
    
	// cmd.PersistentFlags().IntVarP(&port, "port", "p", 89, "port number")
    // cmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "http request timeout")
}

// yaml 配置文件
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig();
	if  err != nil {
		klog.Println("'config.yaml' file read error:", err)

		os.Exit(0)
	} else {
		conf.ConfDBServer = viper.GetString("dbserver.datasource")

		conf.Gin.Enable = viper.GetBool("modules.gin.enable")
		conf.Gin.Port = viper.GetInt("modules.gin.port")

		conf.TcpServer.Enable = viper.GetBool("modules.tcpserver.enable")
		conf.TcpServer.Port = viper.GetInt("modules.tcpserver.port")
		
		conf.Vendors = viper.GetStringSlice("vendors")

		klog.Println(conf.Gin.Enable, conf.Gin.Port)
		klog.Println(conf.Vendors)
	}

	// //设置监听回调函数 似乎调用了2次
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	//klog.Printf("config is change :%s \n", e.String())
	// 	fullstate := viper.GetString("dbserver.timeout.fullstate")
	// 	klog.Println("fullstate--: ", fullstate)
	// 	//cancel()
	// })

	// viper.WatchConfig()

}

// 注册模块
func registerModules() {
    gin.Register()
    tcpp.Register()
}