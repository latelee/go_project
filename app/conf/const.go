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
    GroupName  = "config.test.io"
    APIVersion = "v1alpha1"
    Kind       = "myKind"
)

const (
    DataBaseDriverName = "sqlite3"
    DataBaseAliasName = "default"
    DataBaseDataSource = "/var/lib/my.db"
    
)
