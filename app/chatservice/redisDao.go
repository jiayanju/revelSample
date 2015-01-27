package chatservice

import (
	"github.com/garyburd/redigo/redis"
	"github.com/revel/revel"
	"time"
)

type ChatRedisDao interface {
	SaveValueToList(key string, value string)

	GetValues(key string) []string
}

type ChatRedisDaoImpl struct {
	pool *redis.Pool
}

func (c *ChatRedisDaoImpl) SaveValueToList(key string, value string) {
	conn := c.pool.Get()
	defer conn.Close()

	conn.Do("RPUSH", key, value)
}

func (c *ChatRedisDaoImpl) GetValues(key string) []string {
	conn := c.pool.Get()
	defer conn.Close()

	values, err := redis.Strings(conn.Do("LRANGE", key, 0, -1))
	if err != nil {
		revel.TRACE.Println(err)
	}
	return values
}

func NewChatRedisDao(server string) ChatRedisDao {
	pool := newPool(server)
	return &ChatRedisDaoImpl{pool: pool}
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
