package chatservice

import (
	"github.com/revel/revel"
)

type ChatRedisService interface {
	SaveValueToList(key string, value string)

	GetValues(key string) []string
}

type ChatRedisServiceImpl struct {
	chatRedisDao ChatRedisDao
}

func (c *ChatRedisServiceImpl) SaveValueToList(key string, value string) {
	c.chatRedisDao.SaveValueToList(key, value)
}

func (c *ChatRedisServiceImpl) GetValues(key string) []string {
	return c.chatRedisDao.GetValues(key)
}

func NewChatRedisService() ChatRedisService {
	server := revel.Config.StringDefault("redis.server.address", "11.11.11.11:6379")
	return &ChatRedisServiceImpl{
		chatRedisDao: NewChatRedisDao(server)}
}
