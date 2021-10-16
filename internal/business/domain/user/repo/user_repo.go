package repo

import (
	"context"
	"gim/internal/business/domain/user/entity"
)

type userRepo struct{}

var UserRepo = new(userRepo)

// Get 获取单个用户
func (*userRepo) Get(ctx context.Context, userId int64) (*entity.User, error) {
	user, err := UserCache.Get(userId)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	user, err = UserDao.Get(userId)
	if err != nil {
		return nil, err
	}

	if user != nil {
		err = UserCache.Set(*user)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

func (*userRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error) {
	return UserDao.GetByPhoneNumber(phoneNumber)
}

// GetByIds 获取多个用户
func (*userRepo) GetByIds(ctx context.Context, userIds []int64) ([]entity.User, error) {
	return UserDao.GetByIds(userIds)
}

// Search 搜索用户
func (*userRepo) Search(ctx context.Context, key string) ([]entity.User, error) {
	return UserDao.Search(key)
}

// Save 保存用户
func (*userRepo) Save(ctx context.Context, user *entity.User) error {
	userId := user.Id
	err := UserDao.Save(user)
	if err != nil {
		return err
	}

	if userId != 0 {
		err = UserCache.Del(user.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
