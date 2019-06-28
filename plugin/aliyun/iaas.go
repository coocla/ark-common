package aliyun

import (
	"ark-common/constants"
	"ark-common/param"
	"ark-common/resource/navite"
	"ark-common/utils/tool"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	log "github.com/sirupsen/logrus"
)

// AliyunResource 阿里云驱动
type AliyunResource struct {
	client  *ecs.Client
	account *navite.CloudAccount
}

var rateLimit = map[string]int{
	constants.HandleSyncRegion:            100,
	constants.HandleSyncZone:              100,
	constants.HandleSyncInstanceSpec:      100,
	constants.HandleSyncImage:             100,
	constants.HandleSyncInstance:          100,
	constants.HandleSyncSecurityGroup:     100,
	constants.HandleSyncSecurityGroupRule: 100,
	constants.HandleSyncDisk:              100,
	constants.HandleSyncKeypair:           100,
	constants.HandleSyncVPC:               100,
	constants.HandleSyncSubnet:            100,
	constants.HandleSyncEip:               100,
}

// RateLimit 获取对应账号执行action的每秒并发数
func (ali *AliyunResource) RateLimit(action string) int {
	if rate, ok := rateLimit[action]; ok {
		return rate
	}
	// 默认并发数
	return 100
}

// NewAliyunPlugin 初始化阿里云驱动
func NewAliyunPlugin(ac *navite.CloudAccount) *AliyunResource {
	return &AliyunResource{
		client:  initClient(ac),
		account: ac,
	}
}

func initClient(ac *navite.CloudAccount) *ecs.Client {
	client, err := ecs.NewClientWithAccessKey(ac.RunRegionID, ac.AccessKey, ac.GetSK())
	if err != nil {
		log.Errorf("initialize clint failed: %v", err)
	}
	return client
}

// GetCloudName 返回云商名字
func (ali *AliyunResource) GetCloudName() string {
	return constants.Aliyun
}

// SyncJobs 返回自动同步的作业
func (ali *AliyunResource) SyncJobs() []string {
	return []string{
		// constants.HandleSyncZone,
		// constants.HandleSyncInstanceSpec,
		// constants.HandleSyncImage,
		// constants.HandleSyncInstance,
		// constants.HandleSyncDisk,
		// constants.HandleSyncKeypair,
		// constants.HandleSyncSecurityGroup,
		// constants.HandleSyncSecurityGroupRule,
		constants.HandleSyncVPC,
		constants.HandleSyncSubnet,
		constants.HandleSyncEip,
	}
}

// GetRegionList 获取地域列表
func (ali *AliyunResource) GetRegionList() (regionList []*navite.CloudRegion) {
	req := ecs.CreateDescribeRegionsRequest()
	resp, err := ali.client.DescribeRegions(req)
	if err != nil {
		log.Errorf("aliyun describe regions failed: %v", err)
		return
	}
	for _, res := range resp.Regions.Region {
		region := &navite.CloudRegion{
			RegionID:   res.RegionId,
			RegionName: res.LocalName,
			CloudName:  constants.Aliyun,
			SyncedTime: time.Now(),
		}
		regionList = append(regionList, region)
	}
	return regionList
}

// GetZoneList 获取可用区列表
func (ali *AliyunResource) GetZoneList() (zoneList []*navite.CloudZone) {
	req := ecs.CreateDescribeZonesRequest()
	resp, err := ali.client.DescribeZones(req)
	if err != nil {
		log.Errorf("aliyun describe zones failed: %v", err)
		return
	}
	for _, res := range resp.Zones.Zone {
		zone := &navite.CloudZone{
			CloudName:  constants.Aliyun,
			RegionID:   ali.account.RunRegionID,
			ZoneID:     res.ZoneId,
			ZoneName:   res.LocalName,
			SyncedTime: time.Now(),
		}
		zoneList = append(zoneList, zone)
	}
	return zoneList
}

// GetImageList 获取镜像列表
func (ali *AliyunResource) GetImageList(pageSize, currentPage int) (count int, imgs []*navite.Image) {
	req := ecs.CreateDescribeImagesRequest()
	req.PageSize = requests.NewInteger(pageSize)
	req.PageNumber = requests.NewInteger(currentPage)
	resp, err := ali.client.DescribeImages(req)
	if err != nil {
		log.Errorf("aliyun describe images failed: %v", err)
		return
	}
	for _, res := range resp.Images.Image {
		img := &navite.Image{
			RegionID:     resp.RegionId,
			AccountID:    ali.account.AccountID(),
			CloudName:    constants.Aliyun,
			ImageID:      res.ImageId,
			ImageName:    res.ImageName,
			DiskSize:     res.Size,
			Owner:        res.ImageOwnerAlias,
			OSType:       res.OSType,
			OSName:       res.OSName,
			ImageVersion: res.ImageVersion,
			Description:  res.Description,
			SyncedTime:   time.Now(),
		}
		imgs = append(imgs, img)
	}
	return resp.TotalCount, imgs
}

// GetInstanceList 获取实例列表
func (ali *AliyunResource) GetInstanceList(pageSize, currentPage int) (count int, instanceList []*navite.Instance) {
	req := ecs.CreateDescribeInstancesRequest()
	req.PageSize = requests.NewInteger(pageSize)
	req.PageNumber = requests.NewInteger(currentPage)
	resp, err := ali.client.DescribeInstances(req)
	if err != nil {
		log.Errorf("aliyun describe instance failed: %v", err)
		return
	}
	for _, res := range resp.Instances.Instance {
		instance := &navite.Instance{
			CloudName:         constants.Aliyun,
			AccountID:         ali.account.AccountID(),
			RegionID:          res.RegionId,
			ZoneID:            res.ZoneId,
			VPCID:             res.VpcAttributes.VpcId,
			InstanceID:        res.InstanceId,
			InstanceName:      res.InstanceName,
			InstanceType:      res.InstanceType,
			Status:            res.Status,
			HostName:          res.HostName,
			CPU:               res.Cpu,
			Memory:            res.Memory,
			OSName:            res.OSName,
			SecurityGroupList: res.SecurityGroupIds.SecurityGroupId,
			DeleteProtection:  res.DeletionProtection,
			Description:       res.Description,
			EipAddress:        res.EipAddress.IpAddress,
			ImageID:           res.ImageId,
			ChargeType:        res.InstanceChargeType,
			NetworkType:       res.InstanceNetworkType,
			KeyPairList:       []string{res.KeyPairName},
			CreatedTime:       tool.TimeForISO8601(res.CreationTime),
			SyncedTime:        time.Now(),
		}
		instanceList = append(instanceList, instance)
	}
	return resp.TotalCount, instanceList
}

// GetSecurityGroupList 获取安全组列表
func (ali *AliyunResource) GetSecurityGroupList(pageSize, currentPage int) (count int, sgList []*navite.SecurityGroup) {
	req := ecs.CreateDescribeSecurityGroupsRequest()
	req.PageSize = requests.NewInteger(pageSize)
	req.PageNumber = requests.NewInteger(currentPage)
	resp, err := ali.client.DescribeSecurityGroups(req)
	if err != nil {
		log.Errorf("aliyun describe securityGroup failed: %v", err)
		return
	}
	for _, res := range resp.SecurityGroups.SecurityGroup {
		sg := &navite.SecurityGroup{
			CloudName:   constants.Aliyun,
			GroupID:     res.SecurityGroupId,
			GroupName:   res.SecurityGroupName,
			RegionID:    ali.account.RunRegionID,
			AccountID:   ali.account.AccountID(),
			VPCID:       res.VpcId,
			Description: res.Description,
			CreatedTime: tool.TimeForISO8601(res.CreationTime),
			SyncedTime:  time.Now(),
		}
		sgList = append(sgList, sg)
	}
	return int(resp.TotalCount), sgList
}

// GetDiskList 获取磁盘列表
func (ali *AliyunResource) GetDiskList(pageSize, currentPage int) (count int, diskList []*navite.Disk) {
	req := ecs.CreateDescribeDisksRequest()
	req.PageSize = requests.NewInteger(pageSize)
	req.PageNumber = requests.NewInteger(currentPage)
	resp, err := ali.client.DescribeDisks(req)
	if err != nil {
		log.Errorf("aliyun describe disks failed: %v", err)
		return
	}
	for _, res := range resp.Disks.Disk {
		disk := &navite.Disk{
			CloudName:        constants.Aliyun,
			AccountID:        ali.account.AccountID(),
			RegionID:         ali.account.RunRegionID,
			DiskID:           res.DiskId,
			DiskName:         res.DiskName,
			DiskType:         res.Category,
			ChargeType:       res.DiskChargeType,
			IsEncrypted:      res.Encrypted,
			DiskSize:         res.Size,
			DiskIOPS:         res.IOPS,
			Status:           res.Status,
			AttachInstanceID: res.InstanceId,
			Device:           res.Device,
			AttachedTime:     tool.TimeForISO8601(res.AttachedTime),
			DetachedTime:     tool.TimeForISO8601(res.DetachedTime),
			Description:      res.Description,
			CreatedTime:      tool.TimeForISO8601(res.CreationTime),
			SyncedTime:       time.Now(),
		}
		diskList = append(diskList, disk)
	}
	return int(resp.TotalCount), diskList
}

// GetKeypairList 获取密钥对
func (ali *AliyunResource) GetKeypairList(pageSize, currentPage int) (count int, keypairList []*navite.Keypair) {
	req := ecs.CreateDescribeKeyPairsRequest()
	req.PageSize = requests.NewInteger(pageSize)
	req.PageNumber = requests.NewInteger(currentPage)
	resp, err := ali.client.DescribeKeyPairs(req)
	if err != nil {
		log.Errorf("aliyun describe keypairs failed: %v", err)
		return
	}
	for _, res := range resp.KeyPairs.KeyPair {
		keypair := &navite.Keypair{
			CloudName:   constants.Aliyun,
			AccountID:   ali.account.AccountID(),
			RegionID:    ali.account.RunRegionID,
			KeypairID:   res.KeyPairName,
			KeypairName: res.KeyPairName,
			SyncedTime:  time.Now(),
		}
		keypairList = append(keypairList, keypair)
	}
	return int(resp.TotalCount), keypairList
}

// GetSecurityGroupRuleList 获取安全组规则
func (ali *AliyunResource) GetSecurityGroupRuleList(securityGroupID string) (sgrList []*navite.SecurityGroupRule) {
	req := ecs.CreateDescribeSecurityGroupAttributeRequest()
	req.SecurityGroupId = securityGroupID
	resp, err := ali.client.DescribeSecurityGroupAttribute(req)
	if err != nil {
		log.Errorf("aliyun describe securityGroupRules failed: %v", err)
		return
	}
	for _, res := range resp.Permissions.Permission {
		sgr := &navite.SecurityGroupRule{
			CloudName:    constants.Aliyun,
			GroupID:      securityGroupID,
			GroupName:    res.DestGroupName,
			DestCidrIP:   res.DestCidrIp,
			SourceCidrIP: res.SourceCidrIp,
			Direction:    res.Direction,
			Protocol:     res.IpProtocol,
			PortRange:    res.PortRange,
			Priority:     res.Priority,
			Action:       strings.ToUpper(res.Policy),
			Description:  res.Description,
			SyncedTime:   time.Now(),
		}
		sgrList = append(sgrList, sgr)
	}
	return sgrList
}

// GetInstanceSpecsList 获取实例规格
func (ali *AliyunResource) GetInstanceSpecsList() (instantSpecList []*navite.InstanceSpec) {
	req := ecs.CreateDescribeInstanceTypesRequest()
	resp, err := ali.client.DescribeInstanceTypes(req)
	if err != nil {
		log.Errorf("aliyun describe instanceTypes failed: %v", err)
		return
	}
	for _, res := range resp.InstanceTypes.InstanceType {
		spec := &navite.InstanceSpec{
			CloudName:        constants.Aliyun,
			AccountID:        ali.account.AccountID(),
			RegionID:         ali.account.RunRegionID,
			InstanceSpecID:   res.InstanceTypeId,
			InstanceSpecName: res.InstanceType,
			InstanceFamily:   res.InstanceTypeFamily,
			CPU:              res.CpuCoreCount,
			Memory:           res.MemorySize,
			SyncedTime:       time.Now(),
		}
		instantSpecList = append(instantSpecList, spec)
	}
	return
}

// GetVPCList 获取VPC列表
func (ali *AliyunResource) GetVPCList(pageSize, currentPage int) (count int, vpcList []*navite.VPC) {
	req := ecs.CreateDescribeVpcsRequest()
	req.PageSize = requests.NewInteger(pageSize)
	req.PageNumber = requests.NewInteger(currentPage)
	resp, err := ali.client.DescribeVpcs(req)
	if err != nil {
		log.Errorf("aliyun describe vpcs failed: %v", err)
		return
	}
	for _, res := range resp.Vpcs.Vpc {
		v := &navite.VPC{
			CloudName:   constants.Aliyun,
			AccountID:   ali.account.AccountID(),
			RegionID:    ali.account.RunRegionID,
			VPCID:       res.VpcId,
			VPCName:     res.VpcName,
			IsDefault:   res.IsDefault,
			CidrBlock:   res.CidrBlock,
			RouterID:    res.VRouterId,
			Status:      res.Status,
			Description: res.Description,
			CreatedTime: tool.TimeForISO8601(res.CreationTime),
			SyncedTime:  time.Now(),
		}
		vpcList = append(vpcList, v)
	}
	return resp.TotalCount, vpcList
}

// GetSubnetList 获取子网列表
func (ali *AliyunResource) GetSubnetList(pageSize, currentPage int) (count int, subnetList []*navite.Subnet) {
	req := ecs.CreateDescribeVSwitchesRequest()
	req.PageSize = requests.NewInteger(pageSize)
	req.PageNumber = requests.NewInteger(currentPage)
	resp, err := ali.client.DescribeVSwitches(req)
	if err != nil {
		log.Errorf("aliyun describe vswitch failed: %v", err)
		return
	}
	for _, res := range resp.VSwitches.VSwitch {
		subnet := &navite.Subnet{
			CloudName:               constants.Aliyun,
			AccountID:               ali.account.AccountID(),
			RegionID:                ali.account.RunRegionID,
			VPCID:                   res.VpcId,
			SubnetID:                res.VSwitchId,
			SubnetName:              res.VSwitchName,
			CidrBlock:               res.CidrBlock,
			IsDefault:               res.IsDefault,
			ZoneID:                  res.ZoneId,
			AvailableIPAddressCount: res.AvailableIpAddressCount,
			Description:             res.Description,
			CreatedTime:             tool.TimeForISO8601(res.CreationTime),
			SyncedTime:              time.Now(),
		}
		subnetList = append(subnetList, subnet)
	}
	return resp.TotalCount, subnetList
}

// GetEipList 获取弹性公网IP列表
func (ali *AliyunResource) GetEipList(pageSize, currentPage int) (count int, eipList []*navite.Eip) {
	req := ecs.CreateDescribeEipAddressesRequest()
	req.PageSize = requests.NewInteger(pageSize)
	req.PageNumber = requests.NewInteger(currentPage)
	resp, err := ali.client.DescribeEipAddresses(req)
	if err != nil {
		log.Errorf("aliyun describe eip failed: %v", err)
		return
	}
	for _, res := range resp.EipAddresses.EipAddress {
		var bandWidth int
		bandWidth, _ = strconv.Atoi(res.Bandwidth)
		eip := &navite.Eip{
			CloudName:        constants.Aliyun,
			AccountID:        ali.account.AccountID(),
			RegionID:         ali.account.RunRegionID,
			ChargeType:       res.ChargeType,
			AddressID:        res.AllocationId,
			AddressName:      "",
			AddressStatus:    res.Status,
			AddressIP:        res.IpAddress,
			BandWidth:        int64(bandWidth),
			BindInstanceID:   res.InstanceId,
			BindInstanceType: res.InstanceType,
			CreatedTime:      tool.TimeForISO8601(res.AllocationTime),
			SyncedTime:       time.Now(),
		}
		eipList = append(eipList, eip)
	}
	return resp.TotalCount, eipList
}

// NewKeypair 创建密钥对
func (ali *AliyunResource) NewKeypair(keypair *navite.Keypair) (err error) {
	req := ecs.CreateImportKeyPairRequest()
	req.KeyPairName = keypair.KeypairName
	req.PublicKeyBody = keypair.PublicKey
	resp, err := ali.client.ImportKeyPair(req)
	if err != nil {
		log.Errorf("aliyun import keypair [%s] failed: %v", req.GetQueryParams(), err)
		return err
	}
	keypair.KeypairID = resp.KeyPairName
	return
}

// DeleteKeypair 删除密钥对
func (ali *AliyunResource) DeleteKeypair(keypairIDList ...string) (err error) {
	req := ecs.CreateDeleteKeyPairsRequest()
	b, _ := json.Marshal(keypairIDList)
	req.KeyPairNames = string(b)
	_, err = ali.client.DeleteKeyPairs(req)
	if err != nil {
		log.Errorf("aliyun delete keypair [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// NewSecurityGroup 创建安全组
func (ali *AliyunResource) NewSecurityGroup(sg *navite.SecurityGroup) (err error) {
	req := ecs.CreateCreateSecurityGroupRequest()
	req.SecurityGroupName = sg.GroupName
	req.Description = sg.Description
	req.VpcId = sg.VPCID
	resp, err := ali.client.CreateSecurityGroup(req)
	if err != nil {
		log.Errorf("aliyun create securityGroup [%s] failed: %v", req.GetQueryParams(), err)
		return
	}
	sg.GroupID = resp.SecurityGroupId
	return nil
}

// DeleteSecurityGroup 删除安全组
func (ali *AliyunResource) DeleteSecurityGroup(sgID string) (err error) {
	req := ecs.CreateDeleteSecurityGroupRequest()
	req.SecurityGroupId = sgID
	_, err = ali.client.DeleteSecurityGroup(req)
	if err != nil {
		log.Errorf("aliyun delete securityGroup [%s] failed: %v", req.GetQueryParams(), err)
	}
	return
}

// NewSecurityGroupRule 创建安全组规则
//
// * PortRange 需要 ?/? 格式
func (ali *AliyunResource) NewSecurityGroupRule(rule *navite.SecurityGroupRule) (err error) {
	switch rule.Direction {
	case constants.FlowIngress:
		ingressReq := ecs.CreateAuthorizeSecurityGroupRequest()
		ingressReq.SourceCidrIp = rule.SourceCidrIP
		ingressReq.PortRange = rule.PortRange
		ingressReq.IpProtocol = rule.Protocol
		ingressReq.Priority = rule.Priority
		ingressReq.Policy = rule.Action
		ingressReq.SecurityGroupId = rule.GroupID
		if _, err = ali.client.AuthorizeSecurityGroup(ingressReq); err != nil {
			log.Errorf("aliyun create securityGroupRule [%s] failed: %v", ingressReq.GetQueryParams(), err)
			return err
		}
	case constants.FlowEgress:
		egressReq := ecs.CreateAuthorizeSecurityGroupEgressRequest()
		egressReq.DestCidrIp = rule.DestCidrIP
		egressReq.PortRange = rule.PortRange
		egressReq.IpProtocol = rule.Protocol
		egressReq.Priority = rule.Priority
		egressReq.Policy = rule.Action
		egressReq.SecurityGroupId = rule.GroupID
		if _, err = ali.client.AuthorizeSecurityGroupEgress(egressReq); err != nil {
			log.Errorf("aliyun create securityGroupRule [%s] failed: %v", egressReq.GetQueryParams(), err)
			return err
		}
	}
	return nil
}

// DeleteSecurityGroupRule 删除安全组规则
func (ali *AliyunResource) DeleteSecurityGroupRule(rule *navite.SecurityGroupRule) (err error) {
	switch rule.Direction {
	case constants.FlowIngress:
		req := ecs.CreateRevokeSecurityGroupRequest()
		req.IpProtocol = rule.Protocol
		req.PortRange = rule.PortRange
		req.NicType = "internet"
		req.Policy = rule.Action
		req.SourceCidrIp = rule.SourceCidrIP
		req.SecurityGroupId = rule.GroupID
		if _, err = ali.client.RevokeSecurityGroup(req); err != nil {
			log.Errorf("aliyun revoke ingress securityGroupRule [%s] failed: %v", req.GetQueryParams(), err)
			return err
		}
	case constants.FlowEgress:
		req := ecs.CreateRevokeSecurityGroupEgressRequest()
		req.IpProtocol = rule.Protocol
		req.PortRange = rule.PortRange
		req.NicType = "internet"
		req.Policy = rule.Action
		req.DestCidrIp = rule.DestCidrIP
		req.SecurityGroupId = rule.GroupID
		if _, err = ali.client.RevokeSecurityGroupEgress(req); err != nil {
			log.Errorf("aliyun revoke egress securityGroupRule [%s] failed: %v", req.GetQueryParams(), err)
			return err
		}
	}
	return nil
}

// NewVPC 创建虚拟专用网络
func (ali *AliyunResource) NewVPC(vpc *navite.VPC) (err error) {
	req := ecs.CreateCreateVpcRequest()
	req.VpcName = vpc.VPCName
	req.CidrBlock = vpc.CidrBlock
	req.Description = vpc.Description
	resp, err := ali.client.CreateVpc(req)
	if err != nil {
		log.Errorf("aliyun create vpc [%s] failed: %v", req.GetQueryParams(), err)
		return
	}
	vpc.VPCID = resp.VpcId
	return
}

// DeleteVPC 删除虚拟专用网络
func (ali *AliyunResource) DeleteVPC(vpcID string) (err error) {
	req := ecs.CreateDeleteVpcRequest()
	req.VpcId = vpcID
	_, err = ali.client.DeleteVpc(req)
	if err != nil {
		log.Errorf("aliyun delete vpc [%s] failed: %v", req.GetQueryParams(), err)
	}
	return
}

// NewSubnet 创建子网
func (ali *AliyunResource) NewSubnet(subnet *navite.Subnet) (err error) {
	req := ecs.CreateCreateVSwitchRequest()
	req.CidrBlock = subnet.CidrBlock
	req.VpcId = subnet.VPCID
	req.ZoneId = subnet.ZoneID
	req.VSwitchName = subnet.SubnetName
	req.Description = subnet.Description
	resp, err := ali.client.CreateVSwitch(req)
	if err != nil {
		log.Errorf("aliyun create vswitch [%s] failed: %v", req.GetQueryParams(), err)
		return
	}
	subnet.SubnetID = resp.VSwitchId
	return
}

// DeleteSubnet 删除子网
func (ali *AliyunResource) DeleteSubnet(subnetID string) (err error) {
	req := ecs.CreateDeleteVSwitchRequest()
	req.VSwitchId = subnetID
	_, err = ali.client.DeleteVSwitch(req)
	if err != nil {
		log.Errorf("aliyun delete vswitch [%s] failed: %v", req.GetQueryParams(), err)
	}
	return
}

// NewDisk 创建云盘
func (ali *AliyunResource) NewDisk(disk *navite.Disk) (err error) {
	req := ecs.CreateCreateDiskRequest()
	req.DiskName = disk.DiskName
	req.Description = disk.Description
	req.DiskCategory = disk.DiskType
	req.Size = requests.NewInteger(disk.DiskSize)
	req.Encrypted = requests.NewBoolean(disk.IsEncrypted)
	req.ZoneId = disk.ZoneID
	resp, err := ali.client.CreateDisk(req)
	if err != nil {
		log.Errorf("aliyun create disk [%s] failed: %v", req.GetQueryParams(), err)
		return
	}
	disk.DiskID = resp.DiskId
	return
}

// DeleteDisk 删除云盘
//
// * 注意: 这里删除参数中第一个ID, 不会为其他ID发起删除请求
func (ali *AliyunResource) DeleteDisk(diskIDList ...string) (err error) {
	req := ecs.CreateDeleteDiskRequest()
	req.DiskId = diskIDList[0]
	_, err = ali.client.DeleteDisk(req)
	if err != nil {
		log.Errorf("aliyun delete disk [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// NewEIP 申请弹性公网IP
func (ali *AliyunResource) NewEIP(eip *navite.Eip) (err error) {
	req := ecs.CreateAllocateEipAddressRequest()
	req.Bandwidth = strconv.Itoa(int(eip.BandWidth))
	req.InternetChargeType = eip.BandWidthChargeType // 带宽的计费方式
	resp, err := ali.client.AllocateEipAddress(req)
	if err != nil {
		log.Errorf("aliyun create eip [%s] failed: %v", req.GetQueryParams(), err)
		return err
	}
	eip.AddressID = resp.AllocationId
	eip.AddressIP = resp.EipAddress
	return
}

// ReleaseEIP 释放弹性公网IP
//
// * 注意: 这里只删除参数中第一个ID, 不会为其他ID发起删除请求
func (ali *AliyunResource) ReleaseEIP(eipIDList ...string) (err error) {
	req := ecs.CreateReleaseEipAddressRequest()
	req.AllocationId = eipIDList[0]
	_, err = ali.client.ReleaseEipAddress(req)
	if err != nil {
		log.Errorf("aliyun release eip [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// RunInstance 创建实例
func (ali *AliyunResource) RunInstance(instance *param.RunInstanceParam) (instanceIDList []string, err error) {
	req := ecs.CreateRunInstancesRequest()
	// 1. 镜像, 机型(配置)
	req.ImageId = instance.ImageID           // 镜像，系统
	req.InstanceType = instance.InstanceType // 机型，内存/CPU
	// 2. 登陆sshkey
	req.KeyPairName = instance.KeyPairID

	// 3. 安全组, 网络
	req.SecurityGroupId = instance.SecurityGroupID
	req.VSwitchId = instance.SubnetID

	// 4. 磁盘设置
	if instance.DiskSize != 0 {
		dataDiskList := []ecs.RunInstancesDataDisk{
			ecs.RunInstancesDataDisk{
				Size:     strconv.Itoa(instance.DiskSize),
				DiskName: instance.InstanceName,
				Category: instance.DiskType,
			},
		}
		req.DataDisk = &dataDiskList
	}

	// 5. 自定义设置
	req.InstanceName = instance.InstanceName
	req.HostName = instance.HostName

	req.Amount = requests.NewInteger(instance.Numbers)

	resp, err := ali.client.RunInstances(req)
	if err != nil {
		log.Errorf("aliyun runInstance [%s] failed: %v", req.GetQueryParams(), err)
		return nil, err
	}
	return resp.InstanceIdSets.InstanceIdSet, nil
}

// DeleteInstance 删除实例
//
// * 此接口不能批量删除
func (ali *AliyunResource) DeleteInstance(instanceIDList ...string) (err error) {
	req := ecs.CreateDeleteInstanceRequest()
	req.InstanceId = instanceIDList[0]
	req.Force = requests.NewBoolean(true)
	_, err = ali.client.DeleteInstance(req)
	if err != nil {
		log.Errorf("aliyun deleteInstance [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// StartInstance 启动实例
//
// * 此接口不能批量操作
func (ali *AliyunResource) StartInstance(instanceIDList ...string) (err error) {
	req := ecs.CreateStartInstanceRequest()
	req.InstanceId = instanceIDList[0]
	_, err = ali.client.StartInstance(req)
	if err != nil {
		log.Errorf("aliyun startInstance [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// StopInstance 停止实例
//
// * 此接口不能批量操作
func (ali *AliyunResource) StopInstance(instanceIDList ...string) (err error) {
	req := ecs.CreateStopInstanceRequest()
	req.InstanceId = instanceIDList[0]
	_, err = ali.client.StopInstance(req)
	if err != nil {
		log.Errorf("aliyun stopInstance [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// RebotInstance 停止实例
//
// * 此接口不能批量操作
func (ali *AliyunResource) RebotInstance(instanceIDList ...string) (err error) {
	req := ecs.CreateRebootInstanceRequest()
	req.InstanceId = instanceIDList[0]
	_, err = ali.client.RebootInstance(req)
	if err != nil {
		log.Errorf("aliyun rebotInstance [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// AttachDisk 挂载磁盘到实例上
func (ali *AliyunResource) AttachDisk(instance *navite.Instance, disk *navite.Disk) (err error) {
	req := ecs.CreateAttachDiskRequest()
	req.DiskId = disk.DiskID
	req.InstanceId = instance.InstanceID
	_, err = ali.client.AttachDisk(req)
	if err != nil {
		log.Errorf("aliyun attachDisk [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// DetachDisk 解挂载磁盘
func (ali *AliyunResource) DetachDisk(instance *navite.Instance, disk *navite.Disk) (err error) {
	req := ecs.CreateDetachDiskRequest()
	req.DiskId = disk.DiskID
	req.InstanceId = instance.InstanceID
	_, err = ali.client.DetachDisk(req)
	if err != nil {
		log.Errorf("aliyun detachDisk [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// AttachEipToInstance 绑定弹性公网IP到实例上
func (ali *AliyunResource) AttachEipToInstance(instance *navite.Instance, eip *navite.Eip) (err error) {
	req := ecs.CreateAssociateEipAddressRequest()
	req.AllocationId = eip.AddressID
	req.InstanceId = instance.InstanceID
	req.InstanceType = "Ecs" // 可以取值: Nat|Slb|Ecs
	_, err = ali.client.AssociateEipAddress(req)
	if err != nil {
		log.Errorf("aliyun attachEipToInstance [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// DetachEipFromInstance 从实例上解绑弹性公网IP
func (ali *AliyunResource) DetachEipFromInstance(instance *navite.Instance, eip *navite.Eip) (err error) {
	req := ecs.CreateUnassociateEipAddressRequest()
	req.AllocationId = eip.AddressID
	req.InstanceId = instance.InstanceID
	req.InstanceType = "Ecs" // 可以取值: Nat|Slb|Ecs
	_, err = ali.client.UnassociateEipAddress(req)
	if err != nil {
		log.Errorf("aliyun detachEipFormInstance [%s] failed: %v", req.GetQueryParams(), err)
	}
	return err
}

// ModifyEIPBandWidth 调整弹性公网IP的带宽
func (ali *AliyunResource) ModifyEIPBandWidth(eip *navite.Eip, bandWidth int64) (err error) {
	req := ecs.CreateModifyEipAddressAttributeRequest()
	req.Bandwidth = strconv.Itoa(int(bandWidth))
	_, err = ali.client.ModifyEipAddressAttribute(req)
	if err != nil {
		log.Errorf("aliyun modifyEIPBandWidth [%s] failed: %v", req.GetQueryParams(), err)
		return err
	}
	eip.BandWidth = bandWidth
	return
}
