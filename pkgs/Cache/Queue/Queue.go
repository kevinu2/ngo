package Queue

import (
	"github.com/gomodule/redigo/redis"
	"github.com/kevinu2/ngo/v2/pkgs/Default"
	"github.com/kevinu2/ngo/v2/pkgs/RedisGo"
)

type Queue struct {
	QueueName    string
	Redis        *RedisGo.Cacher
	RedisTimeout int
}

func (q Queue) Produce(data []string) error {
	for _, v := range data {
		err := q.Redis.LPush(q.QueueName, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (q Queue) Consume() (string, error) {
	for {
		rs, err := q.Redis.BRPop(q.QueueName, q.RedisTimeout)
		if err != nil {
			return Default.StringEmpty, err
		} else {
			if rs == nil {
				continue
				//return Default.DEFAULT_EMPTY, enum.ErrorQueueEmpty.GetMsg("q.QueueName")
			} else {
				return redis.String(rs, err)
			}
		}
	}
}
