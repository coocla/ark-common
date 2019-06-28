package tencent_test

import (
	"ark-common/clients/mgo"
	"ark-common/constants"
	"ark-common/param"
	"ark-common/plugin/tencent"
	"ark-common/resource/manage"
	"ark-common/resource/navite"
	"ark-common/utils/tool"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	db      *mgo.Client
	account *navite.CloudAccount
	driver  *tencent.TencentResource
)

func init() {
	mgoURN := os.Getenv("MGO_URL")
	if db == nil {
		db = mgo.NewMgo(mgoURN)
	}
	_, accountList := manage.ListCloudAccount(db, constants.Tencent, "", 1, 1)
	account = accountList[0]
	account.RunRegionID = "ap-beijing"
	driver = tencent.NewTencentPlugin(account)
}

func TestKeyPair(t *testing.T) {
	var (
		publicKey []byte
		keypair   *navite.Keypair
	)

	Convey("测试 tencent 密钥对", t, func() {
		Convey("生成RSA密钥对", func() {
			_, publicKey = tool.NewRSAKeyPair()
			So(publicKey, ShouldNotEqual, "")
		})
		Convey("将密钥对导入到云端", func() {
			keypair = &navite.Keypair{
				CloudName:   constants.Tencent,
				AccountID:   account.AccountID(),
				RegionID:    account.RunRegionID,
				KeypairName: "test",
				PublicKey:   string(publicKey),
				CreatedTime: time.Now(),
			}
			err := driver.NewKeypair(keypair)
			So(err, ShouldBeNil)
		})
		time.Sleep(5)
		Convey("删除密钥对", func() {
			err := driver.DeleteKeypair(keypair.KeypairID)
			So(err, ShouldBeNil)
		})
	})
}

func TestDisk(t *testing.T) {
	var (
		disk *navite.Disk
	)
	Convey("测试 tencent 云盘", t, func() {
		Convey("创建云盘", func() {
			disk = &navite.Disk{
				CloudName:  constants.Tencent,
				AccountID:  account.AccountID(),
				RegionID:   account.RunRegionID,
				DiskName:   "test",
				DiskSize:   50,
				DiskType:   "cloud_premium",
				ChargeType: "postpaid_by_hour",
				ZoneID:     "ap-beijing-1",
			}
			err := driver.NewDisk(disk)
			So(err, ShouldBeNil)
		})
		time.Sleep(20)
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
				CloudName:   constants.Tencent,
				AccountID:   account.AccountID(),
				RegionID:    account.RunRegionID,
				GroupName:   "TestSG",
				Description: "test",
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
		GroupID:      "sg-mhwek5rn",
		SourceCidrIP: "10.10.0.1/24",
		Direction:    constants.FlowIngress,
		Protocol:     "TCP",
		PortRange:    "8801",
		Action:       "DROP",
	}

	egressSgr := &navite.SecurityGroupRule{
		CloudName:  account.CloudName,
		GroupID:    "sg-mhwek5rn",
		DestCidrIP: "10.11.0.1/24",
		Direction:  constants.FlowEgress,
		Protocol:   "TCP",
		PortRange:  "8801",
		Action:     "DROP",
	}

	Convey("创建安全组规则", t, func() {
		Convey("创建入站规则", func() {
			ingressErr := driver.NewSecurityGroupRule(ingressSgr)
			So(ingressErr, ShouldBeNil)
		})
		Convey("创建出站规则", func() {
			egressErr := driver.NewSecurityGroupRule(egressSgr)
			So(egressErr, ShouldBeNil)
		})
	})
}

func TestDeleteSecurityGroupRule(t *testing.T) {
	ingressSgr := &navite.SecurityGroupRule{
		CloudName:    account.CloudName,
		GroupID:      "sg-mhwek5rn",
		SourceCidrIP: "10.10.0.1/24",
		Direction:    constants.FlowIngress,
		Protocol:     "TCP",
		PortRange:    "8801",
		Action:       "DROP",
	}
	egressSgr := &navite.SecurityGroupRule{
		CloudName:  account.CloudName,
		GroupID:    "sg-mhwek5rn",
		DestCidrIP: "10.11.0.1/24",
		Direction:  constants.FlowEgress,
		Protocol:   "TCP",
		PortRange:  "8801",
		Action:     "DROP",
	}
	Convey("删除安全组规则", t, func() {
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
	Convey("测试私有网络", t, func() {
		Convey("创建vpc", func() {
			vpc = &navite.VPC{
				CloudName: constants.Tencent,
				AccountID: account.AccountID(),
				RegionID:  account.RunRegionID,
				VPCName:   "TestVPC",
				CidrBlock: "192.168.0.0/16",
			}
			err := driver.NewVPC(vpc)
			So(err, ShouldBeNil)
		})
		Convey("创建子网", func() {
			subnet = &navite.Subnet{
				CloudName:  constants.Tencent,
				AccountID:  account.AccountID(),
				RegionID:   account.RunRegionID,
				SubnetName: "TestSubnet",
				VPCID:      vpc.VPCID,
				CidrBlock:  "192.168.1.0/24",
				ZoneID:     "ap-beijing-1",
			}
			err := driver.NewSubnet(subnet)
			So(err, ShouldBeNil)
			time.Sleep(10)
		})
		Convey("删除子网", func() {
			err := driver.DeleteSubnet(subnet.SubnetID)
			So(err, ShouldBeNil)
		})
		Convey("删除vpc", func() {
			err := driver.DeleteVPC(vpc.VPCID)
			So(err, ShouldBeNil)
		})
	})
}

func TestRunInstance(t *testing.T) {
	p := &param.RunInstanceParam{
		AccountID:    account.AccountID(),
		RegionID:     account.RunRegionID,
		ZoneID:       "ap-beijing-1",
		ImageID:      "img-9qabwvbn",
		InstanceType: "S1.4XLARGE16",
		HostName:     "TestArk",
		InstanceName: "TestArk",
		KeyPairID:    "skey-ksckh9ff",
		Numbers:      1,
	}
	_, err := driver.RunInstance(p)
	Convey("创建虚拟机", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestDeleteInstance(t *testing.T) {
	instanceIDList := []string{
		"ins-e9awl403",
	}
	err := driver.DeleteInstance(instanceIDList...)
	Convey("删除实例", t, func() {
		So(err, ShouldBeNil)
	})
}
