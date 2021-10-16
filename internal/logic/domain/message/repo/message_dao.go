package repo

import (
	"fmt"
	"gim/internal/logic/domain/message/entity"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

const messageTableNum = 1

type messageRepo struct{}

var MessageRepo = new(messageRepo)

func (*messageRepo) tableName(userId int64) string {
	return fmt.Sprintf("message_%03d", userId%messageTableNum)
}

// Save 插入一条消息
func (d *messageRepo) Save(message entity.Message) error {
	err := db.DB.Table(d.tableName(message.UserId)).Create(&message).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// ListBySeq 根据类型和id查询大于序号大于seq的消息
func (d *messageRepo) ListBySeq(userId, seq, limit int64) ([]entity.Message, bool, error) {
	db := db.DB.Table(d.tableName(userId)).
		Where("user_id = ? and seq > ?", userId, seq)

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, false, gerrors.WrapError(err)
	}
	if count == 0 {
		return nil, false, nil
	}

	var messages []entity.Message
	err = db.Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, false, gerrors.WrapError(err)
	}
	return messages, count > limit, nil
}
