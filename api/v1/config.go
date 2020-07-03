package v1

import "gitlab.gosccap.cn/bourse/avian/pkg/app"

// 全局配置
func Config(c *app.Context) app.Err {
	data := map[string]interface{}{
		"title":  "global config",
		"detail": "xxx",
	}
	c.Output(data)
	return nil
}
