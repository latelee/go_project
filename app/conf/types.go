package conf


///////////////////////////////////////

type AppCoreConfig struct {
    TypeMeta
    DataBase *DataBase `json:"database,omitempty"`
    Modules  *Modules `json:"modules,omitempty"`
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
    Gin *Gin `json:"gin,omitempty"`
    UpdServer *UpdServer `json:"udpServer,omitempty"`
    TcpServer *TcpServer `json:"tcpServer,omitempty"`
    DevServer *DevServer `json:"devServer,omitempty"`
}

type Gin struct {
    Enable bool `json:"enable,omitempty"`
    Port   int `json:"port,omitempty"`
}

type UpdServer struct {
    Enable bool `json:"enable,omitempty"`
    Port   int `json:"port,omitempty"`
}

type TcpServer struct {
    Enable bool `json:"enable,omitempty"`
    Port   int `json:"port,omitempty"`
}

type DevServer struct {
    Enable   bool `json:"enable,omitempty"`
    Name     string `json:"name,omitempty"`
    Protocol string `json:"protocol,omitempty"`
	Port     int `json:"port,omitempty"`
}