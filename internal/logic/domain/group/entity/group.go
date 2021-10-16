package entity

import (
	"context"
	"gim/internal/logic/domain/group/repo"
	"gim/internal/logic/proxy"
	"gim/pkg/gerrors"
	"gim/pkg/grpclib"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"gim/pkg/util"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// Group 群组
type Group struct {
	Id           int64     // 群组id
	Name         string    // 组名
	AvatarUrl    string    // 头像
	Introduction string    // 群简介
	UserNum      int32     // 群组人数
	Extra        string    // 附加字段
	CreateTime   time.Time // 创建时间
	UpdateTime   time.Time // 更新时间
}

func (g *Group) ToProto() *pb.Group {
	if g == nil {
		return nil
	}

	return &pb.Group{
		GroupId:      g.Id,
		Name:         g.Name,
		AvatarUrl:    g.AvatarUrl,
		Introduction: g.Introduction,
		UserMum:      g.UserNum,
		Extra:        g.Extra,
		CreateTime:   g.CreateTime.Unix(),
		UpdateTime:   g.UpdateTime.Unix(),
	}
}

func (g *Group) Update(ctx context.Context, userId int64, in *pb.UpdateGroupReq) error {
	g.Name = in.Name
	g.AvatarUrl = in.AvatarUrl
	g.Introduction = in.Introduction
	g.Extra = in.Extra

	userResp, err := rpc.BusinessIntClient.GetUser(ctx, &pb.GetUserReq{UserId: userId})
	if err != nil {
		return err
	}
	err = g.PushToGroup(ctx, g.Id, pb.PushCode_PC_UPDATE_GROUP, &pb.UpdateGroupPush{
		OptId:        userId,
		OptName:      userResp.User.Nickname,
		Name:         g.Name,
		AvatarUrl:    g.AvatarUrl,
		Introduction: g.Introduction,
		Extra:        g.Extra,
	}, true)
	if err != nil {
		return err
	}
	return nil
}

// SendToGroup 消息发送至群组
func (g *Group) SendToGroup(ctx context.Context, sender *pb.Sender, req *pb.SendMessageReq) (int64, error) {
	users, err := repo.GroupUserRepo.GetUsers(ctx, req.ReceiverId)
	if err != nil {
		return 0, err
	}

	if sender.SenderType == pb.SenderType_ST_USER && !IsInGroup(users, sender.SenderId) {
		logger.Sugar.Error(ctx, sender.SenderId, req.ReceiverId, "不在群组内")
		return 0, gerrors.ErrNotInGroup
	}

	// 如果发送者是用户，将消息发送给发送者,获取用户seq
	var userSeq int64
	if sender.SenderType == pb.SenderType_ST_USER {
		userSeq, err = proxy.MessageProxy.SendToUser(ctx, sender, sender.SenderId, req)
		if err != nil {
			return 0, err
		}
	}

	go func() {
		defer util.RecoverPanic()
		// 将消息发送给群组用户，使用写扩散
		for _, user := range users {
			// 前面已经发送过，这里不需要再发送
			if sender.SenderType == pb.SenderType_ST_USER && user.UserId == sender.SenderId {
				continue
			}
			_, err := proxy.MessageProxy.SendToUser(grpclib.NewAndCopyRequestId(ctx), sender, user.UserId, req)
			if err != nil {
				return
			}
		}
	}()

	return userSeq, nil
}

func IsInGroup(users []GroupUser, userId int64) bool {
	for i := range users {
		if users[i].UserId == userId {
			return true
		}
	}
	return false
}

// PushToGroup 向群组推送消息
func (g *Group) PushToGroup(ctx context.Context, groupId int64, code pb.PushCode, message proto.Message, isPersist bool) error {
	logger.Logger.Debug("push_to_group",
		zap.Int64("request_id", grpclib.GetCtxRequestId(ctx)),
		zap.Int64("group_id", groupId),
		zap.Int32("code", int32(code)),
		zap.Any("message", message))

	messageBuf, err := proto.Marshal(message)
	if err != nil {
		return gerrors.WrapError(err)
	}

	commandBuf, err := proto.Marshal(&pb.Command{Code: int32(code), Data: messageBuf})
	if err != nil {
		return gerrors.WrapError(err)
	}

	_, err = g.SendToGroup(ctx,
		&pb.Sender{
			SenderType: pb.SenderType_ST_SYSTEM,
			SenderId:   0,
		},
		&pb.SendMessageReq{
			ReceiverType:   pb.ReceiverType_RT_GROUP,
			ReceiverId:     groupId,
			ToUserIds:      nil,
			MessageType:    pb.MessageType_MT_COMMAND,
			MessageContent: commandBuf,
			SendTime:       util.UnixMilliTime(time.Now()),
			IsPersist:      isPersist,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

type GroupUserUpdate struct {
	GroupId int64  `json:"group_id"` // 群组id
	UserId  int64  `json:"user_id"`  // 用户id
	Label   string `json:"label"`    // 用户标签
	Extra   string `json:"extra"`    // 群组用户附件属性
}
