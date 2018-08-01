// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2018-05-11 16:54
// @Version : 0.1
// @Software: GoLand

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const sep = "；"

var d datas
var conf config

type config struct {
	HttpListenAddr       string
	LogFileName          string
	DiagnosisChannelSize int
}

type parameterJson struct {
	Symptoms []string
	Gender   string
	Age      int
}

type resultClassify struct {
	Classify1 []string
	Classify2 []string
	Classify3 []string
	Classify4 []string
}

type datas struct {
	diseaseMap        map[string]set      // 疾病: 症状集合
	diseaseGenderMap  map[string]set      // 性别: 疾病集合  keyset={'', '男性', '女性'}
	diseaseAgeMap     map[string]set      // 年龄: 疾病集合  keyset={'', '<1', '<3', '<=18', '<55'}
	symptomMap        map[string]set      // 症状: 疾病集合
	symptomGradeMap   map[string][]string // 症状: 分级列表  ([父级, 子级1, 子级2, ...])
	symptomGenderMap  map[string]set      // 性别: 症状集合  keyset={'', '男性', '女性'}
	symptomAgeMap     map[string]set      // 年龄: 症状集合  keyset={'', '婴幼儿或儿童'}
	symptomSynonymMap map[string]string   // 症状: 症状同义词
	diseaseList       []string
	diseaseCount      int
	symptomCount      int
}

type set struct {
	m map[string]string
	sync.RWMutex
}

func Set() *set {
	return &set{m: map[string]string{}}
}

func (s *set) add(item string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.m[item]; !ok {
		s.m[item] = ""
	}
}

func (s *set) addAll(items []string) {
	for _, item := range items {
		s.add(item)
	}
}

func (s *set) remove(item string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.m[item]; ok {
		delete(s.m, item)
	}
}

func (s *set) copy() *set {
	jstr, _ := json.Marshal(s.m)
	m := map[string]string{}
	json.Unmarshal(jstr, &m)
	return &set{m: m}
}

func (s *set) len() int {
	return len(s.m)
}

func (s *set) list() []string {
	s.RLock()
	defer s.RUnlock()
	list := make([]string, len(s.m))
	var i int
	for item := range s.m {
		list[i] = item
		i++
	}
	return list
}

func (s *set) contain(item string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *set) intersection(s2 set) {
	for item := range s.m {
		if !s2.contain(item) {
			s.remove(item)
		}
	}
}

func (s *set) union(s2 set) {
	for item := range s2.m {
		s.add(item)
	}
}

// 返回高层级的症状列表（包含输入症状）
func (d *datas) highLevelSymptoms(symptom string) (symptoms []string) {
	for symptom != "" {
		symptoms = append(symptoms, symptom)
		if tmp, ok := d.symptomGradeMap[symptom]; ok {
			symptom = tmp[0]
		}
	}
	for from, to := 0, len(symptoms)-1; from < to; from, to = from+1, to-1 {
		symptoms[from], symptoms[to] = symptoms[to], symptoms[from]
	}
	return
}

// 返回低层级的症状列表（包含输入症状）
func (d *datas) lowLevelSymptoms(symptom string) (symptoms []string) {
	symptoms = []string{symptom}
	symptomGrade := d.symptomGradeMap[symptom]
	if len(symptomGrade) > 1 {
		for _, symptom := range symptomGrade[1:] {
			symptoms = append(symptoms, d.lowLevelSymptoms(symptom)...)
		}
	}
	return
}

// 返回全部层级的症状列表（包含输入症状）
func (d *datas) allLevelSymptoms(symptom string) []string {
	high := d.highLevelSymptoms(symptom)
	low := d.lowLevelSymptoms(symptom)
	if len(high) > 1 {
		return append(high[:len(high)-1], low...)
	} else if len(low) > 1 {
		return append(high, low[1:]...)
	}
	return high
}

// 症状过滤器
func (d *datas) symptomFilter(symptom, gender string, age int) bool {
	ageKey := ""
	if age < 12 {
		ageKey = "婴幼儿或儿童"
	}
	symptomsGender := d.symptomGenderMap[gender]
	symptomsGender.union(d.symptomGenderMap[""])
	symptomsAge := d.symptomAgeMap[ageKey]
	symptomsAge.union(d.symptomAgeMap[""])
	return inList(symptom, symptomsGender.list()) && inList(symptom, symptomsAge.list())
}

// 疾病过滤器
func (d *datas) diseaseFilter(disease, gender string, age int) bool {
	ageKey := ""
	if age < 1 {
		ageKey = "<1"
	} else if age < 3 {
		ageKey = "<3"
	} else if age <= 18 {
		ageKey = "<=18"
	} else if age < 55 {
		ageKey = "<55"
	}
	diseases1 := d.diseaseGenderMap[gender]
	diseases2 := d.diseaseGenderMap[""]
	diseases3 := d.diseaseAgeMap[ageKey]
	diseases4 := d.diseaseAgeMap[""]
	return (diseases1.contain(disease) || diseases2.contain(disease)) && (diseases3.contain(disease) || diseases4.contain(disease))
}

// 获取伴随症状列表
func (d *datas) getConcomitantSymptoms(iSymptoms []string, gender string, age int) (int, int, []string) {
	if len(iSymptoms) == 0 {
		return 0, 0, []string{}
	}
	val := d.symptomMap[iSymptoms[0]]
	diseases := val.copy()
	for _, symptom := range iSymptoms[1:] {
		diseases.intersection(d.symptomMap[symptom])
	}
	if len(diseases.m) == 0 {
		return 0, 0, []string{}
	}
	cSymptomsMap := map[string]int{}
	for _, disease := range diseases.list() {
		if !d.diseaseFilter(disease, gender, age) {
			continue
		}
		symptoms := d.diseaseMap[disease]
		for _, symptom := range symptoms.list() {
			if inList(symptom, iSymptoms) {
				continue
			}
			for _, sSymptoms := range strings.Split(symptom, sep) {
				if _, ok := cSymptomsMap[sSymptoms]; ok {
					cSymptomsMap[sSymptoms] += 1
				} else {
					cSymptomsMap[sSymptoms] = 1
				}
			}
		}
	}
	var cSymptoms []string
	for symptom, count := range cSymptomsMap {
		if d.symptomFilter(symptom, gender, age) {
			cSymptoms = append(cSymptoms, fmt.Sprintf("%03d%s", count, symptom))
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(cSymptoms))) // 倒序排列
	rcSymptoms := make([]string, len(cSymptoms))
	for i, symptom := range cSymptoms {
		rcSymptoms[i] = symptom[3:]
	}
	return len(diseases.m), len(cSymptoms), rcSymptoms
}

// 根据症状相关性进行分组
func (d *datas) symptomsGroupProc(symptoms []string, gender string, age int) [][]string {
	f := func() [][]string {
		var fResult [][]string
		for i := len(symptoms); i > 0; i-- {
			for _, comb := range combinations(symptoms, i) {
				if dLength, cLength, _ := d.getConcomitantSymptoms(comb, gender, age); dLength > 0 || cLength > 0 {
					fResult = append(fResult, comb)
				}
			}
			if len(fResult) > 0 {
				return fResult
			}
		}
		return fResult
	}
	if len(symptoms) <= 1 {
		return [][]string{symptoms}
	}
	var gSymptoms [][]string
	for len(symptoms) > 0 {
		result := f()
		gSymptoms = append(gSymptoms, result...)
		symptomsSet := Set()
		for _, res := range result {
			symptomsSet.addAll(res)
		}
		for i := 0; i < len(symptoms); i++ {
			if symptomsSet.contain(symptoms[i]) {
				symptoms = append(symptoms[:i], symptoms[i+1:]...)
				i--
			}
		}
	}
	return gSymptoms
}

// 诊断疾病
func (d *datas) diagnosis(symptoms []string, gender string, age int) string {
	// 症状处理函数
	notInSymptoms := func(niSymptoms, refSymptoms []string) []string {
		var rel []string
		for _, symptom := range niSymptoms {
			if !inList(symptom, refSymptoms) {
				rel = append(rel, symptom)
			}
		}
		return rel
	}
	// 疾病诊断函数
	fDiagnosis := func(pSymptoms []string, resultChan chan [][]string) {
		gDiseases := Set()
		for _, symptom := range pSymptoms {
			gDiseases.union(d.symptomMap[symptom])
		}
		diseases := make([][]string, 4)
		for _, gDisease := range gDiseases.list() {
			if !d.diseaseFilter(gDisease, gender, age) {
				continue
			}
			gSymptoms := d.diseaseMap[gDisease]
			gSymptomsList := gSymptoms.list()
			// 输入中不包含在疾病中的症状
			rel1 := notInSymptoms(pSymptoms, gSymptomsList)
			// 疾病中不包含在输入中的症状
			rel2 := notInSymptoms(gSymptomsList, pSymptoms)
			rel1Len := len(rel1)
			rel2Len := len(rel2)
			if rel1Len == 0 {
				if rel2Len == 0 {
					diseases[0] = append(diseases[0], fmt.Sprintf("000%s", gDisease))
				} else {
					diseases[1] = append(diseases[1], fmt.Sprintf("%03d%s(%d)", rel2Len, gDisease, rel2Len))
				}
			} else {
				if rel2Len == 0 {
					diseases[2] = append(diseases[2], fmt.Sprintf("%03d%s", rel1Len, gDisease))
				}
			}
		}
		resultChan <- diseases
	}

	resultsParser := func(input set) []string {
		resultList := input.list()
		sort.Strings(resultList)
		output := make([]string, len(resultList))
		for i, resultL := range resultList {
			output[i] = string([]rune(resultL)[3:])
		}
		return output
	}

	dResult := make(map[string]resultClassify)

	fDiagnosisResultsChan := make(chan [][]string, conf.DiagnosisChannelSize)
	symptomsGroups := d.symptomsGroupProc(symptoms, gender, age)
	for _, symptomsGroup := range symptomsGroups {
		symptomsGroupGrade := make([][]string, len(symptomsGroup))
		for i, symptomG := range symptomsGroup {
			for _, symptom := range d.allLevelSymptoms(symptomG) {
				if _, ok := d.symptomMap[symptom]; ok {
					symptomsGroupGrade[i] = append(symptomsGroupGrade[i], symptom)
				}
			}
		}

		goroutineCount := 0
		for _, pSymptoms := range product(symptomsGroupGrade) {
			go fDiagnosis(pSymptoms, fDiagnosisResultsChan)
			goroutineCount++
		}

		results := []set{*Set(), *Set(), *Set(), *Set()}
		for i := 0; i < goroutineCount; i++ {
			for i, result := range <-fDiagnosisResultsChan {
				results[i].addAll(result)
			}
		}
		dResultKey := strings.Join(symptomsGroup, " ")
		dResultVal := resultClassify{
			Classify1: resultsParser(results[0]),
			Classify2: resultsParser(results[1]),
			Classify3: resultsParser(results[2]),
			Classify4: resultsParser(results[3]),
		}
		dResult[dResultKey] = dResultVal
	}
	close(fDiagnosisResultsChan)
	b, _ := json.Marshal(dResult)
	return string(b)
}

func (d *datas) readDatas(dir string) {
	fReadSet := func(name string) map[string]set {
		file, _ := os.OpenFile(name, os.O_RDONLY, 0644)
		defer file.Close()
		b, _ := ioutil.ReadAll(file)
		m := make(map[string][]string)
		json.Unmarshal(b, &m)
		mss := make(map[string]set)
		for k, v := range m {
			ds := Set()
			ds.addAll(v)
			mss[k] = *ds
		}
		return mss
	}
	fReadSlice := func(name string) map[string][]string {
		file, _ := os.OpenFile(name, os.O_RDONLY, 0644)
		defer file.Close()
		b, _ := ioutil.ReadAll(file)
		m := make(map[string][]string)
		json.Unmarshal(b, &m)
		return m
	}
	fReadString := func(name string) map[string]string {
		file, _ := os.OpenFile(name, os.O_RDONLY, 0644)
		defer file.Close()
		b, _ := ioutil.ReadAll(file)
		m := make(map[string]string)
		json.Unmarshal(b, &m)
		return m
	}
	d.diseaseMap = fReadSet(path.Join(dir, "jsonfile/disease_dict.json"))
	d.diseaseGenderMap = fReadSet(path.Join(dir, "jsonfile/disease_gender_dict.json"))
	d.diseaseAgeMap = fReadSet(path.Join(dir, "jsonfile/disease_age_dict.json"))
	d.symptomMap = fReadSet(path.Join(dir, "jsonfile/symptom_dict.json"))
	d.symptomGradeMap = fReadSlice(path.Join(dir, "jsonfile/symptom_grade_dict.json"))
	d.symptomGenderMap = fReadSet(path.Join(dir, "jsonfile/symptom_gender_dict.json"))
	d.symptomAgeMap = fReadSet(path.Join(dir, "jsonfile/symptom_age_dict.json"))
	d.symptomSynonymMap = fReadString(path.Join(dir, "jsonfile/symptom_synonym_dict.json"))
	d.diseaseList = nil
	d.diseaseCount = 0
	d.symptomCount = 0
}

// 返回输入元素是否存在于输入列表中
func inList(ele string, lst []string) bool {
	lstSet := Set()
	for _, l := range lst {
		lstSet.addAll(strings.Split(l, sep))
	}
	for _, eleSP := range strings.Split(ele, sep) {
		if lstSet.contain(eleSP) {
			return true
		}
	}
	return false
}

func combineLoop(arr []string, r []string, i int, n int, output chan<- []string) {
	if n <= 0 {
		return
	}
	rlen := len(r) - n
	alen := len(arr)
	for j := i; j < alen; j++ {
		r[rlen] = arr[j]
		if n == 1 {
			or := make([]string, len(r))
			copy(or, r)
			output <- or
		} else {
			combineLoop(arr, r, j+1, n-1, output)
		}
	}
}

//对数组进行组合
func combinations(arr []string, n int) (combs [][]string) {
	output := make(chan []string)
	r := make([]string, n)
	go func() {
		combineLoop(arr, r, 0, n, output)
		close(output)
	}()
	for comb := range output {
		combs = append(combs, comb)
	}
	return combs
}

func productLoop(strs []string, slice [][]string, output chan<- []string) {
	sLen := len(slice)
	for _, str := range slice[0] {
		if sLen == 1 {
			output <- append(strs, str)
		} else {
			productLoop(append(strs, str), slice[1:], output)
		}
	}
}

func product(slice [][]string) (prods [][]string) {
	output := make(chan []string)
	var strs []string
	go func() {
		productLoop(strs, slice, output)
		close(output)
	}()
	for prod := range output {
		prods = append(prods, prod)
	}
	return prods
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析url传递的参数，对于POST则解析响应包的主体（request body）
	if r.Method == "POST" {
		rBody, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		//结构已知，解析到结构体
		var pj parameterJson
		json.Unmarshal([]byte(rBody), &pj)
		begin := time.Now().UnixNano()
		diagnosisResult := d.diagnosis(pj.Symptoms, pj.Gender, pj.Age)
		log.Printf("Parameter: %s  Diagnosis(%.3fms): %s", rBody, float64(time.Now().UnixNano()-begin)/1000000, diagnosisResult)
		fmt.Fprint(w, string(diagnosisResult))
	} else {
		fmt.Fprint(w, "错误的请求! 请使用POST方法。数据格式: {\"Symptoms\":[\"咳嗽\",\"发热\",\"乏力\"],\"Gender\":\"男性\",\"Age\":18}")
	}
}

func main() {
	// 获取执行文件的绝对路径
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// 读取配置文件，不能存在时使用默认值
	confFilePath := path.Join(dir, "diagnosis_config.json")
	confFile, err := os.OpenFile(confFilePath, os.O_RDONLY, 0644)
	defer confFile.Close()
	if err != nil {
		conf = config{HttpListenAddr: "localhost:8008", LogFileName: "diagnosis_service.log", DiagnosisChannelSize: 8}
		confFileDefault, _ := os.OpenFile(confFilePath, os.O_CREATE, 0644)
		defer confFileDefault.Close()
		b, _ := json.Marshal(conf)
		writer := bufio.NewWriter(confFileDefault)
		writer.Write(b)
		writer.Flush()
	} else {
		b, _ := ioutil.ReadAll(confFile)
		json.Unmarshal(b, &conf)
	}
	// 打开log文件
	logFile, _ := os.OpenFile(path.Join(dir, conf.LogFileName), os.O_APPEND|os.O_CREATE, 0644)
	defer logFile.Close()
	// 设置log输出到文件和控制台
	im := io.MultiWriter([]io.Writer{logFile, os.Stdout}...)
	log.SetOutput(im)
	// 读取字典数据
	d.readDatas(dir)
	log.Println(fmt.Sprintf("Server started at '%s'", conf.HttpListenAddr))
	// 设置HTTP地址的解析函数; 开启服务
	http.HandleFunc("/diagnosis", handler)
	http.ListenAndServe(conf.HttpListenAddr, nil)
}
