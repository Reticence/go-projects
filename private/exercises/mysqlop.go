// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-11-01 11:28
// @Version : 1.0
// @Software: Gogland

package exercises

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func RunMysqlOp() {
	db, _ = sql.Open("mysql", "root:123456@tcp(10.1.1.102:3306)/?charset=utf8")
	rows, _ := db.Query("SELECT user_age, user_sex FROM test.gouser")
	var userAge, userSex int32
	for rows.Next() {
		rows.Scan(&userAge, &userSex)
		fmt.Println(userAge, userSex)
	}
	//insert()
	query()
	fmt.Println(`string test`)
}

//插入demo
func insert() {
	stmt, err := db.Prepare("INSERT test.gouser (user_name,user_age,user_sex) values (?,?,?)")
	checkErr(err)
	res, err := stmt.Exec("tony", 20, 1)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println(id)
}

//查询demo
func query() {
	rows, err := db.Query("SELECT * FROM test.gouser")
	checkErr(err)

	//普通demo
	//for rows.Next() {
	//	var userId int
	//	var userName string
	//	var userAge int
	//	var userSex int

	//	rows.Columns()
	//	err = rows.Scan(&userId, &userName, &userAge, &userSex)
	//	checkErr(err)

	//	fmt.Println(userId)
	//	fmt.Println(userName)
	//	fmt.Println(userAge)
	//	fmt.Println(userSex)
	//}

	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns)) // 存放循环scan的指针
	values := make([]interface{}, len(columns))   // 存放scan后的结果值
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		fmt.Println(record)
	}
}

//更新数据
func update() {
	stmt, err := db.Prepare(`UPDATE test.gouser SET user_age=?,user_sex=? WHERE user_id=?`)
	checkErr(err)
	res, err := stmt.Exec(21, 2, 1)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(num)
}

//删除数据
func remove() {
	stmt, err := db.Prepare(`DELETE FROM test.gouser WHERE user_id=?`)
	checkErr(err)
	res, err := stmt.Exec(1)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(num)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
