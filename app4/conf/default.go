package conf

import (
    "fmt"
    "os"
    "path"
    "io/ioutil"
    //"github.com/latelee/go_project/pkg/com"
    "gopkg.in/yaml.v2"
    "k8s.io/klog"
)

// 解析文件，如不存在，使用默认值
// 如部分不存在，则使用默认值
func Config() *EdgeCoreConfig {
    cfg := newDefaultConfig()
	if err := cfg.parse(path.Join(DefaultConfigDir, DefaultCOnfigFile)); err != nil {
        klog.Print("config file not exist or parse error, using default one")
    }
    return cfg
}

func PrintConfig(config interface{}) {
    data, err := yaml.Marshal(config)
    if err != nil {
        fmt.Printf("Marshal min config to yaml error %v\n", err)
        return
    }
    fmt.Println("### config:")
    fmt.Printf("\n%v\n\n", string(data))
}

func PrintDefaultAndExit(config interface{}) {
    data, err := yaml.Marshal(config)
    if err != nil {
        fmt.Printf("Marshal min config to yaml error %v\n", err)
        os.Exit(1)
    }
    fmt.Println("### default config:")
    fmt.Printf("\n%v\n\n", string(data))
    os.Exit(0)
}

func (c *EdgeCoreConfig) parse(filename string) error {
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

func newDefaultConfig() *EdgeCoreConfig {
    return &EdgeCoreConfig{
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
            Edged: &Edged{
				Enable:         true,
            },
            Host: &Host{
				InterfaceName:  "eth0",
                IP:             "127.0.0.1",
            },
            Device: &Device{
				IP:             "127.0.0.1",
            },
        }, // end of Modules
    }
}