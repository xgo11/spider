package spider

import (
	"errors"
	"fmt"
)
import (
	"github.com/xgo11/redis4g"
	"github.com/xgo11/spider/core"
)

var (
	_ core.IQueue = &redisQueue{}
)

type redisQueue struct {
	name   string
	qName  string
	client *redis4g.WrapClient
}

func NewRedisQueue(name string, confPath string) core.IQueue {
	client := redis4g.Connect(confPath)
	if client == nil {
		panic(fmt.Sprintf("connect queue %v fail", confPath))
	}
	return &redisQueue{name: name, qName: "sys:" + name, client: client}
}

func (rq *redisQueue) Name() string {
	return rq.name
}

func (rq *redisQueue) Put(message ...string) error {

	var its = make([]interface{}, len(message))
	for i, msg := range message {
		its[i] = msg
	}

	if cnt := rq.client.LPush(rq.qName, its...); cnt > 0 {
		return nil
	}

	return errors.New("put message fail")
}

func (rq *redisQueue) Pop(count ...int) []string {
	cnt := 1
	if len(count) > 0 && count[0] > 0 {
		cnt = count[0]
	}

	msgArr := make([]string, 0, cnt)
	for ; cnt > 0; cnt-- {
		if msg := rq.client.RPop(rq.qName); msg != "" {
			msgArr = append(msgArr, msg)
			continue
		}
		break
	}
	return msgArr

}

func (rq *redisQueue) Size() int {
	return int(rq.client.LLen(rq.qName))
}

func (rq *redisQueue) Limit() int {
	return 0
}
