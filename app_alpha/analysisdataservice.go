// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-11-15 11:21
// @Version : 0.1
// @Software: GoLand

package main

import adu "go-projects/app_alpha/analysisdatautils"

func main() {
	adu.InitializationAll()
	//adu.syncMeridianIndexList()
	//adu.Ins2DB([]string{"(1,'a')"})
	//adu.Ins2DB([]string{"(1,'a')", "(2,'b')", "(3,'c')"})
	adu.Test()
}
