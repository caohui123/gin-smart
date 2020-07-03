package util

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// 判断文件夹是否存在
func IsPathExists(dirPath string) (bool, error) {
	_, err := os.Stat(dirPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 读取文件内容
func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

// 使用io.WriteString()函数进行数据的写入
func AppendToFile(filename, content string) error {
	content = fmt.Sprintf("%s\n", content)
	fileObj, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o777)
	defer fileObj.Close()
	if err != nil {
		return err
	}
	if _, err := io.WriteString(fileObj, content); err == nil {
		return err
	}
	return nil
}

func Dir(file string) string {
	data := strings.Split(file, "/")
	dir := strings.Join(data[0:len(data)-1], "/")
	return dir
}

// 正则匹配每一行
func MatchFileLine(fileName string, reg *regexp.Regexp) ([]string, error) {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	// 当前行所在的section
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if reg.MatchString(line) {
			return reg.FindStringSubmatch(line), nil
		}
	}
}

// 获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	// suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			/*ok := strings.HasSuffix(fi.Name(), ".go")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}*/

			files = append(files, dirPth+PthSep+fi.Name())
		}
	}

	// 读取子目录下文件
	/*	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}*/

	return files, nil
}
