package repo

import (
	"gim/internal/logic/domain/group/entity"
	"gim/pkg/db"
	"gim/pkg/gerrors"
	"time"

	"github.com/jinzhu/gorm"
)

type groupDao struct{}

var GroupDao = new(groupDao)

// Get 获取群组信息
func (*groupDao) Get(groupId int64) (*entity.Group, error) {
	var group = entity.Group{Id: groupId}
	err := db.DB.First(&group).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, gerrors.WrapError(err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &group, nil
}

// Save 插入一条群组
func (*groupDao) Save(group *entity.Group) error {
	group.CreateTime = time.Now()
	group.UpdateTime = time.Now()
	err := db.DB.Save(&group).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Update 更新群组信息
func (*groupDao) Update(groupId int64, name, avatarUrl, introduction, extra string) error {
	err := db.DB.Exec("update `group` set name = ?,avatar_url = ?,introduction = ?,extra = ? where id = ?",
		name, avatarUrl, introduction, extra, groupId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// UpdateUserNum 更新群组信息
func (*groupDao) UpdateUserNum(groupId int64, userNum int) error {
	err := db.DB.Exec("update `group` set user_num = user_num + ? where id = ?",
		userNum, groupId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}
