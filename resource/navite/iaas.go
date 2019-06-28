package navite

import (
	"time"
)

// CloudImageTable 云商镜像表
const (
	ImageTable             = "images"
	InstanceTable          = "instances"
	InstanceSpecTable      = "instanceSpecs"
	DiskTable              = "disks"
	KeyPairTable           = "keypairs"
	SecurityGroupTable     = "securityGroups"
	SecurityGroupRuleTable = "securityGroupRules"
	VPCTable               = "vpcs"
	SubnetTable            = "subnets"
	EIPTable               = "eips"
)

// Image 云镜像
type Image struct {
	RegionID     string    `bson:"regionId" json:"regionId"`
	AccountID    string    `bson:"accountId" json:"accountId"` // 引用Account.ID
	CloudName    string    `bson:"cloudName" json:"cloudName"`
	ImageID      string    `bson:"imageId" json:"imageId"`
	ImageName    string    `bson:"imageName" json:"imageName"`
	ImageVersion string    `bson:"imageVersion" json:"imageVersion"`
	OSType       string    `bson:"osType" json:"osType"`
	OSName       string    `bson:"osName" json:"osName"`
	DiskSize     int       `bson:"diskSize" json:"diskSize"`
	Owner        string    `bson:"owner" json:"owner"` // 镜像所有者
	Description  string    `bson:"description" json:"description"`
	CreatedTime  time.Time `bson:"createdTime" json:"createdTime"`
	SyncedTime   time.Time `bson:"syncedTime" json:"syncedTime"` // 同步下来的时间
}

// InstanceSpec 实例规格
type InstanceSpec struct {
	CloudName        string    `bson:"cloudName" json:"cloudName"`
	AccountID        string    `bson:"accountId" json:"accountId"`
	RegionID         string    `bson:"regionId" json:"regionId"`
	ZoneID           string    `bson:"zoneId" json:"zoneId"`
	InstanceSpecID   string    `bson:"instanceSpecId" json:"instanceSpecId"`
	InstanceSpecName string    `bson:"instanceSpecName" json:"instanceSpecName"`
	InstanceFamily   string    `bson:"instanceFamily" json:"instanceFamily"`
	CPU              int       `bson:"cpu" json:"cpu"`
	Memory           float64   `bson:"memory" json:"memory"`
	Status           string    `bson:"status" json:"status"`
	SyncedTime       time.Time `bson:"syncedTime" json:"syncedTime"`
}

// Instance 计算实例
type Instance struct {
	CloudName         string    `bson:"cloudName" json:"cloudName"`
	AccountID         string    `bson:"accountId" json:"accountId"`
	RegionID          string    `bson:"regionId" json:"regionId"`
	ZoneID            string    `bson:"zoneId" json:"zoneId"`
	VPCID             string    `bson:"vpcId" json:"vpcId"`
	InstanceID        string    `bson:"instanceId" json:"instanceId"`
	InstanceName      string    `bson:"instanceName" json:"instanceName"`
	Status            string    `bson:"status" json:"status"`
	HostName          string    `bson:"hostname" json:"hostname"`
	CPU               int       `bson:"cpu" json:"cpu"`
	Memory            int       `bson:"memory" json:"memory"`
	OSName            string    `bson:"osName" json:"osName"`
	DeleteProtection  bool      `bson:"deleteProtection" json:"deleteProtection"` // 是否允许通过API控制
	Description       string    `bson:"description" json:"description"`
	EipAddress        string    `bson:"eipAddress" json:"eipAddress"`
	ImageID           string    `bson:"imageId" json:"imageId"`
	InnerIPAddress    string    `bson:"innerIpAddress" json:"innerIpAddress"`
	ChargeType        string    `bson:"chargeType" json:"chargeType"`
	InstanceType      string    `bson:"instanceType" json:"instanceType"`
	NetworkType       string    `bson:"networkType" json:"networkType"`
	KeyPairList       []string  `bson:"keypairList" json:"keypairList"`
	SecurityGroupList []string  `bson:"securityGroupList" json:"securityGroupList"`
	CreatedTime       time.Time `bson:"createdTime" json:"createdTime"`
	SyncedTime        time.Time `bson:"syncedTime" json:"syncedTime"`
}

// SecurityGroup 安全组
type SecurityGroup struct {
	CloudName   string    `bson:"cloudName" json:"cloudName"`
	AccountID   string    `bson:"accountId" json:"accountId"`
	RegionID    string    `bson:"regionId" json:"regionId"`
	GroupID     string    `bson:"groupId" json:"grouopId"`
	GroupName   string    `bson:"groupName" json:"groupName"`
	IsDefault   bool      `bson:"isDefault" json:"isDefault"`
	VPCID       string    `bson:"vpcId" json:"vpcId"` // ali空代表经典网络, 非空代表专有网络
	Description string    `bson:"description" json:"description"`
	CreatedTime time.Time `bson:"createdTime" json:"createdTime"`
	SyncedTime  time.Time `bson:"syncedTime" json:"syncedTime"`
}

// SecurityGroupRule 安全组规则
type SecurityGroupRule struct {
	CloudName    string    `bson:"cloudName" json:"cloudName"`
	GroupID      string    `bson:"groupId" json:"grouopId"`
	GroupName    string    `bson:"groupName" json:"groupName"`
	DestCidrIP   string    `bson:"destCidrIp" json:"destCidrIp"`
	SourceCidrIP string    `bson:"sourceCidrIp" json:"sourceCidrIp"`
	Direction    string    `bson:"direction" json:"direction"`
	Protocol     string    `bson:"protocol" json:"protocol"` // 支持: TCP, UDP, ICMP, GRE, ALL, *
	PortRange    string    `bson:"portRange" json:"portRange"`
	Priority     string    `bson:"priority" json:"priority"`
	Action       string    `bson:"action" json:"action"` // ACCEPT/DROP
	Description  string    `bson:"description" json:"description"`
	SyncedTime   time.Time `bson:"syncedTime" json:"syncedTime"`
}

// Disk 块存储
type Disk struct {
	CloudName        string    `bson:"cloudName" json:"cloudName"`
	RegionID         string    `bson:"regionId" json:"regionId"`
	AccountID        string    `bson:"accountId" json:"accountId"`
	ZoneID           string    `bson:"zoneId" json:"zoneId"`
	DiskID           string    `bson:"diskId" json:"diskId"`
	DiskName         string    `bson:"diskName" json:"diskName"`
	DiskType         string    `bson:"diskType" json:"diskType"`
	ChargeType       string    `bson:"chargeType" json:"chargeType"`
	IsEncrypted      bool      `bson:"isEncrypted" json:"isEncrypetd"`
	Shareable        bool      `bson:"shareable" json:"shareable"`
	DiskSize         int       `bson:"size" json:"size"`
	DiskIOPS         int       `bson:"iops" json:"iops"`
	Status           string    `bson:"status" json:"status"`
	AttachInstanceID string    `bson:"attachInstanceId" json:"attachInstanceId"`
	Device           string    `bson:"device" json:"device"`
	AttachedTime     time.Time `bson:"attachedTime" json:"attachedTime"`
	DetachedTime     time.Time `bson:"detachedTime" json:"detachedTime"`
	Description      string    `bson:"description" json:"description"`
	CreatedTime      time.Time `bson:"createdTime" json:"createdTime"`
	SyncedTime       time.Time `bson:"syncedTime" json:"syncedTime"`
}

// Keypair 密钥对
type Keypair struct {
	CloudName   string    `bson:"cloudName" json:"cloudName"`
	RegionID    string    `bson:"regionId" json:"regionId"`
	AccountID   string    `bson:"accountId" json:"accountId"`
	KeypairID   string    `bson:"keypairId" json:"keypairId"`
	KeypairName string    `bson:"keypairName" json:"keypairName"`
	PublicKey   string    `bson:"publicKey" json:"publicKey"`
	Description string    `bson:"description" json:"description"`
	CreatedTime time.Time `bson:"createdTime" json:"createdTime"`
	SyncedTime  time.Time `bson:"syncedTime" json:"syncedTime"`
}

// VPC vpc私有网络
type VPC struct {
	CloudName   string    `bson:"cloudName" json:"cloudName"`
	RegionID    string    `bson:"regionId" json:"regionId"`
	AccountID   string    `bson:"accountId" json:"accountId"`
	VPCID       string    `bson:"vpcId" json:"vpcId"`
	VPCName     string    `bson:"vpcName" json:"vpcName"`
	IsDefault   bool      `bson:"isDefault" json:"isDefault"`
	CidrBlock   string    `bson:"cidrBlock" json:"cidrBlock"`
	RouterID    string    `bson:"routerId" json:"routerId"`
	Status      string    `bson:"status" json:"status"`
	Description string    `bson:"description" json:"description"`
	CreatedTime time.Time `bson:"createdTime" json:"createdTime"`
	SyncedTime  time.Time `bson:"syncedTime" json:"syncedTime"`
}

// Subnet 私有网络中的子网
type Subnet struct {
	CloudName               string    `bson:"cloudName" json:"cloudName"`
	RegionID                string    `bson:"regionId" json:"regionId"`
	AccountID               string    `bson:"accountId" json:"accountId"`
	VPCID                   string    `bson:"vpcId" json:"vpcId"`
	SubnetID                string    `bson:"subnetId" json:"subnetId"`
	SubnetName              string    `bson:"subnetName" json:"subnetName"`
	CidrBlock               string    `bson:"cidrBlock" json:"cidrBlock"`
	IsDefault               bool      `bson:"isDefault" json:"isDefault"`
	EnableBroadcast         bool      `bson:"enableBroadcast" json:"enableBroadcast"`
	ZoneID                  string    `bson:"zoneId" json:"zoneId"`
	AvailableIPAddressCount int       `bson:"avaiableIpCount" json:"avaiableIpCount"`
	IsVPCSnat               bool      `bson:"isVpcSnat" json:"isVpcSnat"`
	Description             string    `bson:"description" json:"description"`
	CreatedTime             time.Time `bson:"createdTime" json:"createdTime"`
	SyncedTime              time.Time `bson:"syncedTime" json:"syncedTime"`
}

// Eip 弹性公网IP
type Eip struct {
	CloudName           string    `bson:"cloudName" json:"cloudName"`
	RegionID            string    `bson:"regionId" json:"regionId"`
	AccountID           string    `bson:"accountId" json:"accountId"`
	ZoneID              string    `bson:"zoneId" json:"zoneId"`
	ChargeType          string    `bson:"chargeType" json:"chargeType"`                   // 付费方式
	BandWidthChargeType string    `bson:"bandWidthChargeType" json:"bandWidthChargeType"` // 计费方式
	AddressID           string    `bson:"addressId" json:"addressId"`
	AddressName         string    `bson:"addressName" json:"addressName"`
	AddressStatus       string    `bson:"addressStatus" json:"addressStatus"`
	AddressIP           string    `bson:"addressIp" json:"addressIp"`
	BandWidth           int64     `bson:"bandWidth" json:"bandWidth"`
	BindInstanceID      string    `bson:"bindInstanceId" json:"bindInstanceId"`
	BindInstanceType    string    `bson:"bindInstanceType" json:"bindInstanceType"`
	NetworkInterfaceID  string    `bson:"networkInterfaceId" json:"networkInterfaceId"`
	AddressType         string    `bson:"addressType" json:"addressType"`
	Description         string    `bson:"description" json:"description"`
	CreatedTime         time.Time `bson:"createdTime" json:"createdTime"`
	SyncedTime          time.Time `bson:"syncedTime" json:"syncedTime"`
}
