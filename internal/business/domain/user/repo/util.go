package repo

import (
	"gim/pkg/db"
	"gim/pkg/util"
)

var RedisUtil = util.NewRedisUtil(db.RedisCli)
