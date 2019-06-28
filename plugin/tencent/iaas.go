package tencent

import (
	"ark-common/constants"
	"ark-common/param"
	"ark-common/resource/navite"
	"ark-common/utils/tool"
	"fmt"
	"strconv"
	"strings"
	"time"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	log "github.com/sirupsen/logrus"
)

// TencentResource 腾讯云驱动
type TencentResource struct {
	cvm     *cvm.Client
	vpc     *vpc.Client
	cbs     *cbs.Client
	account *navite.CloudAccount
}

var rateLimit = map[string]int{
	constants.HandleSyncRegion:            20,  // https://cloud.tencent.com/document/api/213/15708
	constants.HandleSyncZone:              20,  // https://cloud.tencent.com/document/api/213/15707
	constants.HandleSyncInstanceSpec:      40,  // https://cloud.tencent.com/document/api/213/17378
	constants.HandleSyncImage:             40,  // https://cloud.tencent.com/document/api/213/15715
	constants.HandleSyncInstance:          40,  // https://cloud.tencent.com/document/api/213/15728
	constants.HandleSyncSecurityGroup:     100, // https://cloud.tencent.com/document/api/215/15808
	constants.HandleSyncSecurityGroupRule: 100, // https://cloud.tencent.com/document/api/215/15804
	constants.HandleSyncDisk:              20,  // https://cloud.tencent.com/document/api/362/16315
	constants.HandleSyncKeypair:           10,  // https://cloud.tencent.com/document/api/213/15699
	constants.HandleSyncVPC:               100, // https://cloud.tencent.com/document/api/215/15778
	constants.HandleSyncSubnet:            100, // https://cloud.tencent.com/document/api/215/15784
	constants.HandleSyncEip:               10,  // https://cloud.tencent.com/document/api/215/16702
	constants.HandleCreateEip:             10,  // https://cloud.tencent.com/document/api/215/16699
}

// RateLimit 获取对应账号执行action的每秒并发数
func (ten *TencentResource) RateLimit(action string) int {
	if rate, ok := rateLimit[action]; ok {
		return rate
	}
	// 默认并发数
	return 100
}

// NewTencentPlugin 初始化阿里云驱动
func NewTencentPlugin(ac *navite.CloudAccount) *TencentResource {
	credential := common.NewCredential(ac.AccessKey, ac.GetSK())
	client := &TencentResource{
		account: ac,
	}
	client.Connect(credential)
	return client
}

// Connect 初始化客户端连接
func (ten *TencentResource) Connect(credential *common.Credential) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	cvm, err := cvm.NewClient(credential, ten.account.RunRegionID, cpf)
	if err != nil {
		log.Errorf("inititenze cvm client failed: %v", err)
	}
	cpf = profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	vpc, err := vpc.NewClient(credential, ten.account.RunRegionID, cpf)
	if err != nil {
		log.Errorf("inititenze vpc client failed: %v", err)
	}
	cpf = profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cbs.tencentcloudapi.com"
	cbs, err := cbs.NewClient(credential, ten.account.RunRegionID, cpf)
	if err != nil {
		log.Errorf("inititenze cbs client failed: %v", err)
	}
	ten.cvm = cvm
	ten.vpc = vpc
	ten.cbs = cbs
}

// GetCloudName 返回云商名字
func (ten *TencentResource) GetCloudName() string {
	return constants.Tencent
}

// SyncJobs 返回自动同步的作业
func (ten *TencentResource) SyncJobs() []string {
	return []string{
		constants.HandleSyncZone,
		constants.HandleSyncInstanceSpec,
		constants.HandleSyncImage,
		constants.HandleSyncInstance,
		constants.HandleSyncDisk,
		constants.HandleSyncKeypair,
		constants.HandleSyncSecurityGroup,
		constants.HandleSyncSecurityGroupRule,
		constants.HandleSyncVPC,
		constants.HandleSyncSubnet,
		constants.HandleSyncEip,
	}
}

// GetPageLimitUint64 获取分页参数
func GetPageLimitUint64(pageSize, currentPage int) (limit, offset *uint64) {
	_limit := uint64(pageSize)
	if currentPage > 0 {
		currentPage = currentPage - 1
	}
	_offset := uint64(pageSize * currentPage)
	return &_limit, &_offset
}

// GetPageLimitInt64 获取分页参数，返回int64
func GetPageLimitInt64(pageSize, currentPage int) (limit, offset *int64) {
	_limit := int64(pageSize)
	if currentPage > 0 {
		currentPage = currentPage - 1
	}
	_offset := int64(pageSize * currentPage)
	return &_limit, &_offset
}

// GetPageLimitString 获取分页参数，返回string
func GetPageLimitString(pageSize, currentPage int) (limit, offset *string) {
	_limit := strconv.Itoa(pageSize)
	if currentPage > 0 {
		currentPage = currentPage - 1
	}
	_offset := strconv.Itoa(pageSize * currentPage)
	return &_limit, &_offset
}

// GetRegionList 获取地域列表
func (ten *TencentResource) GetRegionList() (regionList []*navite.CloudRegion) {
	req := cvm.NewDescribeRegionsRequest()
	resp, err := ten.cvm.DescribeRegions(req)
	if err != nil {
		log.Errorf("tencent describe regions failed: %v", err)
		return
	}
	for _, res := range resp.Response.RegionSet {
		region := &navite.CloudRegion{
			RegionID:   *res.Region,
			RegionName: *res.RegionName,
			CloudName:  constants.Tencent,
			SyncedTime: time.Now(),
		}
		regionList = append(regionList, region)
	}
	return regionList
}

// GetZoneList 获取可用区列表
func (ten *TencentResource) GetZoneList() (zoneList []*navite.CloudZone) {
	req := cvm.NewDescribeZonesRequest()
	resp, err := ten.cvm.DescribeZones(req)
	if err != nil {
		log.Errorf("tencent describe zones failed: %v", err)
		return
	}
	for _, res := range resp.Response.ZoneSet {
		zone := &navite.CloudZone{
			CloudName:  constants.Tencent,
			RegionID:   ten.account.RunRegionID,
			ZoneID:     *res.Zone,
			ZoneName:   *res.ZoneName,
			SyncedTime: time.Now(),
		}
		zoneList = append(zoneList, zone)
	}
	return zoneList
}

// GetImageList 获取镜像列表
func (ten *TencentResource) GetImageList(pageSize, currentPage int) (count int, imgs []*navite.Image) {
	req := cvm.NewDescribeImagesRequest()
	req.Limit, req.Offset = GetPageLimitUint64(pageSize, currentPage)
	resp, err := ten.cvm.DescribeImages(req)
	if err != nil {
		log.Errorf("tencent describe images failed: %v", err)
		return
	}
	for _, res := range resp.Response.ImageSet {
		fmt.Printf("%+v\n", res)
		img := &navite.Image{
			RegionID:    ten.account.RunRegionID,
			AccountID:   ten.account.AccountID(),
			CloudName:   constants.Tencent,
			ImageID:     *res.ImageId,
			ImageName:   *res.ImageName,
			DiskSize:    int(*res.ImageSize),
			OSType:      *res.Platform,
			OSName:      *res.OsName,
			Description: *res.ImageDescription,
			SyncedTime:  time.Now(),
		}
		imgs = append(imgs, img)
	}
	count = int(*resp.Response.TotalCount)
	return
}

// GetInstanceList 获取实例列表
func (ten *TencentResource) GetInstanceList(pageSize, currentPage int) (count int, instanceList []*navite.Instance) {
	req := cvm.NewDescribeInstancesRequest()
	req.Limit, req.Offset = GetPageLimitInt64(pageSize, currentPage)
	resp, err := ten.cvm.DescribeInstances(req)
	if err != nil {
		log.Errorf("tencent describe instance failed: %v", err)
		return
	}
	for _, res := range resp.Response.InstanceSet {
		instance := &navite.Instance{
			CloudName:    constants.Tencent,
			AccountID:    ten.account.AccountID(),
			RegionID:     ten.account.RunRegionID,
			ZoneID:       *res.Placement.Zone,
			VPCID:        *res.VirtualPrivateCloud.VpcId,
			InstanceID:   *res.InstanceId,
			InstanceName: *res.InstanceName,
			InstanceType: *res.InstanceType,
			Status:       *res.InstanceState,
			CPU:          int(*res.CPU),
			Memory:       int(*res.Memory) * 1024,
			OSName:       *res.OsName,
			SecurityGroupList: func() []string {
				sList := []string{}
				for _, v := range res.SecurityGroupIds {
					sList = append(sList, *v)
				}
				return sList
			}(),
			EipAddress: func() string {
				eList := []string{}
				for _, v := range res.PublicIpAddresses {
					eList = append(eList, *v)
				}
				return strings.Join(eList, ",")
			}(),
			KeyPairList: func() []string {
				kList := []string{}
				for _, v := range res.LoginSettings.KeyIds {
					kList = append(kList, *v)
				}
				return kList
			}(),
			ImageID:     *res.ImageId,
			ChargeType:  *res.InstanceChargeType,
			CreatedTime: tool.TimeForISO8601(*res.CreatedTime),
			SyncedTime:  time.Now(),
		}
		instanceList = append(instanceList, instance)
	}
	count = int(*resp.Response.TotalCount)
	return
}

// GetSecurityGroupList 获取安全组列表
func (ten *TencentResource) GetSecurityGroupList(pageSize, currentPage int) (count int, sgList []*navite.SecurityGroup) {
	req := vpc.NewDescribeSecurityGroupsRequest()
	req.Limit, req.Offset = GetPageLimitString(pageSize, currentPage)
	resp, err := ten.vpc.DescribeSecurityGroups(req)
	if err != nil {
		log.Errorf("tencent describe securityGroup failed: %v", err)
		return
	}
	for _, res := range resp.Response.SecurityGroupSet {
		sg := &navite.SecurityGroup{
			CloudName:   constants.Tencent,
			AccountID:   ten.account.AccountID(),
			RegionID:    ten.account.RunRegionID,
			GroupID:     *res.SecurityGroupId,
			GroupName:   *res.SecurityGroupName,
			IsDefault:   *res.IsDefault,
			Description: *res.SecurityGroupDesc,
			CreatedTime: tool.TimeForISO8601(*res.CreatedTime),
			SyncedTime:  time.Now(),
		}
		sgList = append(sgList, sg)
	}
	count = int(*resp.Response.TotalCount)
	return
}

// GetDiskList 获取磁盘列表
func (ten *TencentResource) GetDiskList(pageSize, currentPage int) (count int, diskList []*navite.Disk) {
	req := cbs.NewDescribeDisksRequest()
	req.Limit, req.Offset = GetPageLimitUint64(pageSize, currentPage)
	resp, err := ten.cbs.DescribeDisks(req)
	if err != nil {
		log.Errorf("tencent describe disks failed: %v", err)
		return
	}
	for _, res := range resp.Response.DiskSet {
		disk := &navite.Disk{
			CloudName:        constants.Tencent,
			AccountID:        ten.account.AccountID(),
			RegionID:         ten.account.RunRegionID,
			DiskID:           *res.DiskId,
			DiskName:         *res.DiskName,
			DiskType:         *res.DiskType,
			ChargeType:       *res.DiskChargeType,
			IsEncrypted:      *res.Encrypt,
			Shareable:        *res.Shareable,
			DiskSize:         int(*res.DiskSize),
			Status:           *res.DiskState,
			AttachInstanceID: *res.InstanceId,
			CreatedTime:      tool.TimeForISO8601(*res.CreateTime),
			SyncedTime:       time.Now(),
		}
		diskList = append(diskList, disk)
	}
	count = int(*resp.Response.TotalCount)
	return
}

// GetKeypairList 获取密钥对
func (ten *TencentResource) GetKeypairList(pageSize, currentPage int) (count int, keypairList []*navite.Keypair) {
	req := cvm.NewDescribeKeyPairsRequest()
	req.Limit, req.Offset = GetPageLimitInt64(pageSize, currentPage)
	resp, err := ten.cvm.DescribeKeyPairs(req)
	if err != nil {
		log.Errorf("tencent describe keypairs failed: %v", err)
		return
	}
	for _, res := range resp.Response.KeyPairSet {
		keypair := &navite.Keypair{
			CloudName:   constants.Tencent,
			AccountID:   ten.account.AccountID(),
			RegionID:    ten.account.RunRegionID,
			KeypairID:   *res.KeyId,
			KeypairName: *res.KeyName,
			PublicKey:   *res.PublicKey,
			Description: *res.Description,
			CreatedTime: tool.TimeForISO8601(*res.CreatedTime),
			SyncedTime:  time.Now(),
		}
		keypairList = append(keypairList, keypair)
	}
	count = int(*resp.Response.TotalCount)
	return
}

// GetSecurityGroupRuleList 获取安全组规则
func (ten *TencentResource) GetSecurityGroupRuleList(securityGroupID string) (sgrList []*navite.SecurityGroupRule) {
	req := vpc.NewDescribeSecurityGroupPoliciesRequest()
	req.SecurityGroupId = &securityGroupID
	resp, err := ten.vpc.DescribeSecurityGroupPolicies(req)
	if err != nil {
		log.Errorf("tencnet describe securityGroupRules failed: %v", err)
		return
	}
	for _, res := range resp.Response.SecurityGroupPolicySet.Egress {
		sgr := &navite.SecurityGroupRule{
			CloudName:   constants.Tencent,
			GroupID:     securityGroupID,
			DestCidrIP:  *res.CidrBlock,
			Direction:   constants.FlowEgress,
			Protocol:    *res.Protocol,
			PortRange:   *res.Port,
			Action:      *res.Action,
			Description: *res.PolicyDescription,
			SyncedTime:  time.Now(),
		}
		sgrList = append(sgrList, sgr)
	}
	for _, res := range resp.Response.SecurityGroupPolicySet.Ingress {
		sgr := &navite.SecurityGroupRule{
			CloudName:    constants.Tencent,
			GroupID:      securityGroupID,
			SourceCidrIP: *res.CidrBlock,
			Direction:    constants.FlowIngress,
			Protocol:     *res.Protocol,
			PortRange:    *res.Port,
			Action:       *res.Action,
			Description:  *res.PolicyDescription,
			SyncedTime:   time.Now(),
		}
		sgrList = append(sgrList, sgr)
	}
	return sgrList
}

// GetInstanceSpecsList 获取实例规格
func (ten *TencentResource) GetInstanceSpecsList() (instantSpecList []*navite.InstanceSpec) {
	req := cvm.NewDescribeZoneInstanceConfigInfosRequest()
	resp, err := ten.cvm.DescribeZoneInstanceConfigInfos(req)
	if err != nil {
		log.Errorf("tencent describe instanceTypes failed: %v", err)
		return
	}
	for _, res := range resp.Response.InstanceTypeQuotaSet {
		spec := &navite.InstanceSpec{
			CloudName:        constants.Tencent,
			AccountID:        ten.account.AccountID(),
			RegionID:         ten.account.RunRegionID,
			ZoneID:           *res.Zone,
			InstanceSpecID:   *res.InstanceType,
			InstanceSpecName: *res.TypeName,
			InstanceFamily:   *res.InstanceFamily,
			CPU:              int(*res.Cpu),
			Memory:           float64(*res.Memory),
			Status:           *res.Status,
			SyncedTime:       time.Now(),
		}
		instantSpecList = append(instantSpecList, spec)
	}
	return
}

// GetVPCList 获取VPC资源列表
func (ten *TencentResource) GetVPCList(pageSize, currentPage int) (count int, vpcList []*navite.VPC) {
	req := vpc.NewDescribeVpcsRequest()
	req.Limit, req.Offset = GetPageLimitString(pageSize, currentPage)
	resp, err := ten.vpc.DescribeVpcs(req)
	if err != nil {
		log.Errorf("tencnet describe vpcs failed: %v", err)
		return
	}
	for _, res := range resp.Response.VpcSet {
		v := &navite.VPC{
			CloudName:   constants.Tencent,
			AccountID:   ten.account.AccountID(),
			RegionID:    ten.account.RunRegionID,
			VPCID:       *res.VpcId,
			VPCName:     *res.VpcName,
			IsDefault:   *res.IsDefault,
			CidrBlock:   *res.CidrBlock,
			CreatedTime: tool.TimeForISO8601(*res.CreatedTime),
			SyncedTime:  time.Now(),
		}
		vpcList = append(vpcList, v)
	}
	return int(*resp.Response.TotalCount), vpcList
}

// GetSubnetList 获取子网列表
func (ten *TencentResource) GetSubnetList(pageSize, currentPage int) (count int, subnetList []*navite.Subnet) {
	req := vpc.NewDescribeSubnetsRequest()
	req.Limit, req.Offset = GetPageLimitString(pageSize, currentPage)
	resp, err := ten.vpc.DescribeSubnets(req)
	if err != nil {
		log.Errorf("tencnet describe subnets failed: %v", err)
		return
	}
	for _, res := range resp.Response.SubnetSet {
		subnet := &navite.Subnet{
			CloudName:               constants.Tencent,
			AccountID:               ten.account.AccountID(),
			RegionID:                ten.account.RunRegionID,
			VPCID:                   *res.VpcId,
			ZoneID:                  *res.Zone,
			SubnetID:                *res.SubnetId,
			SubnetName:              *res.SubnetName,
			CidrBlock:               *res.CidrBlock,
			IsDefault:               *res.IsDefault,
			EnableBroadcast:         *res.EnableBroadcast,
			AvailableIPAddressCount: int(*res.AvailableIpAddressCount),
			IsVPCSnat:               *res.IsRemoteVpcSnat,
			CreatedTime:             tool.TimeForISO8601(*res.CreatedTime),
			SyncedTime:              time.Now(),
		}
		subnetList = append(subnetList, subnet)
	}
	return int(*resp.Response.TotalCount), subnetList
}

// GetEipList 获取弹性公网IP列表
func (ten *TencentResource) GetEipList(pageSize, currentPage int) (count int, eipList []*navite.Eip) {
	req := vpc.NewDescribeAddressesRequest()
	req.Limit, req.Offset = GetPageLimitInt64(pageSize, currentPage)
	resp, err := ten.vpc.DescribeAddresses(req)
	if err != nil {
		log.Errorf("tencent describe eips failed: %v", err)
		return
	}
	for _, res := range resp.Response.AddressSet {
		eip := &navite.Eip{
			CloudName:          constants.Tencent,
			AccountID:          ten.account.AccountID(),
			RegionID:           ten.account.RunRegionID,
			AddressID:          *res.AddressId,
			AddressName:        *res.AddressName,
			AddressStatus:      *res.AddressStatus,
			AddressIP:          *res.AddressIp,
			AddressType:        *res.AddressType,
			BindInstanceID:     *res.InstanceId,
			NetworkInterfaceID: *res.NetworkInterfaceId,
			CreatedTime:        tool.TimeForISO8601(*res.CreatedTime),
			SyncedTime:         time.Now(),
		}
		eipList = append(eipList, eip)
	}
	return int(*resp.Response.TotalCount), eipList
}

// NewKeypair 创建新的密钥对
func (ten *TencentResource) NewKeypair(keypair *navite.Keypair) (err error) {
	req := cvm.NewImportKeyPairRequest()
	var defaultProject int64
	req.KeyName = &keypair.KeypairName
	req.PublicKey = &keypair.PublicKey
	req.ProjectId = &defaultProject
	resp, err := ten.cvm.ImportKeyPair(req)
	if err != nil {
		log.Errorf("tencent import keypair failed: %v", err)
		return err
	}
	keypair.KeypairID = *resp.Response.KeyId
	return
}

// DeleteKeypair 删除密钥对
func (ten *TencentResource) DeleteKeypair(keypairIDList ...string) (err error) {
	req := cvm.NewDeleteKeyPairsRequest()
	var keyIds []*string
	for _, i := range keypairIDList {
		keyIds = append(keyIds, &i)
	}
	req.KeyIds = keyIds
	_, err = ten.cvm.DeleteKeyPairs(req)
	if err != nil {
		log.Errorf("tencent delete keypair [%s] failed: %v", req.ToJsonString(), err)
	}
	return err
}

// NewSecurityGroup 创建安全组
func (ten *TencentResource) NewSecurityGroup(sg *navite.SecurityGroup) (err error) {
	req := vpc.NewCreateSecurityGroupRequest()
	req.GroupName = &sg.GroupName
	req.GroupDescription = &sg.Description
	resp, err := ten.vpc.CreateSecurityGroup(req)
	if err != nil {
		log.Errorf("tencnet create securityGroup [%s] failed: %v", req.ToJsonString(), err)
		return err
	}
	sg.GroupID = *resp.Response.SecurityGroup.SecurityGroupId
	return nil
}

// DeleteSecurityGroup 删除安全组
func (ten *TencentResource) DeleteSecurityGroup(sgID string) (err error) {
	req := vpc.NewDeleteSecurityGroupRequest()
	req.SecurityGroupId = &sgID
	_, err = ten.vpc.DeleteSecurityGroup(req)
	if err != nil {
		log.Errorf("tencent create securityGroup [%s] failed: %v", req.ToJsonString(), err)
	}
	return
}

// NewSecurityGroupRule 创建安全组规则
//
// * Protocol、Action 要大写
func (ten *TencentResource) NewSecurityGroupRule(rule *navite.SecurityGroupRule) (err error) {
	var (
		egress  []*vpc.SecurityGroupPolicy
		ingress []*vpc.SecurityGroupPolicy
	)

	switch strings.ToLower(rule.Direction) {
	case constants.FlowEgress:
		egress = append(egress, &vpc.SecurityGroupPolicy{
			Port:              &rule.PortRange,
			Protocol:          &rule.Protocol,
			CidrBlock:         &rule.DestCidrIP,
			Action:            &rule.Action,
			PolicyDescription: &rule.Description,
		})
	case constants.FlowIngress:
		ingress = append(ingress, &vpc.SecurityGroupPolicy{
			Port:              &rule.PortRange,
			Protocol:          &rule.Protocol,
			CidrBlock:         &rule.SourceCidrIP,
			Action:            &rule.Action,
			PolicyDescription: &rule.Description,
		})
	}
	req := vpc.NewCreateSecurityGroupPoliciesRequest()
	req.SecurityGroupId = &rule.GroupID
	req.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
		Egress:  egress,
		Ingress: ingress,
	}
	_, err = ten.vpc.CreateSecurityGroupPolicies(req)
	if err != nil {
		log.Errorf("tencent create securityGroupRule [%s] failed: %v", req.ToJsonString(), err)
		return err
	}
	return nil
}

// DeleteSecurityGroupRule 删除安全组规则
func (ten *TencentResource) DeleteSecurityGroupRule(rule *navite.SecurityGroupRule) (err error) {
	req := vpc.NewDeleteSecurityGroupPoliciesRequest()
	req.SecurityGroupId = &rule.GroupID
	switch rule.Direction {
	case constants.FlowIngress:
		req.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
			Ingress: []*vpc.SecurityGroupPolicy{
				&vpc.SecurityGroupPolicy{
					Protocol:  &rule.Protocol,
					Port:      &rule.PortRange,
					CidrBlock: &rule.SourceCidrIP,
					Action:    &rule.Action,
				},
			},
		}
	case constants.FlowEgress:
		req.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
			Egress: []*vpc.SecurityGroupPolicy{
				&vpc.SecurityGroupPolicy{
					Protocol:  &rule.Protocol,
					Port:      &rule.PortRange,
					CidrBlock: &rule.DestCidrIP,
					Action:    &rule.Action,
				},
			},
		}
	}
	if _, err := ten.vpc.DeleteSecurityGroupPolicies(req); err != nil {
		log.Errorf("tencent delete securityGroupRule [%s] failed: %v", req.ToJsonString(), err)
		return err
	}
	return nil
}

// NewVPC 创建虚拟专用网络
func (ten *TencentResource) NewVPC(v *navite.VPC) (err error) {
	req := vpc.NewCreateVpcRequest()
	req.VpcName = &v.VPCName
	req.CidrBlock = &v.CidrBlock
	resp, err := ten.vpc.CreateVpc(req)
	if err != nil {
		log.Errorf("tencent create vpc [%s] failed: %v", req.ToJsonString(), err)
		return err
	}
	v.VPCID = *resp.Response.Vpc.VpcId
	v.IsDefault = *resp.Response.Vpc.IsDefault
	v.CreatedTime = tool.TimeForISO8601(*resp.Response.Vpc.CreatedTime)
	return
}

// DeleteVPC 删除虚拟专用网络
func (ten *TencentResource) DeleteVPC(vpcID string) (err error) {
	req := vpc.NewDeleteVpcRequest()
	req.VpcId = &vpcID
	_, err = ten.vpc.DeleteVpc(req)
	if err != nil {
		log.Errorf("tencent delete vpc [%s] failed: %v", req.ToJsonString(), err)
	}
	return
}

// NewSubnet 创建子网
func (ten *TencentResource) NewSubnet(subnet *navite.Subnet) (err error) {
	req := vpc.NewCreateSubnetRequest()
	req.VpcId = &subnet.VPCID
	req.SubnetName = &subnet.SubnetName
	req.CidrBlock = &subnet.CidrBlock
	req.Zone = &subnet.ZoneID
	resp, err := ten.vpc.CreateSubnet(req)
	if err != nil {
		log.Errorf("tencent create subnet [%s] failed: %v", req.ToJsonString(), err)
		return err
	}
	subnet.SubnetID = *resp.Response.Subnet.SubnetId
	subnet.EnableBroadcast = *resp.Response.Subnet.EnableBroadcast
	subnet.IsDefault = *resp.Response.Subnet.IsDefault
	subnet.IsVPCSnat = *resp.Response.Subnet.IsRemoteVpcSnat
	subnet.AvailableIPAddressCount = int(*resp.Response.Subnet.AvailableIpAddressCount)
	return
}

// DeleteSubnet 删除子网
func (ten *TencentResource) DeleteSubnet(subnetID string) (err error) {
	req := vpc.NewDeleteSubnetRequest()
	req.SubnetId = &subnetID
	_, err = ten.vpc.DeleteSubnet(req)
	if err != nil {
		log.Errorf("tencent delete subnet [%s] failed: %v", req.ToJsonString(), err)
	}
	return
}

// NewDisk 创建云盘
func (ten *TencentResource) NewDisk(disk *navite.Disk) (err error) {
	req := cbs.NewCreateDisksRequest()
	chargeType := strings.ToUpper(disk.ChargeType)
	diskType := strings.ToUpper(disk.DiskType)
	req.DiskType = &diskType
	req.DiskChargeType = &chargeType
	req.Placement = &cbs.Placement{
		Zone: &disk.ZoneID,
	}
	if disk.DiskSize < 50 {
		disk.DiskSize = 50
	}
	size := uint64(disk.DiskSize)
	req.DiskName = &disk.DiskName
	req.DiskSize = &size
	if disk.IsEncrypted {
		encrypted := "ENCRYPT"
		req.Encrypt = &encrypted
	}
	req.Shareable = &disk.Shareable
	resp, err := ten.cbs.CreateDisks(req)
	if err != nil {
		log.Errorf("tencent create subnet [%s] failed: %v", req.ToJsonString(), err)
		return err
	}
	disk.DiskID = *resp.Response.DiskIdSet[0]
	return
}

// DeleteDisk 删除云盘
func (ten *TencentResource) DeleteDisk(diskIDList ...string) (err error) {
	req := cbs.NewTerminateDisksRequest()
	var diskIds []*string
	for _, i := range diskIDList {
		diskIds = append(diskIds, &i)
	}
	req.DiskIds = diskIds
	_, err = ten.cbs.TerminateDisks(req)
	if err != nil {
		log.Errorf("tencent delete disk [%s] failed: %v", req.ToJsonString(), err)
	}
	return err
}

// NewEIP 申请弹性公网IP
func (ten *TencentResource) NewEIP(eip *navite.Eip) (err error) {
	numbers := int64(1)
	req := vpc.NewAllocateAddressesRequest()
	req.AddressCount = &numbers
	resp, err := ten.vpc.AllocateAddresses(req)
	if err != nil {
		log.Errorf("tencent create eip [%s] failed: %v", req.ToJsonString(), err)
		return err
	}
	eip.AddressID = *resp.Response.AddressSet[0]
	return
}

// ReleaseEIP 释放弹性公网IP
func (ten *TencentResource) ReleaseEIP(eipIDList ...string) (err error) {
	req := vpc.NewReleaseAddressesRequest()
	var eipIds []*string
	for _, i := range eipIDList {
		eipIds = append(eipIds, &i)
	}
	req.AddressIds = eipIds
	_, err = ten.vpc.ReleaseAddresses(req)
	if err != nil {
		log.Errorf("tencent release eip [%s] failed: %v", req.ToJsonString(), err)
	}
	return err
}

// RunInstance 创建实例
func (ten *TencentResource) RunInstance(instance *param.RunInstanceParam) (instanceIDList []string, err error) {
	req := cvm.NewRunInstancesRequest()
	// 1. 位置区域
	req.Placement = &cvm.Placement{
		Zone: &instance.ZoneID,
	}
	// 2. 镜像, 机型(配置)
	req.ImageId = &instance.ImageID           // 镜像，系统
	req.InstanceType = &instance.InstanceType // 机型，内存/CPU
	// 3. 登陆sshkey
	req.LoginSettings = &cvm.LoginSettings{
		KeyIds: []*string{&instance.KeyPairID},
	}

	// 4. 安全组, 网络
	req.SecurityGroupIds = []*string{&instance.SecurityGroupID}
	req.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		VpcId:    &instance.VPCID,
		SubnetId: &instance.SubnetID,
	}

	// 5. 磁盘设置
	dsize := int64(instance.DiskSize)
	diskList := []*cvm.DataDisk{
		&cvm.DataDisk{
			DiskSize: &dsize,
			DiskType: &instance.DiskType,
		},
	}
	req.DataDisks = diskList

	// 6. 自定义设置
	instanceCount := int64(instance.Numbers)
	req.InstanceName = &instance.InstanceName
	req.HostName = &instance.HostName
	req.InstanceCount = &instanceCount

	resp, err := ten.cvm.RunInstances(req)
	if err != nil {
		log.Errorf("tencent runInstance [%s] failed: %v", req.ToJsonString(), err)
		return nil, err
	}
	for _, v := range resp.Response.InstanceIdSet {
		instanceIDList = append(instanceIDList, *v)
	}
	return instanceIDList, nil
}

// DeleteInstance 删除实例
func (ten *TencentResource) DeleteInstance(instanceIDList ...string) (err error) {
	if len(instanceIDList) > 100 {
		return fmt.Errorf("Exceeding the single operation limit")
	}
	req := cvm.NewTerminateInstancesRequest()
	instanceIds := []*string{}
	for _, v := range instanceIDList {
		instanceIds = append(instanceIds, &v)
	}
	req.InstanceIds = instanceIds
	_, err = ten.cvm.TerminateInstances(req)
	if err != nil {
		log.Errorf("tencent terminateInstance [%s] failed: %v", req.ToJsonString(), err)
	}
	return err
}

// StartInstance 启动实例
//
// * 只有状态为STOPPED的实例才可以进行此操作
func (ten *TencentResource) StartInstance(instanceIDList ...string) (err error) {
	req := cvm.NewStartInstancesRequest()
	instanceIds := []*string{}
	for _, v := range instanceIDList {
		instanceIds = append(instanceIds, &v)
	}
	req.InstanceIds = instanceIds
	_, err = ten.cvm.StartInstances(req)
	if err != nil {
		log.Errorf("tencent startInstance [%s] failed: %v", req.ToJsonString(), err)
	}
	return err
}

// StopInstance 停止实例
//
// * 只有状态为RUNNING的实例才可以进行此操作
func (ten *TencentResource) StopInstance(instanceIDList ...string) (err error) {
	req := cvm.NewStopInstancesRequest()
	instanceIds := []*string{}
	for _, v := range instanceIDList {
		instanceIds = append(instanceIds, &v)
	}
	req.InstanceIds = instanceIds
	_, err = ten.cvm.StopInstances(req)
	if err != nil {
		log.Errorf("tencent stopInstance [%s] failed: %v", req.ToJsonString(), err)
	}
	return err
}

// RebotInstance 停止实例
//
// * 只有状态为RUNNING的实例才可以进行此操作
func (ten *TencentResource) RebotInstance(instanceIDList ...string) (err error) {
	req := cvm.NewRebootInstancesRequest()
	instanceIds := []*string{}
	for _, v := range instanceIDList {
		instanceIds = append(instanceIds, &v)
	}
	req.InstanceIds = instanceIds
	_, err = ten.cvm.RebootInstances(req)
	if err != nil {
		log.Errorf("tencent rebotInstance [%s] failed: %v", req.ToJsonString(), err)
	}
	return err
}

// AttachDisk 挂载磁盘到实例上
func (ten *TencentResource) AttachDisk(instance *navite.Instance, disk *navite.Disk) (err error) {
	req := cbs.NewAttachDisksRequest()
	req.InstanceId = &instance.InstanceID
	req.DiskIds = []*string{&disk.DiskID}
	_, err = ten.cbs.AttachDisks(req)
	if err != nil {
		log.Errorf("tencent attachDisk [%s] failed: %v", req.ToJsonString(), err)
	}
	return err
}

// DetachDisk 解挂载磁盘
func (ten *TencentResource) DetachDisk(instance *navite.Instance, disk *navite.Disk) (err error) {
	req := cbs.NewDetachDisksRequest()
	req.DiskIds = []*string{&disk.DiskID}
	_, err = ten.cbs.DetachDisks(req)
	if err != nil {
		log.Errorf("tencent detachDisk [%s] failed: %v", req.ToJsonString(), err)
	}
	return err
}

// AttachEipToInstance 绑定弹性公网IP到实例上
func (ten *TencentResource) AttachEipToInstance(instance *navite.Instance, eip *navite.Eip) (err error) {
	req := vpc.NewAssociateAddressRequest()
	req.AddressId = &eip.AddressID
	req.InstanceId = &instance.InstanceID
	_, err = ten.vpc.AssociateAddress(req)
	if err != nil {
		log.Errorf("tencent attachEipToInstance [%s] failed: %v", req.ToJsonString(), err)
	}
	return err

}

// DetachEipFromInstance 从实例上解绑弹性公网IP
func (ten *TencentResource) DetachEipFromInstance(instance *navite.Instance, eip *navite.Eip) (err error) {
	req := vpc.NewDisassociateAddressRequest()
	req.AddressId = &eip.AddressID
	_, err = ten.vpc.DisassociateAddress(req)
	if err != nil {
		log.Errorf("tencent detachEipFromInstance [%s] failed: %v", req.ToJsonString(), err)
	}
	return err
}

// ModifyEIPBandWidth 调整弹性公网IP的带宽
func (ten *TencentResource) ModifyEIPBandWidth(eip *navite.Eip, bandWidth int64) (err error) {
	req := vpc.NewModifyAddressesBandwidthRequest()
	req.AddressIds = []*string{&eip.AddressID}
	req.InternetMaxBandwidthOut = &bandWidth
	_, err = ten.vpc.ModifyAddressesBandwidth(req)
	if err != nil {
		log.Errorf("tencent modifyEIPBandWidth [%s] failed: %v", req.ToJsonString(), err)
		return err
	}
	eip.BandWidth = bandWidth
	return
}
