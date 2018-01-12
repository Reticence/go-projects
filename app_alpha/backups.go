// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-12-25 17:17
// @Version : 0.1
// @Software: GoLand

package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

const (
	backupConf         = "BackupConf.json"
	sourcefilesmd5dict = "sourcefilesmd5dict"
)

var m map[string]string

type configurationFiles []configurationFile

type configurationFile struct {
	SourceDir string
	BackupDir string
}

type fileInfo struct {
	cf       configurationFile
	filePath string
	name     string
	size     int64
	isDir    bool
}

func (cfs *configurationFiles) readConfigurationFile() {
	var file *os.File
	var err error
	file, err = os.OpenFile(backupConf, os.O_RDONLY, 0644)
	if err != nil {
		file, err = os.OpenFile(pathJoin("D:/Temporary Space/go-io", backupConf), os.O_RDONLY, 0644)
		if err != nil {
			log.Fatalf("configuration file not found")
		}
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if err := json.Unmarshal(b, cfs); err != nil {
		log.Fatalf(err.Error())
	}
}

func (fi *fileInfo) backupFiles() {
	if fi.isDir || (fi.name == sourcefilesmd5dict && fi.filePath == "") {
		return
	}
	sourceFilePath := pathJoin(fi.cf.SourceDir, fi.filePath, fi.name)
	backupFilePath := pathJoin(fi.cf.BackupDir, fi.filePath, fi.name)
	if fi.size == getSize(backupFilePath) {
		sourceFileMD5 := getMD5(sourceFilePath)
		if value, ok := m[sourceFilePath]; ok && value == sourceFileMD5 && exists(backupFilePath) {
			return
		} else {
			m[sourceFilePath] = sourceFileMD5
		}
	}
	fmt.Print("Source '", sourceFilePath, "' => Backup '", backupFilePath, "' ")

	begin := time.Now().Unix()
	readFile, err := os.OpenFile(sourceFilePath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("open source file error")
		return
	}
	defer readFile.Close()

	os.MkdirAll(pathJoin(fi.cf.BackupDir, fi.filePath), os.ModeDir)
	writeFile, err := os.OpenFile(backupFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("open backup file error")
		return
	}
	defer writeFile.Close()

	if _, err = io.Copy(writeFile, readFile); err != nil {
		fmt.Println("copy error")
		return
	}
	fmt.Println("successful")
	fmt.Println("copy file:", sourceFilePath, time.Now().Unix()-begin)
}

func (fi fileInfo) deleteFiles() {
	sourceFilePath := pathJoin(fi.cf.SourceDir, fi.filePath, fi.name)
	backupFilePath := pathJoin(fi.cf.BackupDir, fi.filePath, fi.name)
	if !exists(sourceFilePath) && exists(backupFilePath) {
		fmt.Print("Remove '", backupFilePath, "' ")
		if fi.isDir {
			if err := os.RemoveAll(backupFilePath); err != nil {
				fmt.Println("delete dir error")
				return
			}
		} else {
			if err := os.Remove(backupFilePath); err != nil {
				fmt.Println("delete file error")
				return
			}
		}
		fmt.Println("successful")
	}
}

func (cf *configurationFile) getFiles(baseDir, filePath string) []fileInfo {
	fileInfos := make([]fileInfo, 0)
	dirList, e := ioutil.ReadDir(pathJoin(baseDir, filePath))
	if e != nil {
		fmt.Println("read dir error")
		return fileInfos
	}
	for _, v := range dirList {
		if v.Name() == sourcefilesmd5dict && filePath == "" {
			continue
		}
		fileInfos = append(fileInfos, fileInfo{cf: *cf, filePath: filePath, name: v.Name(), size: v.Size(), isDir: v.IsDir()})
		if v.IsDir() {
			fileInfos = append(fileInfos, cf.getFiles(baseDir, pathJoin(filePath, v.Name()))...)
		}
	}
	return fileInfos
}

func exists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

func pathJoin(args ...string) string {
	return path.Join(args...)
}

func getSize(filePath string) int64 {
	stat, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return stat.Size()
}

func loadDict(baseDir string) map[string]string {
	file, _ := os.OpenFile(pathJoin(baseDir, sourcefilesmd5dict), os.O_RDONLY, 0644)
	defer file.Close()
	b, _ := ioutil.ReadAll(file)
	json.Unmarshal(b, &m)
	return make(map[string]string, 1)
}

func saveDict(baseDir string) {
	file, _ := os.OpenFile(pathJoin(baseDir, sourcefilesmd5dict), os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	writer := bufio.NewWriter(file)
	b, _ := json.Marshal(m)
	writer.Write(b)
	writer.Flush()
}

func getMD5(filePath string) string {
	begin := time.Now().Unix()
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	hash := md5.New()
	if _, err = io.Copy(hash, reader); err != nil {
		return ""
	}
	fmt.Println("md5 calculate:", filePath, time.Now().Unix()-begin)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func main() {
	cfs := make(configurationFiles, 1)
	cfs.readConfigurationFile()
	for _, cf := range cfs {
		m = make(map[string]string)
		loadDict(cf.BackupDir)
		fileInfos := cf.getFiles(cf.SourceDir, "")
		for _, fi := range fileInfos {
			fi.backupFiles()
		}

		fileInfos = cf.getFiles(cf.BackupDir, "")
		for _, fi := range fileInfos {
			fi.deleteFiles()
		}
		saveDict(cf.BackupDir)
	}
}
