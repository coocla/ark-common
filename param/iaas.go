package param

// SearchRegionParam 搜索地域参数
type SearchRegionParam struct {
	CloudName string `form:"cloudName"`
}

// RunInstanceParam 创建并运行一台vm的参数
type RunInstanceParam struct {
	AccountID       string `json:"accountId" form:"accountId" binding:"required"`
	RegionID        string `json:"regionId" form:"regionId" binding:"required"`
	ZoneID          string `json:"zoneId" form:"zoneId" binding:"required"`
	ImageID         string `json:"imageId" form:"imageId" binding:"required"`
	InstanceType    string `json:"instanceType" form:"instanceType" binding:"required"`
	HostName        string `json:"hostName" form:"hostName" binding:"required"`
	InstanceName    string `json:"instanceName" form:"instanceName" binding:"required"`
	KeyPairID       string `json:"keyPairId" form:"keyPairId" binding:"required"`
	SecurityGroupID string `json:"securityGroupId" form:"securityGroupId"`
	SubnetID        string `json:"subnetId" form:"subnetId"`
	VPCID           string `json:"vpcId" form:"vpcId"`
	DiskType        string `json:"diskType" form:"diskType"`
	DiskSize        int    `json:"diskSize" form:"diskSize"`
	Numbers         int    `json:"numbers,default=1" form:"numbers,default=1"`
}

// SearchInstanceParam 搜索实例参数
type SearchInstanceParam struct {
	RegionID  string `form:"regionId"`
	CloudName string `form:"cloudName"`
}

// SearchInstanceImageParam 搜索镜像参数
type SearchInstanceImageParam struct {
	RegionID      string `form:"regionId" binding:"required"`
	AccountID     string `form:"accountId"`
	ImageCategory string `form:"imageCategory"`
	ImageName     string `form:"imageName"`
	OSType        string `form:"osType"`
	OSName        string `form:"osName"`
}

// SearchSecurityGroupParam 搜索安全组参数
type SearchSecurityGroupParam struct {
	CloudName string `form:"cloudName"`
	RegionID  string `form:"regionId"`
	AccountID string `form:"accountId"`
}

type SearchDiskParam struct{}

type SearchKeypairParam struct{}
