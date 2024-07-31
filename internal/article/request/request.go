package request

type SearchRequest struct {
	Query    string `json:"query" form:"query"`
	Page     int    `json:"page" form:"page"`
	Limit    int    `json:"limit" form:"limit"`
	Category int    `json:"category" form:"category"`
	Kind     int    `json:"kind" form:"kind"`
}
