package Collection

import (
	"fmt"
	"github.com/aiscrm/redisgo"
	"github.com/kevinu2/ngo/constant"
	"github.com/kevinu2/ngo/enum"
	"github.com/kevinu2/ngo/model"
	"time"
)

var c *Collection

type Collection struct {
	Name         string
	GcOffset     int64
	GcMax        int64
	GcOffsetTime int64
	GcMaxTime    int64
	IncrData     model.CollectionData
	LastData     model.CollectionData
	TailsData    []model.CollectionData
	Redis        *redisgo.Cacher
}

func init() {
	c = New()
}

func New() *Collection {
	v := new(Collection)
	v.GcMax = constant.CollectionGcMax
	v.GcOffset = constant.CollectionGcOffset
	return v
}

func AddRedis(redis *redisgo.Cacher) { c.AddRedis(redis) }
func (c *Collection) AddRedis(redis *redisgo.Cacher) {
	c.Redis = redis
}

func AddName(name string) { c.AddName(name) }
func (c *Collection) AddName(name string) {
	c.Name = name
}

func AddGc(max, offset int64) { c.AddGc(offset) }
func (c *Collection) AddGc(offset int64) {
	c.GcOffset = offset
}

func Incr(score int64, member string) { Incr(score, member) }
func (c *Collection) Incr(score int64, member string) error {
	c.IncrData = model.CollectionData{
		Score:  score,
		Member: member,
	}
	err := c.set()
	if err != nil {
		fmt.Printf("Error: Collection %s Set(%d, %s) fails, err: %s", c.Name, score, member, err.Error())
		return err
	}
	if c.GcMax > 0 {
		err = c.calGc(score)
		if err != nil {
			fmt.Printf("Error: Collection %s calGc(%d) fails, err: %s", c.Name, score, err.Error())
		}
	}
	if score > c.GcMaxTime {
		err = c.Gc(c.GcOffsetTime)
		if err != nil {
			fmt.Printf("Error: Collection GC(%s) fails, err: %s", c.Name, err.Error())
			return err
		}
	}
	return nil
}

func (c *Collection) calGc(score int64) error {
	//TODO Auto gc math
	offsetTime := time.Duration(c.GcOffset) * time.Hour
	maxTime := time.Duration(c.GcMax) * time.Hour
	c.GcOffsetTime = score - offsetTime.Milliseconds() - maxTime.Milliseconds()
	c.GcMaxTime = score - maxTime.Milliseconds()
	return nil
}

func (c *Collection) set() error {
	_, err := c.Redis.ZAdd(c.Name, c.IncrData.Score, c.IncrData.Member)
	if err != nil {
		return err
	}
	return nil
}

func Gc(max int64) error { return c.Gc(max) }
func (c *Collection) Gc(max int64) error {
	for {
		rs, err := c.Redis.ZRevrangeByScore(c.Name, max, -1, 0, 100)
		if err != nil {
			return err
		}
		if len(rs) > 0 {
			for k := range rs {
				_, err = c.Redis.ZRem(c.Name, k)
				if err != nil {
					return err
				}
			}
		} else {
			return nil
		}
	}
}
func Last() error { return c.Last() }
func (c *Collection) Last() error {
	rs, err := c.Redis.ZRevrange(c.Name, 0, 0)
	if err != nil {
		return err
	}

	if len(rs) > 0 {
		for k, v := range rs {
			c.LastData.Score = v
			c.LastData.Member = k
		}
	} else {
		c.LastData.Score = 0
		c.LastData.Member = ""
	}
	return nil
}

func Tails(max int64, num int) error { return c.Tails(max, num) }
func (c *Collection) Tails(max int64, num int) error {
	rs, err := c.Redis.ZRevrangeByScore(c.Name, max, -1, 0, num)
	if err != nil {
		return err
	}
	if len(rs) > 0 {
		var tails []model.CollectionData
		for k, v := range rs {
			tails = append(tails, model.CollectionData{
				Score:  v,
				Member: k,
			})
		}
		c.TailsData = tails
		return nil
	} else {
		return enum.ErrorMapEmpty.GetMsg(c.Name)
	}
}
