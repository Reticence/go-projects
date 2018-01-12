// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-10-20 10:34
// @Version : 1.0
// @Software: Gogland

package exercises

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

type Infos struct {
	Id   int
	Name string
}

func (i *Infos) setValue(iid int, iname string) Infos {
	i.Id = iid
	i.Name = iname
	return *i
}

func (i *Infos) setValue2(iid int, iname string) (id int, name string) {
	i.Id = iid
	i.Name = iname
	return i.Id, i.Name
}

func readGzFile(ipath string) {
	f, e := os.Open(ipath)
	if e != nil {
		panic(e)
	}
	defer f.Close()

	rg, _ := gzip.NewReader(f)
	rd := bufio.NewReader(rg)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		fmt.Print(line)
	}

}

func RunTestcodes() {
	var infos Infos
	fmt.Println(infos)
	infos.setValue(1, "tom")
	fmt.Println(infos)
	//mainHTTP()
	//readGzFile("T:/zzz2.sql.gz")
}
