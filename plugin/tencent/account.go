package tencent

import (
	"ark-common/clients/mgo"
	"ark-common/constants"
	"ark-common/resource/navite"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TencentAccount 腾讯云账户
type TencentAccount struct {
	rbd *mgo.Client
}

// NewTencentAccountPlugin 初始化阿里云驱动
func NewTencentAccountPlugin(rbd *mgo.Client) *TencentAccount {
	return &TencentAccount{
		rbd: rbd,
	}
}

// BindAccount 绑定云账号
func (ali *TencentAccount) BindAccount(accountName, ak, sk string) *navite.CloudAccount {
	account := &navite.CloudAccount{
		ID:          primitive.NewObjectID(),
		AccountName: accountName,
		CloudName:   constants.Tencent,
		AccessKey:   ak,
		SecurityKey: sk,
		CreatedTime: time.Now(),
	}
	account.Encryption()
	_, err := ali.rbd.Table(navite.CloudAccountTable).Insert(account)
	if err != nil {
		log.Errorf("bind tencent account [%+v] failed: %v", account, err)
	}
	return account
}
