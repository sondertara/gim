package service

import (
	"context"
	"gim/internal/logic/domain/group/entity"
	"gim/internal/logic/domain/group/repo"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"time"
)

type groupService struct{}

var GroupService = new(groupService)

// Create 创建群组
func (*groupService) Create(ctx context.Context, userId int64, in *pb.CreateGroupReq) (int64, error) {
	now := time.Now()
	group := &entity.Group{
		Name:         in.Name,
		AvatarUrl:    in.AvatarUrl,
		Introduction: in.Introduction,
		Extra:        in.Extra,
		CreateTime:   now,
		UpdateTime:   now,
	}

	err := repo.GroupRepo.Save(ctx, group)
	if err != nil {
		return 0, err
	}

	// 创建者添加为管理员
	repo.GroupUserRepo.Save(ctx, &entity.GroupUser{
		GroupId:    group.Id,
		UserId:     userId,
		MemberType: int(pb.MemberType_GMT_ADMIN),
		CreateTime: now,
		UpdateTime: now,
	})
	if err != nil {
		return 0, err
	}

	// 其让人添加为成员
	for i := range in.MemberIds {
		err = repo.GroupUserRepo.Save(ctx, &entity.GroupUser{
			GroupId:    group.Id,
			UserId:     in.MemberIds[i],
			MemberType: int(pb.MemberType_GMT_MEMBER),
			CreateTime: now,
			UpdateTime: now,
		})
		if err != nil {
			return 0, err
		}
	}
	return group.Id, nil
}

// GetUsers 获取群组用户
func (s *groupService) GetUsers(ctx context.Context, groupId int64) ([]*pb.GroupMember, error) {
	group, err := repo.GroupRepo.Get(ctx, groupId)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	members, err := repo.GroupUserRepo.GetUsers(ctx, groupId)
	if err != nil {
		return nil, err
	}

	userIds := make(map[int64]int32, len(members))
	for i := range members {
		userIds[members[i].UserId] = 0
	}
	resp, err := rpc.BusinessIntClient.GetUsers(ctx, &pb.GetUsersReq{UserIds: userIds})
	if err != nil {
		return nil, err
	}

	var infos = make([]*pb.GroupMember, len(members))
	for i := range members {
		member := pb.GroupMember{
			UserId:     members[i].UserId,
			MemberType: pb.MemberType(members[i].MemberType),
			Remarks:    members[i].Remarks,
			Extra:      members[i].Extra,
		}

		user, ok := resp.Users[members[i].UserId]
		if ok {
			member.Nickname = user.Nickname
			member.Sex = user.Sex
			member.AvatarUrl = user.AvatarUrl
			member.UserExtra = user.Extra
		}
		infos[i] = &member
	}

	return infos, nil
}
