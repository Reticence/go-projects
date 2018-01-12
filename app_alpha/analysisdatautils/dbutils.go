// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-11-15 11:22
// @Version : 0.1
// @Software: GoLand

package analysisdatautils

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var db *sql.DB

func syncMeridianIndexList() {
	var db252 *sql.DB
	db252, _ = sql.Open("mysql", "applications:sU9pK6*upz%B=T2tJS-F@tcp(218.241.151.252:3306)/checkup_library?charset=utf8")
	defer db252.Close()
	var indexCode, indexName, parentCode, indexType, displaySign string
	//var indexType, displaySign int
	rows, err := db252.Query("SELECT index_code, index_name, parent_code, index_type, display_sign FROM meridian_index_list")
	if err != nil {
		panic(err)
	}

	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString("INSERT INTO meridian_index_list(index_code, index_name, parent_code, index_type, display_sign) VALUES\n")
	for rows.Next() {
		rows.Scan(&indexCode, &indexName, &parentCode, &indexType, &displaySign)
		sqlBuffer.WriteString("('")
		sqlBuffer.WriteString(indexCode)
		sqlBuffer.WriteString("','")
		sqlBuffer.WriteString(strings.Replace(indexName, "'", "\\'", -1))
		sqlBuffer.WriteString("','")
		sqlBuffer.WriteString(parentCode)
		sqlBuffer.WriteString("',")
		sqlBuffer.WriteString(indexType)
		sqlBuffer.WriteString(",")
		sqlBuffer.WriteString(displaySign)
		sqlBuffer.WriteString("),\n")
	}
	sqlTruncate := "TRUNCATE TABLE analysis_data_checkup.meridian_index_list;"
	sqlInsert := sqlBuffer.String()[:sqlBuffer.Len()-2]

	db.Exec(sqlTruncate)
	db.Exec(sqlInsert)
}

func queryMaxId() (r, p int) {
	row, err := db.Query("SELECT MAX(related_id), MAX(person_id) FROM basic_information")
	if err != nil {
		panic(err)
	}
	row.Next()
	row.Scan(&r, &p)
	return r, p
}

func Ins2DB(values []string) {
	sqlPrepare := "INSERT INTO preview_examination_results(related_id, index_code, index_source, ref_range_minimum, ref_range_maximum, ref_range_unit, raw_data, clean_data) VALUES "
	fmt.Println(sqlPrepare + strings.Join(values, ","))
}

func ExportData() {

}
