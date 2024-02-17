package models

type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type FiltersInput struct {
	Skip    int64    `json:"skip"`
	Take    int64    `json:"take"`
	Filters []Filter `json:"input"`
}
