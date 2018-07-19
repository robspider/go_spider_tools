package bloomfilter_util

import (
	"testing"
	"github.com/willf/bloom"
	"github.com/robspider/go_spider_tools/tools/sql_util/datasource_config"
	"github.com/robspider/go_spider_tools/tools"
	"github.com/satori/go.uuid"
)

func TestInitBloomFilter(t *testing.T) {
	filter := bloom.New(uint(10000000),32)
	con := datasource_config.NewDataSourceConfig(tools.Mysql_user,tools.Mysql_pwd)
	con.MaxIdleConns = 20
	con.MaxOpenConns = 200
	con.Parameter_str = "charset=utf8"
	con.Host = tools.Mysql_host
	con.Port = tools.Mysql_port
	con.Db = "lrw360-map"
	db,_ := con.ConnectMysql()

	InitBloomFilter(filter,db,40000,"tb_company_text_hc_h5","uuid")
}

func Test2(t *testing.T)  {
	filter := bloom.New(uint(10000000),32)
	con := datasource_config.NewDataSourceConfig(tools.Mysql_user,tools.Mysql_pwd)
	con.MaxIdleConns = 20
	con.MaxOpenConns = 400
	con.Parameter_str = "charset=utf8"
	con.Host = tools.Mysql_host
	con.Port = tools.Mysql_port
	con.Db = "lrw360-map"
	db,_ := con.ConnectMysql()
	InitBloomFilter(filter,db,40000,"tb_spider_company_describe_hc","describe_id")
}

func BenchmarkInitBloomFilter(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		filter := bloom.New(uint(10000000),32)
		for pb.Next(){
			u := uuid.NewV4()
			filter.Add(u.Bytes())
		}
	})
}

func BenchmarkUuid(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			uuid.NewV4().Bytes()
		}
	})
}