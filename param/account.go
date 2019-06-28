package param

// BindCloudAccountParam 绑定云商账号参数
type BindCloudAccountParam struct {
	CloudName   string `json:"cloudName" form:"cloudName" binding:"required"`
	AccountName string `json:"accountName" form:"accountName" binding:"required"`
	AccessKey   string `json:"accessKey" form:"accessKey" binding:"required"`
	SecurityKey string `json:"securityKey" form:"securityKey" binding:"required"`
}

// DoCloudAccountParam 操作云商账号参数
type DoCloudAccountParam struct {
	CloudName string `json:"cloudName" form:"cloudName" binding:"required"`
	AccountID string `json:"accountId" form:"accountId" binding:"required"`
}

// SearchCloudAccountParam 搜索云商账号参数
type SearchCloudAccountParam struct {
	CloudName   string `form:"cloudName"`
	AccountName string `form:"accountName"`
	AccountID   string `form:"accountId"`
}
