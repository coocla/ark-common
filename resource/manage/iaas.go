package manage

import (
	"ark-common/clients/mgo"
	"ark-common/resource/navite"
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// ListRegions 搜索地域
func ListRegions(rbd *mgo.Client, cloudName string, pageSize, currentPage int) (count int, regionList []*navite.CloudRegion) {
	filter := bson.M{}
	if cloudName != "" {
		filter["cloudName"] = cloudName
	}
	regionList = []*navite.CloudRegion{}
	total, err := rbd.Table(navite.CloudRegionTable).Count(filter, nil)
	if err != nil {
		log.Warnf("list [%v] regions failed: %v", filter, err)
		return 0, regionList
	}
	mctx := context.Background()
	cur, err := rbd.Table(navite.CloudRegionTable).Query(filter, pageSize, currentPage, nil)
	if err != nil {
		log.Warnf("list [%v] regions failed: %v", filter, err)
		return 0, regionList
	}
	defer cur.Close(mctx)
	err = cur.All(mctx, &regionList)
	if err != nil {
		log.Errorf("decord mgo document failed: %v", err)
	}
	return int(total), regionList
}

// ListImages 搜索镜像
func ListImages(rbd *mgo.Client, regionID, osType string, pageSize, currentPage int) (count int, imageList []*navite.Image) {
	filter := bson.M{}
	if regionID != "" {
		filter["regionId"] = regionID
	}
	if osType != "" {
		filter["osType"] = osType
	}
	imageList = []*navite.Image{}
	total, err := rbd.Table(navite.ImageTable).Count(filter, nil)
	if err != nil {
		log.Warnf("list [%v] images failed: %v", filter, err)
		return 0, imageList
	}
	mctx := context.Background()
	cur, err := rbd.Table(navite.ImageTable).Query(filter, pageSize, currentPage, nil)
	if err != nil {
		log.Warnf("list [%v] images failed: %v", filter, err)
		return 0, imageList
	}
	defer cur.Close(mctx)
	err = cur.All(mctx, &imageList)
	if err != nil {
		log.Errorf("decord mgo document failed: %v", err)
	}
	return int(total), imageList
}

// ListInstances 列出实例
func ListInstances(rbd *mgo.Client, pageSize, currentPage int) (count int, instanceList []*navite.Instance) {
	filter := bson.M{}
	instanceList = []*navite.Instance{}
	total, err := rbd.Table(navite.InstanceTable).Count(filter, nil)
	if err != nil {
		log.Warnf("list [%v] instances failed: %v", filter, err)
		return 0, instanceList
	}
	mctx := context.Background()
	cur, err := rbd.Table(navite.InstanceTable).Query(filter, pageSize, currentPage, nil)
	if err != nil {
		log.Warnf("list [%v] instances failed: %v", filter, err)
		return 0, instanceList
	}
	defer cur.Close(mctx)
	err = cur.All(mctx, &instanceList)
	if err != nil {
		log.Errorf("decord mgo document failed: %v", err)
	}
	return int(total), instanceList

}

// ListSecurityGroups 安全组列表
func ListSecurityGroups(rbd *mgo.Client, cloudName, accountID, regionID string, pageSize, currentPage int) (count int, sgList []*navite.SecurityGroup) {
	filter := bson.M{}
	if cloudName != "" {
		filter["cloudName"] = cloudName
	}
	if accountID != "" {
		filter["accountId"] = accountID
	}
	if regionID != "" {
		filter["regionId"] = regionID
	}
	sgList = []*navite.SecurityGroup{}
	total, err := rbd.Table(navite.SecurityGroupTable).Count(filter, nil)
	if err != nil {
		log.Warnf("list [%v] securityGroups failed: %v", filter, err)
		return 0, sgList
	}
	mctx := context.Background()
	cur, err := rbd.Table(navite.SecurityGroupTable).Query(filter, pageSize, currentPage, nil)
	if err != nil {
		log.Warnf("list [%v] securityGroups failed: %v", filter, err)
		return 0, sgList
	}
	defer cur.Close(mctx)
	err = cur.All(mctx, &sgList)
	if err != nil {
		log.Errorf("decord mgo document failed: %v", err)
	}
	return int(total), sgList
}

func ListDisks(rbd *mgo.Client, pageSize, currentPage int) (count int, diskList []*navite.Disk) {
	filter := bson.M{}
	diskList = []*navite.Disk{}
	total, err := rbd.Table(navite.DiskTable).Count(filter, nil)
	if err != nil {
		log.Warnf("list [%v] disks failed: %v", filter, err)
		return 0, diskList
	}
	mctx := context.Background()
	cur, err := rbd.Table(navite.DiskTable).Query(filter, pageSize, currentPage, nil)
	if err != nil {
		log.Warnf("list [%v] disks failed: %v", filter, err)
		return 0, diskList
	}
	defer cur.Close(mctx)
	err = cur.All(mctx, &diskList)
	if err != nil {
		log.Errorf("decord mgo document failed: %v", err)
	}
	return int(total), diskList
}

func ListKeypairs(rbd *mgo.Client, pageSize, currentPage int) (count int, keypairList []*navite.Keypair) {
	filter := bson.M{}
	keypairList = []*navite.Keypair{}
	total, err := rbd.Table(navite.KeyPairTable).Count(filter, nil)
	if err != nil {
		log.Warnf("list [%v] keypairs failed: %v", filter, err)
		return 0, keypairList
	}
	mctx := context.Background()
	cur, err := rbd.Table(navite.KeyPairTable).Query(filter, pageSize, currentPage, nil)
	if err != nil {
		log.Warnf("list [%v] keypairs failed: %v", filter, err)
		return 0, keypairList
	}
	defer cur.Close(mctx)
	err = cur.All(mctx, &keypairList)
	if err != nil {
		log.Errorf("decord mgo document failed: %v", err)
	}
	return int(total), keypairList
}
