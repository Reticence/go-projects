// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-11-01 16:04
// @Version : 1.0
// @Software: Gogland

package main

import (
	"bufio"
	"database/sql"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB

type fileInfoSlice []os.FileInfo

func (f fileInfoSlice) Len() int {
	return len(f)
}
func (f fileInfoSlice) Swap(i, j int) { // 重写 Swap() 方法
	f[i], f[j] = f[j], f[i]
}
func (f fileInfoSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return f[j].Size() < f[i].Size()
}

type codeContent struct {
	MeridianIndex map[string]string
	CodeContent   map[string]interface{}
}

type infos struct {
	MeridianIndexByteLen uint64
	CodeContentByteLen   uint64
}

func (cc *codeContent) contains(key string) bool {
	value := cc.CodeContent[key]
	return value != nil
}

func (cc *codeContent) serialize(filePath string) {
	file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	writer := bufio.NewWriter(file)

	infos := &infos{}
	bm, e1 := json.Marshal(cc.MeridianIndex)
	if e1 != nil {
		panic(e1)
	}

	bc, e2 := json.Marshal(cc.CodeContent)
	if e2 != nil {
		panic(e2)
	}
	infos.MeridianIndexByteLen = uint64(len(bm))
	infos.CodeContentByteLen = uint64(len(bc))
	binary.Write(writer, binary.LittleEndian, infos)
	writer.Write(bm)
	writer.Write(bc)
	writer.Flush()
}

func (cc *codeContent) deserialize(filePath string) bool {
	file, _ := os.OpenFile(filePath, os.O_RDONLY, 0644)
	defer file.Close()
	reader := bufio.NewReader(file)

	infos := &infos{}
	err := binary.Read(reader, binary.LittleEndian, infos)
	if err != nil {
		fmt.Println(err)
		return false
	}

	jsonBytes := make([]byte, 0, infos.MeridianIndexByteLen+infos.CodeContentByteLen)
	b := make([]byte, 1024*4)
	for {
		n, err := reader.Read(b)
		jsonBytes = append(jsonBytes, b[:n]...)
		if err == io.EOF {
			break
		}
	}
	if err := json.Unmarshal(jsonBytes[:infos.MeridianIndexByteLen], &cc.MeridianIndex); err != nil {
		fmt.Println(err)
		return false
	}

	if err := json.Unmarshal(jsonBytes[infos.MeridianIndexByteLen:], &cc.CodeContent); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (cc *codeContent) loadCodeContent(rebuild bool, savePath string) {
	db, _ = sql.Open("mysql", "departmenthbd:A$KMb!hi9x+j!PGxmqsL@tcp(218.241.151.253)/?charset=utf8")
	codeContentPath := savePath + "code_content.bin"
	if rebuild {
		begin := time.Now().Unix()
		tx, _ := db.Begin()
		tx.Exec("LOCK TABLE test.code_content WRITE, test.code_classify WRITE, unstructureddb.sys_unstructured_data A READ, unstructureddb.sys_standard_description B READ, unstructureddb.strd_code_corr C READ, checkup_library.meridian_index_list D READ")
		tx.Exec("CALL unstructureddb.`rebuild_library`()")
		tx.Exec("UNLOCK TABLES")
		err := tx.Commit()
		if err != nil {
			tx.Rollback()
		}
		fmt.Println(time.Now().Unix() - begin)
	}

	begin := time.Now().Unix()
	if rebuild || !cc.deserialize(codeContentPath) {
		os.Remove(codeContentPath)

		var indexCode, indexName, codeContent string
		rows, _ := db.Query("SELECT index_code, index_name FROM checkup_library.meridian_index_list")
		for rows.Next() {
			rows.Scan(&indexCode, &indexName)
			cc.MeridianIndex[indexCode] = indexName
		}

		rows, _ = db.Query("SELECT CONCAT(`index_code`, TRIM(`code_content`)) AS `code_content` FROM test.code_content")
		for rows.Next() {
			rows.Scan(&codeContent)
			cc.CodeContent[codeContent] = 1
		}

		cc.serialize(codeContentPath)
	}
	fmt.Println(time.Now().Unix() - begin)
}

func (cc *codeContent) fileProcessing(taskinfo, path string) {
	dirpath := path[:len(path)-5] + "/"
	os.RemoveAll(dirpath)

	indexdict := make(map[int]string)
	results := make(map[string]map[string]interface{})
	titles := make(map[string]string)
	file, _ := xlsx.OpenFile(path)
	sheet := file.Sheets[0]
	for i, row := range sheet.Rows {
		for j, cell := range row.Cells {
			if i == 0 {
				position := strings.Index(cell.Value, "-")
				indexcode := cell.Value[:position]
				position = strings.Index(cell.Value, "@@") + 2
				titles[indexcode] = indexcode + "_" + cell.Value[position:]
				if cc.MeridianIndex[indexcode] != "" {
					if indexcode[:5] != "TJ.1." && indexcode[:6] != "TJ.3.1" {
						indexdict[j] = indexcode
					}
				} else {
					errorfmt := `指标 "` + indexcode + `"(colnum = ` + string(j) + `) 不存在!`
					os.Mkdir(dirpath, os.ModeDir)
					ioutil.WriteFile(dirpath+"Error.fmt", []byte(errorfmt), 0644)
					return
				}
				continue
			}
			if indexdict[j] != "" && cell.Value != "" {
				indexcode := indexdict[j]
				if cc.CodeContent[indexcode+strings.ToUpper(cell.Value)] == nil {
					if results[indexcode] == nil {
						results[indexcode] = make(map[string]interface{})
					}
					results[indexcode][cell.Value] = 1
				}
			}
		}
		if i%1000 == 0 && i > 0 {
			fmt.Printf("taskid/total = %s  rownum = %d  path = %s\n", taskinfo, i, path)
		}
	}

	if len(results) > 0 {
		os.Mkdir(dirpath, os.ModeDir)
		data := make([][]string, 0, len(results))
		data = append(data, []string{"指标名称", "待处理数量"})
		file := xlsx.NewFile()
		sheet, _ := file.AddSheet("Sheet1")
		column := 0
		for indexcode, uncodes := range results {
			rownum := 0
			cell := sheet.Cell(0, column)
			cell.Value = titles[indexcode]
			data = append(data, []string{indexcode, strconv.Itoa(len(uncodes))})
			for uncode := range uncodes {
				rownum++
				cell := sheet.Cell(rownum, column)
				cell.Value = uncode
			}
			column++
		}
		file.Save(dirpath + "unCode.xlsx")
		f, _ := os.Create(dirpath + "Detail.csv") //创建文件
		defer f.Close()
		f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
		w := csv.NewWriter(f)         //创建一个新的写入文件流
		w.WriteAll(data)
		w.Flush()
	}
	fmt.Printf("taskid/total = %s  rownum = %d  path = %s  %s", taskinfo, sheet.MaxRow, path, "Finished!")
}

func sortBySize(fInfos []os.FileInfo) []os.FileInfo {
	fileInfos := make([]os.FileInfo, 0, len(fInfos))
	for _, file := range fInfos {
		if !file.IsDir() {
			fileInfos = append(fileInfos, file)
		}
	}
	sort.Sort(fileInfoSlice(fileInfos))
	return fileInfos
}

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Have request..")
	//fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	//file, _ := os.OpenFile("T:/test.html", os.O_RDONLY, 0644)
	//fb, _ := ioutil.ReadAll(file)
	//w.Write(fb)
	t, _ := template.ParseFiles("T:/test.html")
	t.Execute(w, nil)
}

func main() {
	//baseDir := "T:/test/"
	//files, _ := ioutil.ReadDir(baseDir)
	//files = sortBySize(files)
	//cc := &codeContent{make(map[string]string), make(map[string]interface{})}
	//cc.loadCodeContent(false, "T:/")
	//for _, fi := range files {
	//	fmt.Println(fi.Size(), fi.Name())
	//	cc.fileProcessing("1/1", baseDir+fi.Name())
	//}
	addr := ":8008"
	http.HandleFunc("/", handler)
	http.ListenAndServe(addr, nil)
}
