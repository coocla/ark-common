package plugin

import (
	"ark-common/clients/mgo"
	"ark-common/clients/redis"
	"ark-common/constants"
	"ark-common/param"
	"ark-common/plugin/aliyun"
	"ark-common/plugin/tencent"
	"ark-common/resource/navite"
	"time"

	log "github.com/sirupsen/logrus"
)

// AccountDriver 云账号接口
type AccountDriver interface {
	BindAccount(accountName, ak, sk string) *navite.CloudAccount // 绑定云账号
}

// ResourceDriver 云商资源接口
type ResourceDriver interface {
	RateLimit(action string) int // 返回接口限速
	GetCloudName() string        // 返回插件所属的云商名
	SyncJobs() []string          // 返回资源同步的作业名

	GetRegionList() (regionList []*navite.CloudRegion)                        // 同步地域
	GetZoneList() (zoneList []*navite.CloudZone)                              // 同步可用区
	GetInstanceSpecsList() (instantSpecList []*navite.InstanceSpec)           // 同步实例类型
	GetImageList(pageSize, currentPage int) (count int, imgs []*navite.Image) // 同步镜像

	GetInstanceList(pageSize, currentPage int) (count int, instanceList []*navite.Instance)     // 同步计算实例
	GetSecurityGroupList(pageSize, currentPage int) (count int, sgList []*navite.SecurityGroup) // 同步安全组
	GetSecurityGroupRuleList(securityGroupID string) (sgrList []*navite.SecurityGroupRule)      // 同步安全组规则
	GetDiskList(pageSize, currentPage int) (count int, diskList []*navite.Disk)                 // 同步磁盘
	GetKeypairList(pageSize, currentPage int) (count int, keypairList []*navite.Keypair)        // 同步密钥对
	GetVPCList(pageSize, currentPage int) (count int, vpcList []*navite.VPC)                    // 同步VPC
	GetSubnetList(pageSize, currentPage int) (count int, subnetList []*navite.Subnet)           // 同步子网
	GetEipList(pageSize, currentPage int) (count int, eipList []*navite.Eip)                    // 同步弹性公网

	NewKeypair(keypair *navite.Keypair) (err error)                                    // 创建密钥对
	DeleteKeypair(keypairIDList ...string) (err error)                                 // 删除密钥对
	NewSecurityGroup(sg *navite.SecurityGroup) (err error)                             // 创建安全组
	DeleteSecurityGroup(sgID string) (err error)                                       // 删除安全组
	NewSecurityGroupRule(rule *navite.SecurityGroupRule) (err error)                   // 创建安全组规则
	DeleteSecurityGroupRule(rule *navite.SecurityGroupRule) (err error)                // 删除安全组规则
	NewVPC(vpc *navite.VPC) (err error)                                                // 创建虚拟专用网
	DeleteVPC(vpcID string) (err error)                                                // 删除虚拟专用网
	NewSubnet(subnet *navite.Subnet) (err error)                                       // 创建子网
	DeleteSubnet(subnetID string) (err error)                                          // 删除子网
	NewDisk(disk *navite.Disk) (err error)                                             // 创建磁盘
	DeleteDisk(diskIDList ...string) (err error)                                       // 删除磁盘
	NewEIP(eip *navite.Eip) (err error)                                                // 申请弹性公网IP
	ReleaseEIP(eipIDList ...string) (err error)                                        // 释放弹性公网IP
	ModifyEIPBandWidth(eip *navite.Eip, bandWidth int64) (err error)                   // 调整弹性公网IP的带宽
	RunInstance(instance *param.RunInstanceParam) (instanceIDList []string, err error) // 创建实例
	DeleteInstance(instanceIDList ...string) (err error)                               // 删除实例
	StartInstance(instanceIDList ...string) (err error)                                // 启动实例
	StopInstance(instanceIDList ...string) (err error)                                 // 停止实例
	RebotInstance(instanceIDList ...string) (err error)                                // 重启实例
	AttachDisk(instance *navite.Instance, disk *navite.Disk) (err error)               // 挂载磁盘
	DetachDisk(instance *navite.Instance, disk *navite.Disk) (err error)               // 卸载磁盘
	AttachEipToInstance(instance *navite.Instance, eip *navite.Eip) (err error)        // 绑定弹性公网IP到实例上
	DetachEipFromInstance(instance *navite.Instance, eip *navite.Eip) (err error)      // 从实例上解绑弹性公网IP
}

// GetCloudDriver 返回对应的云商资源驱动
func GetCloudDriver(ac *navite.CloudAccount) ResourceDriver {
	if ac == nil {
		return nil
	}
	switch ac.CloudName {
	case constants.Aliyun:
		return aliyun.NewAliyunPlugin(ac)
	case constants.Tencent:
		return tencent.NewTencentPlugin(ac)
	}
	log.Errorf("not support cloud %s", ac.CloudName)
	return nil
}

// GetCloudAccountDriver 根据云账号返回对应的云商账号驱动
func GetCloudAccountDriver(rbd *mgo.Client, cloudName string) AccountDriver {
	switch cloudName {
	case constants.Aliyun:
		return aliyun.NewAliyunAccountPlugin(rbd)
	case constants.Tencent:
		return tencent.NewTencentAccountPlugin(rbd)
	}
	log.Errorf("not support cloud %s", cloudName)
	return nil
}

// CheckRateLimit 检查对应账号指定动作的限速额度
func CheckRateLimit(ac *navite.CloudAccount, action string) bool {
	driver := GetCloudDriver(ac)
	rateLimit := driver.RateLimit(action)
	limiter := redis.InitQuota(nil)
	// key Example: aliyun-1-SyncEip
	key := ac.CloudName + "-" + ac.AccountID() + "-" + action
	_, _, allowed := limiter.AllowN(key, int64(rateLimit), time.Second, 1)
	return allowed
}
