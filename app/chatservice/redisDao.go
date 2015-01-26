package chatservice

import (
	"github.com/garyburd/redigo/redis"
	"github.com/revel/revel"
	"time"
)

type ChatRedisDao interface {
	SaveValueToList(key string, value string)

	GetValues(key string) string
}

type ChatRedisDaoImpl struct {
	pool *redis.Pool
}

func (c *ChatRedisDaoImpl) SaveValueToList(key string, value string) {
	conn := c.pool.Get()
	defer conn.Close()

	conn.Do("RPUSH", key, value)
}

func (c *ChatRedisDaoImpl) GetValues(key string) (value string) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("Get", key)
	if err == nil {
		revel.TRACE.Println(err)
	}
	return value
}

func NewChatRedisDao(server string) ChatRedisDao {
	return &ChatRedisDaoImpl{
		pool: newPool(server)}
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
