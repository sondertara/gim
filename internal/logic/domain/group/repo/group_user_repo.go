package repo

import (
	"context"
	"gim/internal/logic/domain/group/entity"
)

type groupUserRepo struct{}

var GroupUserRepo = new(groupUserRepo)

// GetUsers 获取群组的所有用户信息
func (*groupUserRepo) GetUsers(ctx context.Context, groupId int64) ([]entity.GroupUser, error) {
	users, err := GroupUserCache.Get(groupId)
	if err != nil {
		return nil, err
	}

	if users != nil {
		return users, nil
	}

	users, err = GroupUserDao.ListUser(groupId)
	if err != nil {
		return nil, err
	}

	err = GroupUserCache.Set(groupId, users)
	if err != nil {
		return nil, err
	}
	return users, err
}

// Get 获取
func (*groupUserRepo) Get(ctx context.Context, groupId, userId int64) (*entity.GroupUser, error) {
	return GroupUserDao.Get(groupId, userId)
}

// Save 获取群组的所有用户信息
func (*groupUserRepo) Save(ctx context.Context, groupUser *entity.GroupUser) error {
	err := GroupUserDao.Save(groupUser)
	if err != nil {
		return err
	}

	err = GroupUserCache.Del(groupUser.GroupId)
	if err != nil {
		return err
	}
	return err
}

func (d *groupUserRepo) Delete(ctx context.Context, groupId int64, userId int64) error {
	err := GroupUserDao.Delete(groupId, userId)
	if err != nil {
		return nil
	}

	err = GroupUserCache.Del(groupId)
	if err != nil {
		return err
	}
	return err
}

func (*groupUserRepo) ListByUserId(ctx context.Context, userId int64) ([]entity.Group, error) {
	groups, err := GroupUserDao.ListByUserId(userId)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (*groupUserRepo) BatchGet(groupId int64, userIds []int64) (map[int64]entity.GroupUser, error) {
	return GroupUserDao.BatchGet(groupId, userIds)
}
