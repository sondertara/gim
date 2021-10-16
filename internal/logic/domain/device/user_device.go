package device

import (
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"gim/pkg/util"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	UserDeviceKey    = "user_device:"
	UserDeviceExpire = 2 * time.Hour
)

var RedisUtil = util.NewRedisUtil(db.RedisCli)

type userDeviceCache struct{}

var UserDeviceCache = new(userDeviceCache)

// Get 获取指定用户的所有在线设备
func (c *userDeviceCache) Get(userId int64) ([]Device, error) {
	var devices []Device
	err := RedisUtil.Get(UserDeviceKey+strconv.FormatInt(userId, 10), &devices)
	if err != nil && err != redis.Nil {
		return nil, gerrors.WrapError(err)
	}

	if err == redis.Nil {
		return nil, nil
	}
	return devices, nil
}

// Set 将指定用户的所有在线设备存入缓存
func (c *userDeviceCache) Set(userId int64, devices []Device) error {
	err := RedisUtil.Set(UserDeviceKey+strconv.FormatInt(userId, 10), devices, UserDeviceExpire)
	return gerrors.WrapError(err)
}

// Del 删除用户的在线设备列表
func (c *userDeviceCache) Del(userIds ...int64) error {
	var ids = make([]string, len(userIds))
	for i := range userIds {
		ids[i] = UserDeviceKey + strconv.FormatInt(userIds[i], 10)
	}

	_, err := db.RedisCli.Del(ids...).Result()
	return gerrors.WrapError(err)
}
