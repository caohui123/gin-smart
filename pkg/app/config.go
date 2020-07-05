package app

// config.ini 文件的struct

type Config struct {
	Env struct {
		Env string `json:"env"`
	} `json:"env"`

	Log struct {
		LogDir string `json:"log_dir"`
	} `json:"log"`

	General struct {
		DefaultPageSize    uint `json:"default_page_size"`
		MaxPageSize        uint `json:"max_page_size"`
		TokenExpireSeconds int `json:"token_expire_seconds"`
	} `json:"general"`

	Server struct {
		Listen int `json:"listen"`
	} `json:"server"`

	Redis struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		DbNum    int    `json:"db_num"`
	} `json:"redis"`

	Database struct {
		Schema   string `json:"schema"`
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"database"`

	Encrypt struct {
		JwtSecret string `json:"jwt_secret"`
	} `json:"encrypt"`
}

func (c *Config) Check() error {
	return nil
}
