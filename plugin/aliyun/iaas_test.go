package aliyun_test

import (
	"ark-common/clients/mgo"
	"ark-common/constants"
	"ark-common/param"
	"ark-common/plugin/aliyun"
	"ark-common/resource/manage"
	"ark-common/resource/navite"
	"ark-common/utils/tool"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	db      *mgo.Client
	account *navite.CloudAccount
	driver  *aliyun.AliyunResource
)

func init() {
	db = mgo.NewMgo("")
	_, accountList := manage.ListCloudAccount(db, constants.Aliyun, "", 1, 1)
	account = accountList[0]
	account.RunRegionID = "cn-beijing"
	driver = aliyun.NewAliyunPlugin(account)
}

func TestKeyPair(t *testing.T) {
	var (
		publicKey []byte
		keypair   *navite.Keypair
	)

	Convey("测试 aliyun 密钥对", t, func() {
		Convey("生成RSA密钥对", func() {
			_, publicKey = tool.NewRSAKeyPair()
			So(publicKey, ShouldNotEqual, "")
		})
		Convey("将密钥对导入到云端", func() {
			keypair = &navite.Keypair{
				CloudName:   constants.Aliyun,
				AccountID:   account.AccountID(),
				RegionID:    account.RunRegionID,
				KeypairName: "Test0001",
				PublicKey:   string(publicKey),
				CreatedTime: time.Now(),
			}
			err := driver.NewKeypair(keypair)
			So(err, ShouldBeNil)
		})
		Convey("删除密钥对", func() {
			err := driver.DeleteKeypair(keypair.KeypairID)
			So(err, ShouldBeNil)
		})
	})
}

func TestDisk(t *testing.T) {
	var disk *navite.Disk
	Convey("测试 aliyun 云盘", t, func() {
		Convey("创建云盘", func() {
			disk = &navite.Disk{
				CloudName: constants.Aliyun,
				AccountID: account.AccountID(),
				RegionID:  account.RunRegionID,
				ZoneID:    "cn-beijing-g",
				DiskName:  "Test00001",
				DiskType:  "cloud_efficiency",
				DiskSize:  20,
			}
			err := driver.NewDisk(disk)
			So(err, ShouldBeNil)
			time.Sleep(20)
		})
		Convey("删除云盘", func() {
			err := driver.DeleteDisk(disk.DiskID)
			So(err, ShouldBeNil)
		})
	})
}

func TestSecurityGroup(t *testing.T) {
	var sg *navite.SecurityGroup
	Convey("测试安全组", t, func() {
		Convey("创建安全组", func() {
			sg = &navite.SecurityGroup{
				CloudName: constants.Aliyun,
				AccountID: account.AccountID(),
				RegionID:  account.RunRegionID,
				GroupName: "TestSG",
			}
			err := driver.NewSecurityGroup(sg)
			So(err, ShouldBeNil)
		})
		Convey("删除安全组", func() {
			err := driver.DeleteSecurityGroup(sg.GroupID)
			So(err, ShouldBeNil)
		})
	})
}

func TestSecurityGroupRule(t *testing.T) {
	ingressSgr := &navite.SecurityGroupRule{
		CloudName:    account.CloudName,
		GroupID:      "sg-2zefpzvxfxkw5nkflw3k",
		SourceCidrIP: "10.10.1.0/24",
		Direction:    constants.FlowIngress,
		PortRange:    "8081/8081",
		Protocol:     "tcp",
		Action:       "drop",
	}

	egressSgr := &navite.SecurityGroupRule{
		CloudName:  account.CloudName,
		GroupID:    "sg-2zefpzvxfxkw5nkflw3k",
		DestCidrIP: "10.10.1.0/24",
		Direction:  constants.FlowEgress,
		PortRange:  "8081/8081",
		Protocol:   "tcp",
		Action:     "drop",
	}
	Convey("测试创建安全组规则", t, func() {
		Convey("创建入站规则", func() {
			err := driver.NewSecurityGroupRule(ingressSgr)
			So(err, ShouldBeNil)
		})
		Convey("创建出站规则", func() {
			err := driver.NewSecurityGroupRule(egressSgr)
			So(err, ShouldBeNil)
		})
	})
}

func TestDeleteSecurityGroupRule(t *testing.T) {
	ingressSgr := &navite.SecurityGroupRule{
		CloudName:    account.CloudName,
		GroupID:      "sg-2zefpzvxfxkw5nkflw3k",
		SourceCidrIP: "10.10.1.0/24",
		Direction:    constants.FlowIngress,
		PortRange:    "8081/8081",
		Protocol:     "tcp",
		Action:       "drop",
	}

	egressSgr := &navite.SecurityGroupRule{
		CloudName:  account.CloudName,
		GroupID:    "sg-2zefpzvxfxkw5nkflw3k",
		DestCidrIP: "10.10.1.0/24",
		Direction:  constants.FlowEgress,
		PortRange:  "8081/8081",
		Protocol:   "tcp",
		Action:     "drop",
	}
	Convey("测试删除安全组规则", t, func() {
		Convey("删除入站规则", func() {
			err := driver.DeleteSecurityGroupRule(ingressSgr)
			So(err, ShouldBeNil)
		})
		Convey("删除出站规则", func() {
			err := driver.DeleteSecurityGroupRule(egressSgr)
			So(err, ShouldBeNil)
		})
	})
}
func TestVPC(t *testing.T) {
	var (
		vpc    *navite.VPC
		subnet *navite.Subnet
	)
	Convey("测试vpc", t, func() {
		SkipConvey("创建vpc", func() {
			vpc = &navite.VPC{
				CloudName: constants.Aliyun,
				AccountID: account.AccountID(),
				RegionID:  account.RunRegionID,
				VPCName:   "TestVPC",
				CidrBlock: "10.16.0.0/12",
			}
			err := driver.NewVPC(vpc)
			So(err, ShouldBeNil)
		})
		SkipConvey("创建子网", func() {
			subnet = &navite.Subnet{
				CloudName:  constants.Aliyun,
				AccountID:  account.AccountID(),
				RegionID:   account.RunRegionID,
				SubnetName: "TestSubnet",
				ZoneID:     "cn-beijing-g",
				VPCID:      "vpc-2zejndwk8qglmkw0phcix",
				CidrBlock:  "10.16.1.0/24",
			}
			err := driver.NewSubnet(subnet)
			So(err, ShouldBeNil)
		})

		Convey("删除子网", func() {
			err := driver.DeleteSubnet("vsw-2zeoz31yyneb80xhkz6u1")
			So(err, ShouldBeNil)
			Convey("删除vpc", func() {
				err := driver.DeleteVPC("vpc-2zejndwk8qglmkw0phcix")
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestRunInstance(t *testing.T) {
	p := &param.RunInstanceParam{
		AccountID:       account.AccountID(),
		RegionID:        account.RunRegionID,
		ZoneID:          "cn-beijing-g",
		ImageID:         "centos_7_03_64_20G_alibase_20170818.vhd",
		InstanceType:    "ecs.g5.large",
		HostName:        "TestArk",
		InstanceName:    "TestArk",
		KeyPairID:       "keypair-01",
		SecurityGroupID: "sg-2ze46vybnzfl08yrhhs2",
		SubnetID:        "vsw-2zeuuokhe86pe0rqe2zjt",
		Numbers:         1,
	}
	_, err := driver.RunInstance(p)
	Convey("创建虚拟机", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestDeleteInstance(t *testing.T) {
	instanceIDList := []string{
		"i-2ze5wlll5cslbfwnia8z",
	}
	err := driver.DeleteInstance(instanceIDList...)
	Convey("删除虚拟机", t, func() {
		So(err, ShouldBeNil)
	})
}
