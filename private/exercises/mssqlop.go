// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-11-20 12:55
// @Version : 0.1
// @Software: GoLand

package exercises

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"net/url"
)

var mssql *sql.DB

func RunMssqlOp() {
	//query := url.Values{}
	//query.Add("connection timeout", fmt.Sprintf("%d", 10))

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword("sa", "123456"),
		Host:   fmt.Sprintf("%s:%d", "10.1.1.102", 3306),
		// Path:  instance, // if connecting to an instance instead of a port
		//RawQuery: query.Encode(),
	}

	connectionString := u.String()
	fmt.Println(connectionString)
	mssql, _ = sql.Open("mssql", connectionString)
	mssql.Close()
}
