package app

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jangozw/gin-smart/pkg/util"
)

// 路由解析的生成文件和内容
type routeGenerateInfo struct {
	apiPackage              string
	apiFile                 string
	apiStructName           string
	paramFile               string
	paramRequestStructName  string
	paramResponseStructName string
	route                   Route
}

// 关闭自动生成文件
func routeManager(rou Route) error {
	return nil

	if !(IsEnvDev() || IsEnvLocal()) {
		return nil
	}
	// 组为根的不管
	if rou.Group == "/" {
		return nil
	}
	info, err := parseRouteGenerateInfo(rou)
	if err != nil {
		PrintConsole(err)
		return err
	}
	if yes, err := IsPathExists(info.apiFile); err != nil {
		PrintConsole(err)
		return err
	} else if yes {
		return nil
	}
	apiTpl, err := apiTplContent(info)
	if err != nil {
		PrintConsole(err)
		return err
	}
	if err := writeToFile(info.apiFile, apiTpl); err != nil {
		PrintConsole(err)
		return err
	} else {
		return writeToParamsFile(info)
	}
}

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

func parseRouteGenerateInfo(r Route) (info routeGenerateInfo, err error) {
	// group := strings.TrimPrefix(r.Group, "/")
	key := r.Group + r.Path
	key = strings.TrimPrefix(key, "/")
	key = strings.TrimSuffix(key, "/")

	pathFields := strings.Split(key, "/")
	if len(pathFields) < 2 {
		return info, nil
	}
	paramsFileName := Setting.Boot + "/param/" + strings.Join(pathFields[0:2], "_") + ".go"
	fileNameFields := pathFields[1:]
	var apiFile string
	num := len(pathFields)
	if num == 2 {
		apiFile = fmt.Sprintf("/api/%s/%s.go", pathFields[0], pathFields[1])
	} else if num > 2 {
		dir := fmt.Sprintf("/api/%s/%s/", pathFields[0], pathFields[1])
		apiFile = dir + strings.Join(fileNameFields, "_") + ".go"
	}
	for i, v := range fileNameFields {
		fileNameFields[i] = util.Capitalize(v)
	}
	apiStructName := strings.Join(fileNameFields, "")

	info = routeGenerateInfo{
		route:                   r,
		apiPackage:              pathFields[0],
		apiFile:                 Setting.Boot + apiFile,
		apiStructName:           apiStructName,
		paramFile:               paramsFileName,
		paramRequestStructName:  fmt.Sprintf("%s", apiStructName),
		paramResponseStructName: fmt.Sprintf("%sResponse", apiStructName),
	}
	return info, nil
}

func getModule() (string, error) {
	reg := regexp.MustCompile(`module\s+(.*?)\n`)
	file := Setting.Boot + "/go.mod"
	content, err := util.ReadAll(file)
	if err != nil {
		return "", err
	}
	if reg.MatchString(string(content)) {
		ma := reg.FindStringSubmatch(string(content))
		return ma[1], nil
	}
	return "", errors.New("unknown")
}

func writeToParamsFile(info routeGenerateInfo) error {
	fileAppends := make([]string, 0)
	allFiles, err := util.GetAllFiles(Setting.Boot + "/param")
	if err != nil {
		return nil
	}
	existsRequestName := 0
	existsResponseName := 0
	reg := regexp.MustCompile(`type\s+` + info.paramRequestStructName + `\s+`)
	reg2 := regexp.MustCompile(`type\s+` + info.paramResponseStructName + `\s+`)
	for _, oneFile := range allFiles {
		by, err := util.ReadAll(oneFile)
		if err != nil {
			return err
		}
		content := string(by)
		// fmt.Println("检查", info.paramResponseStructName, info.paramRequestStructName)
		if reg2.MatchString(content) {
			existsResponseName++
			// fmt.Println("发现", info.paramResponseStructName)
		}
		if reg.MatchString(content) {
			existsRequestName++
			// fmt.Println("发现", info.paramRequestStructName)
		}
	}
	if yes, err := util.IsPathExists(info.paramFile); err != nil {
		return err
	} else if !yes {
		PrintConsole(info.paramFile + " 不存在")
		fileAppends = append(fileAppends, "package param")
		fileAppends = append(fileAppends, "\n")
	}
	if existsRequestName == 0 {
		fileAppends = append(fileAppends, fmt.Sprintf("\n// %s %s%s", info.route.Method, info.route.Group, info.route.Path))
		fileAppends = append(fileAppends, fmt.Sprintf("type %s struct {}", info.paramRequestStructName))
	}
	if existsResponseName == 0 {
		fileAppends = append(fileAppends, fmt.Sprintf("type %s struct {}", info.paramResponseStructName))
	}
	if len(fileAppends) == 0 {
		return nil
	}
	appendContent := strings.Join(fileAppends, "\n")
	if err := writeToFile(info.paramFile, appendContent); err != nil {
		PrintConsole("append param file err :" + err.Error())
		return nil
	}
	// fmt.Println("+++", info.paramFile, ":")
	// fmt.Println(appendContent)
	return nil
}

func writeToFile(file string, content string) error {
	dir := util.Dir(file)
	if dir == "" {
		// PrintConsole()
		return errors.New(file + ":dir is empty!")
	}

	// 文件夹是否存在
	if ok, err := util.IsPathExists(dir); err != nil {
		PrintConsole(err)
		return err
	} else if !ok {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			PrintConsole("mkdir err:" + dir + err.Error())
			return err
		}
	}
	// 文件是否存在
	/*	if ok, err := utils.IsPathExists(file); err != nil {
			return err
		} else if ok {
			return nil
		}*/

	return util.AppendToFile(file, content)
}

func apiTplContent(info routeGenerateInfo) (string, error) {
	tpl := Setting.Boot + "/pkg/tpl/api.tpl"
	bytes, err := util.ReadAll(tpl)
	if err != nil {
		return "", err
	}
	module, err := getModule()
	if err != nil {
		return "", err
	}

	repMap := map[string]string{
		"{package}":     info.apiPackage,
		"{description}": info.route.Method + " " + info.route.FullPath(),
		"{module}":      module,
		"{structName}":  info.apiStructName,
		"{input}":       "param." + info.paramRequestStructName,
		"{output}":      "param." + info.paramResponseStructName,
	}
	tplContent := string(bytes)
	for k, v := range repMap {
		tplContent = strings.ReplaceAll(tplContent, k, v)
	}
	return tplContent, nil
}
