// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-10-12 15:58
// @Version : 1.0
// @Software: Gogland

package exercises

import (
	"encoding/xml"
	"github.com/tealeg/xlsx"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type Result struct {
	Jbxx []Jbxx `xml:"jbxx"`
	Data Data   `xml:"data"`
	Zj   Zj     `xml:"zj"`
}

type Jbxx struct {
	Text string `xml:"text,attr"`
	Tjbh Tjbh   `xml:"tjbh"`
	Name Name   `xml:"name"`
	Sex  Sex    `xml:"sex"`
	Age  Age    `xml:"age"`
	Sfzh Sfzh   `xml:"sfzh"`
}

type Tjbh struct {
	Text string `xml:"text,attr"`
	Tjbh string `xml:",chardata"`
}

type Name struct {
	Text string `xml:"text,attr"`
	Name string `xml:",chardata"`
}

type Sex struct {
	Text string `xml:"text,attr"`
	Sex  string `xml:",chardata"`
}

type Age struct {
	Text string `xml:"text,attr"`
	Age  string `xml:",chardata"`
}

type Sfzh struct {
	Text string `xml:"text,attr"`
	Sfzh string `xml:",chardata"`
}

type Data struct {
	Text string `xml:"text,attr"`
	Ksmc []Ksmc `xml:"ksmc"`
}

type Ksmc struct {
	Text     string     `xml:"text,attr"`
	Xiangmu1 []Xiangmu1 `xml:"xiangmu1"`
	Xj       []Xj       `xml:"xj"`
}

type Xiangmu1 struct {
	Text     string     `xml:"text.attr"` // 项目名称1
	Xiangmu2 []Xiangmu2 `xml:"xiangmu2"`
}

type Xiangmu2 struct {
	Text string `xml:"text,attr"` // 项目名称2
	Jg   string `xml:"jg"`        // 结果
	Dw   string `xml:"dw"`        // 单位
	Ckfw string `xml:"ckfw"`      // 参考范围
	Ycts string `xml:"ycts"`      // 异常提示
	Ycbz string `xml:"ycbz"`      // 异常标志
}

type Xj struct {
	Text string `xml:"text,attr"`
	Xjqk string `xml:"xjqk"`
}

type Zj struct {
	Text string `xml:"text,attr"`
	Zs   string `xml:"zs"`
	Jy   string `xml:"jy"`
	Zjrq string `xml:"zjrq"`
}

func getDatetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func RunLoadxml() {
	dirName := "T:/test/aaa/"
	titleMap := make(map[string]string)

	fileList, _ := ioutil.ReadDir(dirName)
	fileListLen := len(fileList)
	percentTmp := 0
	for i, file := range fileList {
		filePath := dirName + file.Name()
		if strings.Contains(filePath, ".xml") {
			content, _ := ioutil.ReadFile(filePath)
			var result Result
			xml.Unmarshal(content, &result)

			for _, value := range result.Jbxx {
				fmt.Println(value)
			}

			for _, value := range result.Data.Ksmc {
				fmt.Println(value)
			}

			fmt.Println(result.Zj)
		}
		percent := 100 * i / fileListLen
		if percent%10 == 0 && percent > percentTmp {
			percentTmp = percent
			fmt.Println(getDatetime(), fmt.Sprintf("Progressing: %d%s", percent, "%"))
		}
	}

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Value = "haha"
	cell = row.AddCell()
	cell.Value = "xixi"

	err := file.Save("T:/file.xlsx")
	if err != nil {
		panic(err)
	}
	titleMap["a"] = "a"
	titleMap["b"] = "b"
	fmt.Println(getDatetime())
}
