package requests

// AddTypeReq 新增字典类型请求
type AddTypeReq struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// DeleteTypeReq 批量删除type请求
type DeleteTypeReq struct {
	IdList []uint `json:"id_list"`
}

// UpdateTypeReq 修改字典类型请求
type UpdateTypeReq struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// GetTypeReq 获取字典类型请求
type GetTypeReq struct {
	Name          string `json:"name" form:"name"`
	Code          string `json:"code" form:"code"`
	Status        int    `json:"status" form:"status"`
	CreateAtBegin string `json:"create_at_begin" form:"create_at_begin"`
	CreateAtEnd   string `json:"create_at_end" form:"create_at_end"`
	Page          int    `json:"page" form:"page"`
	Limit         int    `json:"limit" form:"limit"`
}

// GetTypeRes 获取字典类型响应
type GetTypeRes struct {
	Id          uint   `json:"id"`
	CreatedAt   string `json:"created_at"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// AddItemReq 新增字典项请求
type AddItemReq struct {
	DictTypeCode string `json:"dict_type_code"`
	Label        string `json:"label"`
	Value        int    `json:"value"`
	Sort         int    `json:"sort"`
	Status       int    `json:"status"`
	Description  string `json:"description"`
	ExtendValue  string `json:"extend_value"`
}

// DeleteItemReq 批量删除item请求
type DeleteItemReq struct {
	IdList       []uint `json:"id_list"`
	DictTypeCode string `json:"dict_type_code"`
}

// UpdateItemReq 修改字典项请求
type UpdateItemReq struct {
	Id          uint   `json:"id"`
	Label       string `json:"label"`
	Value       int    `json:"value"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	Description string `json:"description"`
	ExtendValue string `json:"extend_value"`
}

// GetItemReq 获取字典项请求
type GetItemReq struct {
	CreateAtBegin string `json:"create_at_begin" form:"create_at_begin"`
	CreateAtEnd   string `json:"create_at_end" form:"create_at_end"`
	DictTypeCode  string `json:"dict_type_code" form:"dict_type_code"`
	Label         string `json:"label" form:"label"`
	Status        int    `json:"status" form:"status"`
	Page          int    `json:"page" form:"page"`
	Limit         int    `json:"limit" form:"limit"`
}

// GetItemRes 获取字典项响应
type GetItemRes struct {
	Id           uint   `json:"id"`
	CreateAt     string `json:"create_at"`
	DictTypeCode string `json:"dict_type_code"`
	Label        string `json:"label"`
	Value        int    `json:"value"`
	Sort         int    `json:"sort"`
	Status       int    `json:"status"`
	Description  string `json:"description"`
	ExtendValue  string `json:"extend_value"`
}
