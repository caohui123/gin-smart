package lib

// import "C"
import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Unknwon/goconfig"
)

type Config struct {
	*goconfig.ConfigFile
}

// required config items
var confRequiredFields = map[string][]string{
	"log":      {"log_dir"},
	"server":   {"listen"},
	"database": {"user", "dbName"},
	"redis":    {"host"},
}

func NewConfig(file string) (*Config, error) {
	if ok, err := isPathExists(file); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New(fmt.Sprintf("config file %s is not exists, you can start with arg -config={path} ", file))
	}
	fmt.Println("init config, configPath is:", file)
	if c, err := goconfig.LoadConfigFile(file); err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't load config file %s, %s", file, err.Error()))
	} else {
		conf := &Config{c}
		if err := checkRequired(conf); err != nil {
			return nil, err
		}
		return conf, nil
	}
}

// check required fields in config file
func checkRequired(conf *Config) error {
	var msg []string
	for section, val := range confRequiredFields {
		for _, field := range val {
			if _, err := conf.GetValue(section, field); err != nil {
				msg = append(msg, fmt.Sprintf("Error: check config fail, %s in [%s] is required in config file", field, section))
			}
		}
	}
	logDir := conf.GetLogDir()
	if ok, err := isPathExists(logDir); err != nil {
		msg = append(msg, fmt.Sprintf("log path err: %s", err.Error()))
	} else if !ok {
		// msg = append(msg, fmt.Sprintf("log path not exists: %s", logPath))
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			return errors.New("make log dir err:" + err.Error())
		} else {
			fmt.Println("successfully make log dir:", logDir)
		}
	}
	str := ""
	if len(msg) > 0 {
		for _, v := range msg {
			str += v
		}
	}
	if str != "" {
		return errors.New(str)
	}
	return nil
}

//
func (c *Config) Section(section string) (map[string]string, error) {
	return c.GetSection(section)
}

//
func (c *Config) Get(section string, key string) (string, error) {
	return c.GetValue(section, key)
}

//
func (c *Config) GetInt(section string, key string) (int, error) {
	if v, err := c.Get(section, key); err != nil {
		return 0, err
	} else {
		return strconv.Atoi(v)
	}
}

// 过期时长
func (c *Config) GetTokenExpireSeconds() int64 {
	sec, _ := c.GetInt("encrypt", "token_expire_seconds")
	if sec >= 86400*365 {
		sec = 86400 * 365
	}
	return int64(sec)
}

func (c *Config) WebPort() (port int, err error) {
	return c.GetInt("server", "listen")
}

func (c *Config) GetLogDir() string {
	d, _ := c.Get("log", "log_dir")
	return d
}

// 判断文件夹是否存在
func isPathExists(dirPath string) (bool, error) {
	_, err := os.Stat(dirPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 签发token的签名秘钥
func (c *Config) GetJwtSecret() (s string, err error) {
	if s, err := c.Get("encrypt", "jwt_secret"); err != nil {
		return s, errors.New("couldn't get the config key : jwt_secret")
	} else {
		if len(s) == 0 {
			return s, errors.New("jwt_secret len too short")
		}
		return s, nil
	}
}

func (c *Config) JwtSecret() string {
	v, _ := c.Get("encrypt", "jwt_secret")
	return v
}

func (c *Config) Env() string {
	v, _ := c.Get("env", "env")
	return v
}
