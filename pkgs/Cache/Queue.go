package Cache

import (
	r "github.com/gomodule/redigo/redis"
	"ngo/constant"
	db "ngo/pkgs/Db"
)

type Queue struct {
	QueueName    string
	DB           *db.Pool
	RedisTimeout int
}

func (q Queue) Produce(data []string) error {
	for _, v := range data {
		err := q.DB.Redis.LPush(q.QueueName, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (q Queue) Consume() (string, error) {
	for {
		rs, err := q.DB.Redis.BRPop(q.QueueName, q.RedisTimeout)
		if err != nil {
			return constant.DefaultEmpty, err
		} else {
			if rs == nil {
				continue
				//return constant.DEFAULT_EMPTY, enum.ErrorQueueEmpty.GetMsg("q.QueueName")
			} else {
				return r.String(rs, err)
			}
		}
	}
}
