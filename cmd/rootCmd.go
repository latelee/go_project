package cmd

import (
	"fmt"
	"os"
	"runtime"

	// "io/ioutil"
	// "bytes"
	"path/filepath"
	//"golang.org/x/net/context"

	"webdemo/pkg/com"
	"webdemo/pkg/klog"

	"github.com/spf13/cobra"
	// "github.com/spf13/afero"
	"github.com/spf13/viper"
	// "github.com/fsnotify/fsnotify"
	"github.com/kubeedge/beehive/pkg/core"

	"webdemo/cmd/gin"
	"webdemo/cmd/tcpp"
	conf "webdemo/common/conf"
)

var (
	cfgFile string
	debug   bool
	deamon  bool
	mode    string

	BuildTime        string
	Version          string
	shortDescription = `web demo`
	longDescription  = `  A web demo
`
	example = `  comming soon...
`
)

func getVersion() string {
	return fmt.Sprintf("web  version: %v build: %v\n", Version, BuildTime)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     filepath.Base(os.Args[0]),
	Short:   shortDescription,
	Long:    getVersion() + longDescription,
	Example: example,
	Version: getVersion(),
	Run: func(cmd *cobra.Command, args []string) {
		conf.RunningOS = runtime.GOOS
		conf.RunningARCH = runtime.GOARCH

		conf.Args = args
		fmt.Printf("server start running on Platform: %v %v %v\n", conf.RunningOS, conf.RunningARCH, getVersion())
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
	// TOCHECK 配置文件和命令行指定，谁优先级高？似乎是配置文件
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "specify the config file name")
	rootCmd.PersistentFlags().StringVarP(&conf.RunMode, "mode", "m", "", "mode: client|website")
	// rootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", " ", "run mode: upgrade|normaltest")
	// rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
	// rootCmd.PersistentFlags().BoolVarP(&deamon, "daemon", "d", false, "deamon mode")

	rootCmd.PersistentFlags().IntVarP(&conf.Gin.Port, "port", "p", 9000, "port number")
	// cmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "http request timeout")
}

// yaml 配置文件
func initConfig() {
	conf.Config = viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		conf.Config.SetConfigFile(cfgFile)
	} else {
		conf.Config.AddConfigPath("./")
		conf.Config.SetConfigName("config")
		conf.Config.SetConfigType("yaml")
	}

	conf.Config.AutomaticEnv() // read in environment variables that match

	err := conf.Config.ReadInConfig()
	if err != nil {
		klog.Println("'config.yaml' file read error:", err)

		os.Exit(0)
	} else {
		conf.ConfDBServer = conf.Config.GetString("dbserver.datasource")

		conf.Gin.Enable = conf.Config.GetBool("modules.gin.enable")
		conf.Gin.Port = conf.Config.GetInt("modules.gin.port")

		conf.TcpServer.Enable = conf.Config.GetBool("modules.tcpserver.enable")
		conf.TcpServer.Port = conf.Config.GetInt("modules.tcpserver.port")

		conf.Vendors = conf.Config.GetStringSlice("vendors")

		conf.AppVersion = getVersion()
		tmpstr := conf.Config.GetString("setting.data_file")
		if len(tmpstr) != 0 {
			conf.DataFileDir = tmpstr
		}

		//  https
		conf.HttpsEnable = conf.Config.GetBool("https.enable")
		conf.HttpsCertFile = conf.Config.GetString("https.cert_file")
		conf.HttpsKeyFile = conf.Config.GetString("https.key_file")

		// klog.Println(conf.Gin.Enable, conf.Gin.Port)
		// klog.Println(conf.Vendors)
	}

	viper.BindPFlags(rootCmd.PersistentFlags())

	curDir := com.GetRunningDirectory()

	klog.Println("test config flag:", curDir, conf.Gin.Enable, conf.Gin.Port)

	// //设置监听回调函数 似乎调用了2次
	// conf.Config.OnConfigChange(func(e fsnotify.Event) {
	// 	//klog.Printf("config is change :%s \n", e.String())
	// 	fullstate := conf.Config.GetString("dbserver.timeout.fullstate")
	// 	klog.Println("fullstate--: ", fullstate)
	// 	//cancel()
	// })

	// conf.Config.WatchConfig()

}

// 注册模块
func registerModules() {
	gin.Register()
	tcpp.Register()
}
