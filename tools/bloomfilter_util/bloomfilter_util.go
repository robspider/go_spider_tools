package bloomfilter_util

import (
	"github.com/willf/bloom"
	"database/sql"
	"frank/gosipder/tools/gorpool"
	"time"
	"log"
	"strconv"
)

func InitBloomFilter(filter *bloom.BloomFilter, db *sql.DB, pagesize int, table_name, uuid_colum_name string) {
	work_num := 200
	p0 := time.Now().UnixNano()
	pool := gorpool.NewPool(work_num, 10000).Start().EnableWaitForAll(true)
	count := 0

	row := db.QueryRow("select count(1) from " + table_name)
	row.Scan(&count)
	num := count / pagesize
	for i := 0; i <= num; i++ {
		if (pool.GetPoolWorkCount() == work_num) {
			i--
			continue
		}
		sql_str := "select " + uuid_colum_name + " from " + table_name + " order by id asc limit " + strconv.Itoa(pagesize*i) + "," + strconv.Itoa(pagesize)
		pool.AddJob(func() {
			log.Println(sql_str)
			rows, err := db.Query(sql_str)
			if (err != nil) {
				log.Println(err)
			}
			for (rows.Next()) {
				uuid := ""
				rows.Scan(&uuid)
				pool.AddJob(func() {
					filter.AddString(uuid)
				})
			}
			rows.Close()
		})
	}
	pool.WaitForAll()
	pool.StopAll()
	p1 := time.Now().UnixNano()
	log.Println("布隆过滤器初始化完成，数据总量：" + strconv.Itoa(count) + ",耗时：" + strconv.FormatInt((p1-p0)/1000000, 10) + "ms")
	pool = nil
}
