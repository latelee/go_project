package conf

const (
    DefaultConfigDir  = "./"
    DefaultCOnfigFile = "config.yaml"
)

const (
	GroupName  = "edgecore.config.kubeedge.io"
	APIVersion = "v1alpha1"
	Kind       = "EdgeCore"
)

const (
	// DataBaseDriverName is sqlite3
	DataBaseDriverName = "sqlite3"
	// DataBaseAliasName is default
	DataBaseAliasName = "default"
	// DataBaseDataSource is edge.db
	DataBaseDataSource = "/var/lib/kubeedge/edgecore.db"
)

///////////////////////////////////////

type EdgeCoreConfig struct {
	TypeMeta
	DataBase *DataBase `json:"database,omitempty"`
	Modules *Modules `json:"modules,omitempty"`
}

type TypeMeta struct {
	Kind       string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
    APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
}

// 数据库
type DataBase struct {
	DriverName string `json:"driverName,omitempty"`
	AliasName  string `json:"aliasName,omitempty"`
	DataSource string `json:"dataSource,omitempty"`
}

// 模块
type Modules struct {
	Edged *Edged `json:"edged,omitempty"`
    Host *Host `json:"host,omitempty"`
    Device *Device `json:"device,omitempty"`
}

type Edged struct {
	Enable bool `json:"enable,omitempty"`
}

type Host struct {
    InterfaceName string `json:"interfaceName,omitempty"`
	IP string `json:"ip,omitempty"`
}

type Device struct {
	IP string `json:"ip,omitempty"`
}