package app

import (
	"fmt"
	"strconv"
)

const (
	TimeFormatFullDate = "2006-01-02 15:04:05" // 常规类型
	EnvLocal           = "local"
	EnvDev             = "dev"
	EnvTest            = "test"
	EnvProduction      = "production"
)

type StringNumber string

func (s *StringNumber) Int64() int64 {
	str := string(*s)
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func (s *StringNumber) Int() int {
	str := string(*s)
	i, _ := strconv.Atoi(str)
	return i
}

func (s *StringNumber) Uint() uint {
	return uint(s.Int())
}

type TriggerIF interface {
	Do()
}

// 触发
func Trigger(tg TriggerIF) {
	go tg.Do()
}

func PrintConsole(value interface{}) {
	var content string
	if err, ok := value.(error); ok {
		content = err.Error()
	}
	if content == "" {
		content = fmt.Sprintf("%v", value)
	}
	fmt.Println("---console---", content)
}
