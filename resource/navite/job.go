package navite

import (
	"ark-common/clients/mgo"
	"ark-common/constants"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
)

// JobTable 作业
const JobTable = "jobs"

// Job 资源同步作业
type Job struct {
	Action      string      `bson:"action" json:"action"`
	CloudName   string      `bson:"cloudName" json:"cloudName"`
	AccountID   string      `bson:"accountId" json:"accountId"`
	JobID       string      `bson:"jobId" json:"jobId"`
	Owner       string      `bson:"owner" json:"owner"`
	Status      string      `bson:"status" json:"status"`
	RegionID    string      `bson:"regionId" json:"regionId"`
	Params      interface{} `bson:"params" json:"-"`
	Reason      string      `bson:"reason" json:"reason"`
	CreatedTime time.Time   `bson:"createdTime" json:"createdTime"`
	StatusTime  time.Time   `bson:"statusTime" json:"statusTime"`
}

func (j *Job) SetPending(rbd *mgo.Client) {
	j.Status = constants.PENDINGJOB
	j.CreatedTime = time.Now()
	j.StatusTime = time.Now()
	_, err := rbd.Table(JobTable).Insert(j)
	if err != nil {
		log.Errorf("insert into [%+v] failed: %v", j, err)
	}
}

func (j *Job) SetWorking(rbd *mgo.Client) {
	j.Status = constants.WORKINGJOB
	j.StatusTime = time.Now()
	filter := bson.M{
		"jobId": j.JobID,
	}
	err := rbd.Table(JobTable).Replace(filter, j, nil)
	if err != nil {
		log.Errorf("replace job [%+v] failed: %v", j, err)
	}
}

func (j *Job) SetFailed(rbd *mgo.Client) {
	j.Status = constants.FAILEDJOB
	j.StatusTime = time.Now()
	filter := bson.M{
		"jobId": j.JobID,
	}
	err := rbd.Table(JobTable).Replace(filter, j, nil)
	if err != nil {
		log.Errorf("replace job [%+v] failed: %v", j, err)
	}
}

func (j *Job) SetCancel(rbd *mgo.Client, reason string) {
	j.Status = constants.CANCELJOB
	j.StatusTime = time.Now()
	j.Reason = reason
	filter := bson.M{
		"jobId": j.JobID,
	}
	err := rbd.Table(JobTable).Replace(filter, j, nil)
	if err != nil {
		log.Errorf("replace job [%+v] failed: %v", j, err)
	}
}

func (j *Job) SetSuccess(rbd *mgo.Client) {
	j.Status = constants.SUCCESSJOB
	j.StatusTime = time.Now()
	filter := bson.M{
		"jobId": j.JobID,
	}
	err := rbd.Table(JobTable).Replace(filter, j, nil)
	if err != nil {
		log.Errorf("replace job [%+v] failed: %v", j, err)
	}
}
