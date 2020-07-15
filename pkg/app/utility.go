package app

import (
	"fmt"

	"github.com/jangozw/gin-smart/pkg/util"
)

const (
	EnvLocal      = "local"
	EnvDev        = "dev"
	EnvTest       = "test"
	EnvProduction = "production"
	CtxStartTime  = "ctx-start-time"
)

type Pager struct {
	util.Pager
}

type paramsCheck interface {
	Check() error
}

type TriggerIF interface {
	Do()
}

// 触发
func Trigger(tg TriggerIF) {
	go tg.Do()
}

// api 监听地址
func HttpServeAddr() string {
	if Cfg != nil {
		return fmt.Sprintf(":%d", Cfg.General.ApiPort)
	}
	return ""
}

// 环境
func IsEnvLocal() bool {
	return CurrentEnv() == EnvLocal
}

func IsEnvDev() bool {
	return CurrentEnv() == EnvDev
}

func IsEnvTest() bool {
	return CurrentEnv() == EnvTest
}

func IsEnvProduction() bool {
	return CurrentEnv() == EnvProduction
}

func CurrentEnv() string {
	if Cfg != nil {
		return Cfg.General.Env
	}
	return ""
}
