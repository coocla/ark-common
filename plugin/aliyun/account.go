package aliyun

import (
	"ark-common/clients/mgo"
	"ark-common/constants"
	"ark-common/resource/navite"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AliyunAccount 阿里云账户
type AliyunAccount struct {
	rbd *mgo.Client
}

// NewAliyunAccountPlugin 初始化阿里云驱动
func NewAliyunAccountPlugin(rbd *mgo.Client) *AliyunAccount {
	return &AliyunAccount{
		rbd: rbd,
	}
}

// BindAccount 绑定云账号
func (ali *AliyunAccount) BindAccount(accountName, ak, sk string) *navite.CloudAccount {
	account := &navite.CloudAccount{
		ID:          primitive.NewObjectID(),
		AccountName: accountName,
		CloudName:   constants.Aliyun,
		AccessKey:   ak,
		SecurityKey: sk,
		CreatedTime: time.Now(),
	}
	account.Encryption()
	_, err := ali.rbd.Table(navite.CloudAccountTable).Insert(account)
	if err != nil {
		log.Errorf("bind aliyun account [%+v] failed: %v", account, err)
	}
	return account
}
