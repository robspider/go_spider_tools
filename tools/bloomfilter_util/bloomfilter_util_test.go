package bloomfilter_util

import (
	"testing"
	"github.com/willf/bloom"
	"frank/gosipder/tools/sql_util"
	"frank/gosipder/tools"
)

func TestInitBloomFilter(t *testing.T) {
	filter := bloom.New(uint(10000000),32)
	con := sql_util.NewDataSourceConfig(tools.Mysql_user,tools.Mysql_pwd)
	con.MaxIdleConns = 20
	con.MaxOpenConns = 200
	con.Parameter_str = "charset=utf8"
	con.Host = tools.Mysql_host
	con.Port = tools.Mysql_port
	con.Db = "lrw360-map"
	db,_ := con.ConnectMysql()

	InitBloomFilter(filter,db,10000,"tb_company_text_hc_h5","uuid")
}
