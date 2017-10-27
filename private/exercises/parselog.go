// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-10-24 13:13
// @Version : 1.0
// @Software: Gogland

package exercises

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"flag"
	"io"
	"os"
	"strings"
)

type parse struct {
	s        string
	res      []string
	monthMap map[string]string
}

func getMonthMap() map[string]string {
	return map[string]string{
		"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04", "May": "05", "Jun": "06",
		"Jul": "07", "Aug": "08", "Sep": "09", "Oct": "10", "Now": "11", "Dec": "12",
	}
}

func readGzLog(path string, rch chan string) {
	f, e := os.Open(path)
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
		rch <- line
	}
	close(rch)
}

func write2CSV(path string, wch chan []string) {
	f, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if e != nil {
		panic(e)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)
	for s := range wch {
		w.Write(s)
		w.Flush()
	}
}

func (p *parse) parseProcessing() {
	rs := []rune(p.s)
	date := strings.Split(string(rs[strings.Index(p.s, "[")+1:strings.Index(p.s, ":")]), "/")
	others := strings.Split(p.s, "\"")
	p.res[0] = string(rs[:strings.Index(p.s, " -")])
	p.res[1] = date[2] + p.monthMap[date[1]] + date[0]
	p.res[2] = others[1]
	p.res[3] = others[5]
	p.res[4] = others[7]
	p.res[5] = others[9]
	p.res[6] = others[11]
	p.res[7] = others[13]
}

func (p *parse) run(rch chan string, wch chan []string) {
	wch <- []string{"ip", "时间", "url", "访问代理", "AppVer", "ApiVer", "Apiplatform", "userinfo"}
	for line := range rch {
		p.s = line
		p.parseProcessing()
		wch <- p.res
	}
	close(wch)
}

func RunParselog() {
	inputfile := flag.String("inputfile", "", "Input file path.")
	outputfile := flag.String("outputfile", "./output.csv", "Outputfile file path.")
	flag.Parse()

	rch := make(chan string, 10)
	wch := make(chan []string, 10)
	p := parse{"", make([]string, 8, 8), getMonthMap()}

	go readGzLog(*inputfile, rch)
	go p.run(rch, wch)
	write2CSV(*outputfile, wch)
}
