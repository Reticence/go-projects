// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-10-25 13:44
// @Version : 1.0
// @Software: Gogland

package main

import (
	"github.com/Reticence/go-projects/private/exercises"
)

func main() {
	//studys.RunParacalc()

	//exercises.RunLoadxml()
	//exercises.RunParselog()
	//exercises.RunTestcodes()
	//exercises.RunWinGui()
	exercises.MainHTTP("localhost:4000")
	//exercises.RunMysqlOp()
	//begin := time.Now().Second()
	//time.Sleep(1 * time.Second)
	//fmt.Println(time.Now().Second() - begin)
	//m := make(map[string]string)
	//val := m["abc"]
	//fmt.Println(val, val == "")
}