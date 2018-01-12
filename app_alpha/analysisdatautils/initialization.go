// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-11-20 13:57
// @Version : 0.1
// @Software: GoLand

package analysisdatautils

import (
	"database/sql"
)

func InitializationAll() {
	// 初始化数据库连接
	db, _ = sql.Open("mysql", "root:123456@tcp(10.1.1.102:3306)/analysis_data_checkup?charset=utf8")
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.Ping()

	// 同步meridian_index_list表
	syncMeridianIndexList()

	// 查询最大related_id & person_id
	r, p := queryMaxId()
	// 启动ID获取协程
	nextRelatedIdChan = make(chan int, 5)
	go func(ch chan int) {
		for {
			r += 1
			ch <- r
		}
	}(nextRelatedIdChan)
	nextPersonIdChan = make(chan int, 5)
	go func(ch chan int) {
		for {
			p += 1
			ch <- p
		}
	}(nextPersonIdChan)

	idcardMap = make(map[string]*IdentificationInfo)
	customerIdMap = make(map[string]*IdentificationInfo)
}
