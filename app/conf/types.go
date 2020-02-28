package conf


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