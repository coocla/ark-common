package misc

import (
	"ark-common/clients/mgo"
	"ark-common/clients/rabbitmq"
	"ark-common/constants"
	"ark-common/resource/navite"
	"ark-common/utils/tool"
	"encoding/json"

	"github.com/streadway/amqp"
)

// SendPullRegionJob 发送拉取地域的作业
func SendPullRegionJob(rbd *mgo.Client, q *rabbitmq.RabbitQueue, accountID, cloudName string) {
	job := &navite.Job{
		Action:    constants.HandleSyncRegion,
		CloudName: cloudName,
		AccountID: accountID,
		JobID:     tool.UUID(),
		Owner:     constants.SYSTEMUSER,
	}
	job.SetPending(rbd)
	body, _ := json.Marshal(job)
	q.Push(constants.LeaderExchange, amqp.ExchangeTopic, constants.SyncJobRoutingKey, body)
}
