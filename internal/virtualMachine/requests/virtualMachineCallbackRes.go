package requests

type CreateVMCallback struct {
	Email  string `json:"email"`
	IpAddr string `json:"ip_addr"`
}

// DestroyResp 销毁虚拟机响应
type DestroyResp struct {
	Email string `json:"email"`
}

// NoticeResp 虚拟机通知
type NoticeResp struct {
	Email string `json:"email"`
	Data  string `json:"data"`
}
