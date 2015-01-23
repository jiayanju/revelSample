package chat

import (
	"github.com/garyburd/redigo/redis"
)

type ChatRedisDao interface {
	SaveMsg(msg string) void
}

type ChatRedisDaoImpl struct {
	pool *redis.Pool
}

func newChatRedisDao(server string, password string) *ChatRedisDaoImpl {
	return &ChatRedisDaoImpl {
		pool : newPool(server, password),
}

func newPool(server string, password string) *redis.Pool {
	return &redis.Pool{
        MaxIdle: 10,
        IdleTimeout: 240 * time.Second,
        Dial: func () (redis.Conn, error) {
            c, err := redis.Dial("tcp", server)
            if err != nil {
                return nil, err
            }
            if _, err := c.Do("AUTH", password); err != nil {
                c.Close()
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

func (c *ChatRedisDaoImpl) SaveMsg(msg string) {
	conn := c.pool.get()
	defer conn

}
