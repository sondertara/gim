package repo

import (
	"gim/internal/logic/domain/group/entity"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"gim/pkg/util"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	GroupKey    = "group:"
	GroupExpire = 2 * time.Hour
)

var RedisUtil = util.NewRedisUtil(db.RedisCli)

type groupCache struct{}

var GroupCache = new(groupCache)

// Get 获取群组缓存
func (c *groupCache) Get(groupId int64) (*entity.Group, error) {
	var user entity.Group
	err := RedisUtil.Get(GroupKey+strconv.FormatInt(groupId, 10), &user)
	if err != nil && err != redis.Nil {
		return nil, gerrors.WrapError(err)
	}
	if err == redis.Nil {
		return nil, nil
	}
	return &user, nil
}

// Set 设置群组缓存
func (c *groupCache) Set(group *entity.Group) error {
	err := RedisUtil.Set(GroupKey+strconv.FormatInt(group.Id, 10), group, GroupExpire)
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Del 删除群组缓存
func (c *groupCache) Del(groupId int64) error {
	_, err := db.RedisCli.Del(GroupKey + strconv.FormatInt(groupId, 10)).Result()
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}
