package friend

import (
	"gim/pkg/db"
	"gim/pkg/gerrors"

	"github.com/jinzhu/gorm"
)

type friendRepo struct{}

var FriendRepo = new(friendRepo)

// Get 获取好友
func (*friendRepo) Get(userId, friendId int64) (*Friend, error) {
	friend := Friend{}
	err := db.DB.First(&friend, "user_id = ? and friend_id = ?", userId, friendId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &friend, nil
}

// Save 添加好友
func (*friendRepo) Save(friend *Friend) error {
	return gerrors.WrapError(db.DB.Save(&friend).Error)
}

// Update 更新好友
func (*friendRepo) Update(friend Friend) error {
	err := db.DB.Model(&friend).Where("user_id = ? and friend_id = ?", friend.UserId, friend.FriendId).
		Updates(
			map[string]interface{}{
				"remarks": friend.Remarks,
				"extra":   friend.Extra,
			},
		).Error
	return gerrors.WrapError(err)
}

// UpdateStatus 更新好友状态
func (*friendRepo) UpdateStatus(userId, friendId int64, status int) error {
	err := db.DB.Model(&Friend{}).Where("user_id = ? and friend_id = ?", userId, friendId).
		Updates(map[string]interface{}{
			"status": status,
		}).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// List 获取好友列表
func (*friendRepo) List(userId int64, status int) ([]Friend, error) {
	var friends []Friend
	err := db.DB.Where("user_id = ? and status = ?", userId, status).Find(&friends).Error
	return friends, gerrors.WrapError(err)
}
