package util

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// 查找配置文件位置
// 查询运行参数,当前目录有没有，没有就就找到项目根目录以根目录定位
func FindConfigFile(filename string, flagKey string) (string, error) {
	// 1, 检查运行参数
	if find := FindArgValue(flagKey); find != "" {
		return find, nil
	}
	filename = strings.TrimPrefix(filename, "./")
	filename = strings.TrimPrefix(filename, "/")
	// 2, 当前路径只查找此文件名
	fs := strings.Split(filename, "/")
	last := fs[len(fs)-1]
	dir, err := os.Getwd()
	if err != nil {
		return ``, err
	}
	path := dir + `/` + filename
	if ok, err := IsPathExists(path); err != nil {
		return ``, err
	} else if ok {
		return path, nil
	}
	path = dir + `/` + last
	if ok, err := IsPathExists(path); err != nil {
		return ``, err
	} else if ok {
		return path, nil
	}
	// 3， 根木录定位文件
	root, err := ProjectRoot()
	if err != nil {
		return ``, err
	}
	return root + `/` + filename, nil
}

// 项目根目录，以go.mod 所在目录为根
func ProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return ``, err
	}
	dirs := strings.Split(dir, `/`)
	for i := len(dirs); i > 0; i-- {
		pDir := strings.Join(dirs[0:i-1], `/`)
		if ok, err := IsPathExists(pDir + `/go.mod`); err != nil {
			return ``, err
		} else if ok {
			return pDir, nil
		}
	}
	return ``, errors.New(`couldn't find project root'`)
}

func FindArgValue(name string) string {
	reg1 := regexp.MustCompile(`-{0,2}` + name + `=(.*)`)
	args := os.Args[1:len(os.Args)]
	for i, key := range args {
		if reg1.MatchString(key) {
			m := reg1.FindStringSubmatch(key)
			if len(m) >= 2 {
				return m[1]
			}
		} else if key == name || key == `-`+name || key == `--`+name {
			if len(args) >= i+2 {
				return args[i+1]
			}
		}
	}
	return ""
}

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
	if err != nil {
		return err
	}
	defer fileObj.Close()
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
