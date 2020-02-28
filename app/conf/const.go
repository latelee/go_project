package conf

const (
    ST_FAULT int = iota
    ST_IDLE
    ST_CHARGING
    ST_UNKNOWN
    ST_STOP
)

const (
    T_REALTIME int = iota
    T_INNERTEST
    T_INIT
    T_CONTROL
    T_CHARGEINFO

)

///////////////////////////////

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
