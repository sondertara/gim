package repo

import (
	"gim/internal/business/domain/user/entity"
)

type authRepo struct{}

var AuthRepo = new(authRepo)

func (*authRepo) Get(userId, deviceId int64) (*entity.Device, error) {
	return AuthCache.Get(userId, deviceId)
}

func (*authRepo) Set(userId, deviceId int64, device entity.Device) error {
	return AuthCache.Set(userId, deviceId, device)
}

func (*authRepo) GetAll(userId int64) (map[int64]entity.Device, error) {
	return AuthCache.GetAll(userId)
}
