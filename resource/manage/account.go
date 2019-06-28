package manage

import (
	"ark-common/clients/mgo"
	"ark-common/constants"
	"ark-common/resource/navite"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/sirupsen/logrus"
)

// GetAvailableCloudAccount 获取云商账户
func GetAvailableCloudAccount(rdb *mgo.Client) (count int, accountList []*navite.CloudAccount) {
	filter := bson.M{
		"disabled": false,
		"healthy":  true,
	}
	accountList = []*navite.CloudAccount{}
	total, err := rdb.Table(navite.CloudAccountTable).Count(filter, nil)
	if err != nil {
		log.Warnf("filter [%v] available account failed: %v", filter, err)
		return 0, accountList
	}
	mctx := context.Background()
	cur, err := rdb.Table(navite.CloudAccountTable).Query(filter, 0, 0, nil)
	if err != nil {
		log.Warnf("filter [%v] available account failed: %v", filter, err)
		return 0, accountList
	}
	defer cur.Close(mctx)
	for cur.Next(mctx) {
		v := &navite.CloudAccount{}
		e := cur.Decode(v)
		if e != nil {
			continue
		}
		accountList = append(accountList, v)
	}
	return int(total), accountList
}

// GetCloudAccount 获取对应云商账号
func GetCloudAccount(rdb *mgo.Client, accountID string) (account *navite.CloudAccount, errCode int) {
	objectID, err := primitive.ObjectIDFromHex(accountID)
	if err != nil {
		return nil, constants.InvalidCloudAccountID
	}
	filter := bson.M{
		"_id": objectID,
	}
	account = &navite.CloudAccount{}
	err = rdb.Table(navite.CloudAccountTable).QueryOne(filter, account, nil)
	if err != nil {
		log.Errorf("filter account failed: %v", err)
		errCode = constants.InvalidCloudAccountID
	}
	return
}

// ListCloudAccount 列出符合条件的云账号
func ListCloudAccount(rdb *mgo.Client, cloudName, accountName string, pageSize, currentPage int) (count int, accountList []*navite.CloudAccount) {
	filter := bson.M{}
	if cloudName != "" {
		filter["cloudName"] = cloudName
	}
	if accountName != "" {
		filter["accountName"] = accountName
	}
	accountList = []*navite.CloudAccount{}
	total, err := rdb.Table(navite.CloudAccountTable).Count(filter, nil)
	if err != nil {
		log.Warnf("list [%v] account failed: %v", filter, err)
		return 0, accountList
	}
	mctx := context.Background()
	cur, err := rdb.Table(navite.CloudAccountTable).Query(filter, pageSize, currentPage, nil)
	if err != nil {
		log.Warnf("list [%v] account failed: %v", filter, err)
		return 0, accountList
	}
	defer cur.Close(mctx)
	err = cur.All(mctx, &accountList)
	if err != nil {
		log.Errorf("decord mgo document failed: %v", err)
	}
	return int(total), accountList
}

// DeleteCloudAccount 删除云账号
func DeleteCloudAccount(rdb *mgo.Client, accountID string) error {
	objID, err := primitive.ObjectIDFromHex(accountID)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": objID,
	}
	_, err = rdb.Table(navite.CloudAccountTable).DeleteMany(filter)
	return err
}

// GetRegions 获取云商对应的区域
func GetRegions(rdb *mgo.Client, c *navite.CloudAccount) []*navite.CloudRegion {
	filter := bson.M{
		"cloudName": c.CloudName,
	}
	rs := []*navite.CloudRegion{}
	cur, err := rdb.Table(navite.CloudRegionTable).Query(filter, 0, 0, nil)
	if err != nil {
		log.Errorf("filter [%v] account region failed: %v", filter, err)
	}
	mctx := context.Background()
	defer cur.Close(mctx)
	for cur.Next(mctx) {
		v := &navite.CloudRegion{}
		e := cur.Decode(v)
		if e != nil {
			continue
		}
		rs = append(rs, v)
	}
	return rs
}
