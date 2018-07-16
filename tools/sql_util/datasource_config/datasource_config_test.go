package datasource_config

import (
	"testing"
	"log"
)

func TestDb(t *testing.T)  {
	con := NewDataSourceConfig("root","Lrw360+")
	con.MaxIdleConns = 20
	con.MaxOpenConns = 200
	con.Parameter_str = "charset=utf8"
	con.Host = "192.168.1.121"
	con.Port = 3306
	con.Db = "lrw360-map"

	db,_ := con.ConnectMysql()

	row := db.QueryRow("select count(1) from tb_company_test_hc_h5")
	count := 0

	row.Scan(&count)

	log.Println(count)
}
