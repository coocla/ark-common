package navite

import (
	"ark-common/utils/tool"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CloudAccountTable 云商账号表
const (
	CloudAccountTable = "cloudAccount"
	CloudRegionTable  = "cloudRegions"
	CloudZoneTable    = "cloudZones"
)

// CloudAccount 云平台账户
type CloudAccount struct {
	ID          primitive.ObjectID `bson:"_id" json:"accountId"`
	AccountName string             `bson:"accountName" json:"accountName"`
	CloudName   string             `bson:"cloudName" json:"cloudName"`
	AccessKey   string             `bson:"accessKey" json:"-"`
	SecurityKey string             `bson:"securityKey" json:"-"`
	Balance     float64            `bson:"balance" json:"balance"`
	Disabled    bool               `bson:"disabled" json:"disabled"`
	Healthy     bool               `bson:"healthy" json:"healthy"`
	Status      string             `bson:"status" json:"status"`
	RunRegionID string             `json:"-"`
	CreatedTime time.Time          `bson:"createdTime" json:"createdTime"`
}

// CloudRegion 地域
type CloudRegion struct {
	RegionID   string    `bson:"regionId" json:"regionId"`
	RegionName string    `bson:"regionName" json:"regionName"`
	CloudName  string    `bson:"cloudName" json:"cloudName"`
	SyncedTime time.Time `bson:"syncedTime" json:"syncedTime"`
}

// CloudZone 可用区
type CloudZone struct {
	RegionID   string    `bson:"regionId" json:"regionId"`
	ZoneID     string    `bson:"zoneId" json:"zoneId"`
	ZoneName   string    `bson:"zoneName" json:"zoneName"`
	CloudName  string    `bson:"cloudName" json:"cloudName"`
	SyncedTime time.Time `bson:"syncedTime" json:"syncedTime"`
}

// AccountID 返回账号ID
func (c *CloudAccount) AccountID() string {
	return c.ID.Hex()
}

// GetSK 获取秘钥
func (c *CloudAccount) GetSK() string {
	return tool.DCP(c.SecurityKey)
}

// Encryption 设置秘钥
func (c *CloudAccount) Encryption() {
	c.SecurityKey = tool.ECP(c.SecurityKey)
}
