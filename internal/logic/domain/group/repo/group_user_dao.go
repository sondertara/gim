package repo

import (
	"gim/internal/logic/domain/group/entity"
	"gim/pkg/db"
	"gim/pkg/gerrors"

	"github.com/jinzhu/gorm"
)

type groupUserDao struct{}

var GroupUserDao = new(groupUserDao)

// ListByUserId 获取用户加入的群组信息
func (*groupUserDao) ListByUserId(userId int64) ([]entity.Group, error) {
	var groups []entity.Group
	err := db.DB.Select("g.id,g.name,g.avatar_url,g.introduction,g.user_num,g.extra,g.create_time,g.update_time").
		Table("group_user u").
		Joins("join `group` g on u.group_id = g.id").
		Where("u.user_id = ?", userId).
		Find(&groups).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return groups, nil
}

// ListUser 获取群组用户信息
func (*groupUserDao) ListUser(groupId int64) ([]entity.GroupUser, error) {
	var groupUsers []entity.GroupUser
	err := db.DB.Find(&groupUsers, "group_id = ?", groupId).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return groupUsers, nil
}

// Get 获取群组用户信息,用户不存在返回nil
func (*groupUserDao) Get(groupId, userId int64) (*entity.GroupUser, error) {
	var groupUser entity.GroupUser
	err := db.DB.First(&groupUser, "group_id = ? and user_id = ?", groupId, userId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, gerrors.WrapError(err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &groupUser, nil
}

// BatchGet 批量获取群组用户信息
func (*groupUserDao) BatchGet(groupId int64, userIds []int64) (map[int64]entity.GroupUser, error) {
	var users []entity.GroupUser
	err := db.DB.Find(&users, "group_id = ? and user_id in (?)", groupId, userIds).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	userMap := make(map[int64]entity.GroupUser, len(users))
	for i := range users {
		userMap[users[i].UserId] = users[i]
	}

	return userMap, nil
}

// Save 将用户添加到群组
func (*groupUserDao) Save(groupUser *entity.GroupUser) error {
	err := db.DB.Save(&groupUser).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Delete 将用户从群组删除
func (d *groupUserDao) Delete(groupId int64, userId int64) error {
	err := db.DB.Exec("delete from group_user where group_id = ? and user_id = ?",
		groupId, userId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Update 更新用户群组信息
func (*groupUserDao) Update(user entity.GroupUser) error {
	err := db.DB.Exec("update group_user set member_type = ?,remarks = ?,extra = ? where group_id = ? and user_id = ?",
		user.MemberType, user.Remarks, user.Extra, user.GroupId, user.UserId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}
