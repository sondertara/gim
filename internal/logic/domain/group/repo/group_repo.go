package repo

import (
	"context"
	"gim/internal/logic/domain/group/entity"
)

type groupRepo struct{}

var GroupRepo = new(groupRepo)

// Get 获取群组信息
func (*groupRepo) Get(ctx context.Context, groupId int64) (*entity.Group, error) {
	group, err := GroupCache.Get(groupId)
	if err != nil {
		return nil, err
	}
	if group != nil {
		return group, nil
	}
	group, err = GroupDao.Get(groupId)
	if err != nil {
		return nil, err
	}
	err = GroupCache.Set(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Save 获取群组信息
func (*groupRepo) Save(ctx context.Context, group *entity.Group) error {
	groupId := group.Id
	err := GroupDao.Save(group)
	if err != nil {
		return err
	}

	if groupId != 0 {
		err = GroupCache.Del(groupId)
		if err != nil {
			return err
		}
	}
	return nil
}
