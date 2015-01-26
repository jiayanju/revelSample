package chatservice

type ChatRedisService interface {
	SaveValueToList(key string, value string)

	GetValues(key string) string
}

type ChatRedisServiceImpl struct {
	chatRedisDao ChatRedisDao
}

func (c *ChatRedisServiceImpl) SaveValueToList(key string, value string) {
	c.chatRedisDao.SaveValueToList(key, value)
}

func (c *ChatRedisServiceImpl) GetValues(key string) string {
	return c.chatRedisDao.GetValues(key)
}

func NewChatRedisService() ChatRedisService {
	return &ChatRedisServiceImpl{
		chatRedisDao: NewChatRedisDao("127.0.0.1:6379")}
}
