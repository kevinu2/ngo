package Cache

import (
	"github.com/aiscrm/redisgo"
	"github.com/kevinu2/ngo/enum"
	log "ngo/pkgs/Log"
)

type CollectionData struct {
	Score  int64
	Member string
}

type Collection struct {
	Name      string
	Offset    int64
	IncrData  CollectionData
	LastData  CollectionData
	TailsData []CollectionData
	Redis     *redisgo.Cacher
}

func (c *Collection) Set() error {
	_, err := c.Redis.ZAdd(c.Name, c.IncrData.Score, c.IncrData.Member)
	if err != nil {
		return err
	}

	if c.Offset > 0 {
		err := c.Gc(c.Offset)
		if err != nil {
			log.Logger().Errorf("Collection GC(%s) fails, err: %s", c.Name, err.Error())
		}

	}
	return nil
}

func (c *Collection) Gc(max int64) error {
	// ZRevRangeByScore Z:OPTION:TIMELINE:1:1 1609739148000 -1 LIMIT 0 10
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

func (c *Collection) Tails(max int64, num int) error {
	rs, err := c.Redis.ZRevrangeByScore(c.Name, max, -1, 0, num)
	if err != nil {
		return err
	}
	if len(rs) > 0 {
		for k, v := range rs {
			c.TailsData = append(c.TailsData, CollectionData{
				Score:  v,
				Member: k,
			})
		}
		return nil
	} else {
		return enum.ErrorMapEmpty.GetMsg(c.Name)
	}
}
