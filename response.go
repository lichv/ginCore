package ginCore

type Response struct {
	State   int         `json:"state" form:"state"`
	Data    interface{} `json:"data" form:"data"`
	Message string      `json:"message" form:"message"`
}

type PageResult struct {
	State int         `json:"state" form:"state"`
	List  interface{} `json:"list" form:"list"`
	Total int         `json:"total" form:"total"`
	Last  int         `json:"last" form:"last"`
}

type KeyResult struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
