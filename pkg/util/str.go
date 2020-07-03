package util

import (
	"crypto/md5"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Location struct {
	Country      string
	Province     string
	City         string
	Area         string
	ProvinceCode int
	CityCode     int
	AreaCode     int
}

func (l *Location) MakeLocation() string {
	return fmt.Sprintf("%s,%s,%s,%s", l.Country, l.Province, l.City, l.Area)
}

// 解析 location 里有省市区地址： "中国,北京市,北京市,东城区" "中国,浙江省,杭州市,萧山区"
func ParseLocation(addr string) (Location, error) {
	addr = strings.TrimSpace(addr)
	loc := Location{}
	matches := regexp.MustCompile(`中国,(\S+),(\S+),(\S+)`).FindStringSubmatch(addr)

	if len(matches) != 4 {
		return loc, errors.New(addr + ": 腾讯地图定位解析省市区失败!")
	}
	loc.Province = matches[1]
	loc.City = matches[2]
	loc.Area = matches[3]
	return loc, nil
}

func SpliteToString(a []int64, sep string) string {
	if len(a) == 0 {
		return ""
	}
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(b, sep)
}

// 将字符串加密成 md5
func String2md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has) // 将[]byte转成16进制
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				// fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
