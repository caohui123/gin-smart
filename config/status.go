package config

// 状态组
type stateGroup struct {
	Key  string  `json:"key"`
	Desc string  `json:"desc"`
	List []state `json:"list"`
}

type state struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
	Desc  string `json:"desc"`
}

// 所有状态配置
func GetAllStates() []stateGroup {
	return allStates
}
