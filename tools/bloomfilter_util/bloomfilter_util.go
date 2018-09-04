package bloomfilter_util

import (
	"github.com/willf/bloom"
	"database/sql"
	"time"
	"log"
	"strconv"
	"github.com/robspider/go_spider_tools/tools/gorpool"
)

func InitBloomFilter(filter *bloom.BloomFilter, db *sql.DB, pagesize int, table_name, uuid_colum_name string) {
	work_num := 48
    var p0 int64 = time.Now().UnixNano()
	log.Println("start init BloomFilter ")
    var bloom_filter_pool *gorpool.Pool
	var count int = 0
	row := db.QueryRow("select count(1) from " + table_name)
	row.Scan(&count)
	log.Printf("总数据量：%d",count)
	if(count == 0 ) {
		return
	}
	if(count < work_num) {
		bloom_filter_pool = gorpool.NewPool(count, pagesize*10).EnableWaitForAll(true).Start()
	}else{
		bloom_filter_pool = gorpool.NewPool(work_num, pagesize*10).EnableWaitForAll(true).Start()
	}
	var num int = count / pagesize
	for i := 0; i <= num; i++ {
		sql_str  := "select " + uuid_colum_name + " from " + table_name + " order by id asc limit " + strconv.Itoa(pagesize*i) + "," + strconv.Itoa(pagesize)
		bloom_filter_pool.AddJob(func() {
			updateBloomFilter(sql_str,db,bloom_filter_pool,filter,&count)
		})
	}
	time.Sleep(1*time.Second)
	bloom_filter_pool.WaitForAll()
	bloom_filter_pool.StopAll()
	var p1 int64 = time.Now().UnixNano()
	log.Println("布隆过滤器初始化完成，数据总量：" + strconv.Itoa(count) + ",耗时：" + strconv.FormatInt((p1-p0)/1000/1000, 10) + "ms")
}



func updateBloomFilter(sql_str string,db *sql.DB,bloom_filter_pool *gorpool.Pool,filter *bloom.BloomFilter,count *int){
	rows, err := db.Query(sql_str)
	if (err != nil) {
		panic(err)
	}
	for (rows.Next()) {
		uuid := ""
		err := rows.Scan(&uuid)
		if(err != nil){
			panic(err)
		}
		bloom_filter_pool.AddJob(func() {
			filter.AddString(uuid)
			*count = *count -1
			if(*count%10000 == 0){
				log.Printf("余量:%d",*count)
			}
		})
	}
	rows.Close()
}

