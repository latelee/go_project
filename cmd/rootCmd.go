package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

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

	web "webdemo/cmd/web"
	website "webdemo/cmd/website"
	conf "webdemo/common/conf"
)

var (
	cfgFile          string
	port             int
	BuildTime        string
	Version          string
	shortDescription = `web demo`
	longDescription  = ` A web demo
`
	example = `  comming soon...
`
)

func getVersion() string {
	return fmt.Sprintf("web demo version: %v build: %v\n", Version, BuildTime)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     filepath.Base(os.Args[0]),
	Short:   shortDescription,
	Long:    getVersion() + longDescription,
	Example: example,
	Version: getVersion(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// TODO 找一个好方法添加子命令
func Execute() error {
	// rootCmd.AddCommand(test.NewCmd())
	rootCmd.AddCommand(web.NewCmd())
	rootCmd.AddCommand(website.NewCmd())
	// rootCmd.AddCommand(tcpp.NewCmd())

	return rootCmd.Execute()
}

// 命令行参数
func init() {
	// OnInitialize 会在真正执行命令前调用，但有的配置提前读取，因此单独调用
	// 如此一来，viper监听功能应该无作用
	// initConfig()
	cobra.OnInitialize(initConfig)

	// BoolVarP 支持长短命令，默认为false，输入--print或-p即为true
	// rootCmd.PersistentFlags().StringVarP(&conf.RunMode, "mode", "m", "", "mode: client|website")
	// rootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", " ", "run mode: upgrade|normaltest")
	// rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
	// rootCmd.PersistentFlags().BoolVarP(&deamon, "daemon", "d", false, "deamon mode")

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", conf.CfgFile, "specify the config file name")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", conf.DEFUALT_PORT, "port number")
}

// yaml 配置文件
func initConfig() {
	conf.Config = viper.New()

	if cfgFile == "" {
		cfgFile = conf.CfgFile
	}

	// 如果不是默认值才设置
	if cfgFile != conf.CfgFile {
		conf.CfgFile = cfgFile
	}

	fmt.Println("read file: ", conf.CfgFile)
	conf.Config.SetConfigFile(conf.CfgFile)

	conf.Config.AutomaticEnv() // read in environment variables that match

	// 如果失败，尝试data目录下的
	err := conf.Config.ReadInConfig()
	if err != nil {
		fmt.Println("read error ", err.Error())
		conf.CfgFile = filepath.Join("data", conf.CfgFile)
		conf.Config.SetConfigFile(conf.CfgFile)
		err = conf.Config.ReadInConfig()
		if err != nil {
			fmt.Printf("[%v] file read error: %v\n", conf.CfgFile, err.Error())
			// conf.TcpServer.Port = 8080
			os.Exit(0)
		}
	}

	// 说明：这里先找日志配置，以便能用日志函数
	// 但是即使不初始化，也能使用，只是输出到终端中
	log_dir := "log"
	log_prefix := ""
	log_level := 3
	log_size := 10 * 1024 * 1024
	log_type := 3
	log_time := 5

	tmpstr := conf.Config.GetString("setting.log.log_dir")
	if tmpstr != "" {
		log_dir = tmpstr
	}
	tmpstr = conf.Config.GetString("setting.log.log_prefix")
	if tmpstr != "" {
		log_prefix = tmpstr
	}
	tmpstr = conf.Config.GetString("setting.log.log_level")
	if tmpstr != "" {
		log_level = com.Str2Int(tmpstr)
	}
	tmpstr = conf.Config.GetString("setting.log.log_size")
	if tmpstr != "" {
		log_size = com.Str2Int(tmpstr)
	}
	tmpstr = conf.Config.GetString("setting.log.log_type")
	if tmpstr != "" {
		log_type = com.Str2Int(tmpstr)
	}

	needLogNode := false
	conf.HostName, _ = os.Hostname()

	if conf.Config.GetInt("setting.log.log_node") == 1 {
		needLogNode = true
	}

	if needLogNode == true {
		log_dir = filepath.Join(log_dir, fmt.Sprintf("%v_%v", log_prefix, conf.HostName))
	}

	tmpstr = conf.Config.GetString("setting.log.log_time")
	if tmpstr != "" {
		log_time = com.Str2Int(tmpstr)
	}

	conf.LogDir = log_dir
	klog.Init_normal(log_dir, log_prefix, log_level, log_size, log_type, log_time)

	conf.Gin.Enable = conf.Config.GetBool("modules.gin.enable")

	conf.Gin.Port = conf.Config.GetInt("modules.gin.port")
	// 不为默认值，设置
	if port != conf.DEFUALT_PORT {
		conf.Gin.Port = port
	}

	conf.ConfDBServer = conf.Config.GetString("dbserver.datasource")

	conf.Gin.Enable = conf.Config.GetBool("modules.gin.enable")
	conf.Gin.Port = conf.Config.GetInt("modules.gin.port")

	conf.TcpServer.Enable = conf.Config.GetBool("modules.tcpserver.enable")
	conf.TcpServer.Port = conf.Config.GetInt("modules.tcpserver.port")

	conf.Vendors = conf.Config.GetStringSlice("vendors")

	conf.AppVersion = getVersion()

	tmpstr = conf.Config.GetString("setting.data_file")
	if len(tmpstr) != 0 {
		conf.DataFileDir = tmpstr
	}

	//  https
	conf.HttpsEnable = conf.Config.GetBool("https.enable")
	conf.HttpsCertFile = conf.Config.GetString("https.cert_file")
	conf.HttpsKeyFile = conf.Config.GetString("https.key_file")

	viper.BindPFlags(rootCmd.PersistentFlags())

	conf.CurDir = com.GetRunningDirectory()

	// klog.Println("test config flag:", conf.CurDir, conf.Gin.Enable, conf.Gin.Port)
	showBanner()

	// //设置监听回调函数 似乎调用了2次
	// conf.Config.OnConfigChange(func(e fsnotify.Event) {
	// 	//klog.Printf("config is change :%s \n", e.String())
	// 	fullstate := conf.Config.GetString("dbserver.timeout.fullstate")
	// 	klog.Println("fullstate--: ", fullstate)
	// 	//cancel()
	// })

	// conf.Config.WatchConfig()

}

func showBanner() {
	conf.RunningOS = runtime.GOOS
	conf.RunningArch = runtime.GOARCH

	conf.AppVersion = getVersion()
	conf.StartTime = time.Now()

	klog.Printf("start running on Platform: %v(%v %v) %v\n", conf.RunningOS, conf.RunningArch, conf.HostName, getVersion())

}
