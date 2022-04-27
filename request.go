package ginCore

type Request struct {
}

type SearchRequest struct {
	Keyword string `json:"keyword"`
}

type PageRequest struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}
type IDRequest struct {
	ID int `json:"id" form:"id"`
}

type NameRequest struct {
	Name string `json:"name" form:"name"`
}

type CodeRequest struct {
	Code string `json:"code"`
}
type ModifyFieldRequest struct {
	ID    int         `json:"id" form:"id"`
	Field string      `json:"field" form:"field"`
	Value interface{} `json:"value" form:"value"`
}

type DatabaseRequest struct {
	DB    string `json:"db" form:"db"`
	Table string `json:"table" form:"table"`
	Field string `json:"field" form:"field"`
	Value string `json:"value" form:"value"`
}
