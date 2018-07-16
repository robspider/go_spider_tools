package datasource_config

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"strconv"
)

type DataSourceConfig struct {
	Host string
	Port int
	Db string
	//dbtype string
	Parameter_str string
	user string
	pwd string
	MaxOpenConns int
	MaxIdleConns int
}

func (config *DataSourceConfig) ConnectMysql() (*sql.DB,error) {
	db,err := sql.Open("mysql",config.user+":"+config.pwd+"@tcp("+config.Host+":"+strconv.Itoa(config.Port)+")/"+config.Db+"?"+config.Parameter_str)
	if(db != nil) {
		db.SetMaxOpenConns(config.MaxOpenConns)
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	return db,err
}

func NewDataSourceConfig(user,pwd string) *DataSourceConfig  {
	con := DataSourceConfig{
		user:user,
		pwd:pwd,
	}
	return &con
}