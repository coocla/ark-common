package navite

import "time"

// ServiceTable 服务表
const ServiceTable = "service"

// Service 服务
type Service struct {
	ServiceName    string    `bson:"serviceName" json:"serviceName"`
	Status         string    `bson:"status" json:"status"`
	LastUpdateTime time.Time `bson:"lastUpdateTime" json:"lastUpdateTime"`
}
