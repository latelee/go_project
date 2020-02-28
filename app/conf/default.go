package conf

import (
    "fmt"
    "os"
    "path"
    "io/ioutil"
    //"com"
    "gopkg.in/yaml.v2"
    "k8s.io/klog"
    
    //"github.com/spf13/pflag"
    "github.com/spf13/cobra"
)

var defconf bool
var showconf bool

// 本包内的标志，如打印配置文件
func AddFlag(cmd *cobra.Command) {
    cmd.PersistentFlags().BoolVar(&defconf, "defconfig", false, "print config information")
    cmd.PersistentFlags().BoolVar(&showconf, "showconfig", false, "show config information")
}

// 解析文件，如不存在，使用默认值
// 如部分不存在，则使用默认值
func Config() *AppCoreConfig {
    cfg := newDefaultConfig()
    if err := cfg.parse(path.Join(DefaultConfigDir, DefaultCOnfigFile)); err != nil {
        klog.Print("config file not exist or parse error, using default one")
    }
    return cfg
}

func PrintConfigAndExit(config interface{}) {
    if showconf == true {
        data, err := yaml.Marshal(config)
        if err != nil {
            fmt.Printf("Marshal min config to yaml error %v\n", err)
            return
        }
        fmt.Println("### config:")
        fmt.Printf("\n%v\n\n", string(data))
        os.Exit(0)
    }
}

func PrintDefaultAndExit() {
    if defconf == true {
        config := newDefaultConfig()
        data, err := yaml.Marshal(config)
        if err != nil {
            fmt.Printf("Marshal min config to yaml error %v\n", err)
            os.Exit(1)
        }
        fmt.Println("### default config:")
        fmt.Printf("\n%v\n\n", string(data))
        os.Exit(0)
    }
}

func (c *AppCoreConfig) parse(filename string) error {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        klog.Errorf("Failed to read configfile %s: %v", filename, err)
        return err
    }
    err = yaml.Unmarshal(data, c)
    if err != nil {
        klog.Errorf("Failed to unmarshal configfile %s: %v", filename, err)
        return err
    }
    return nil
}

func newDefaultConfig() *AppCoreConfig {
    return &AppCoreConfig{
        TypeMeta: TypeMeta{
            Kind:       Kind,
            APIVersion: path.Join(GroupName, APIVersion),
        },
        DataBase: &DataBase{
            DriverName: DataBaseDriverName,
            AliasName:  DataBaseAliasName,
            DataSource: DataBaseDataSource,
        },
        Modules: &Modules{
            Gin: &Gin{
                Enable:         true,
                Port:           4000,
            },
            UpdServer: &UpdServer{
                Enable:         true,
                Port:           10086,
            },
            TcpServer: &TcpServer{
                Enable:         true,
                Port:           8080,
            },
            DevServer: &DevServer{
                Enable:         true,
                Name:           "unknown",
                Protocol:       "unknown",
                Port:           9000,
            },
        }, // end of Modules
    }
}