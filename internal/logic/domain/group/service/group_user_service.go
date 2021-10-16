package service

import (
	"context"
	"gim/internal/logic/domain/group/entity"
	"gim/internal/logic/domain/group/repo"
	"gim/pkg/gerrors"
	"gim/pkg/logger"
	"gim/pkg/pb"
	"gim/pkg/rpc"
	"time"
)

type groupUserService struct{}

var GroupUserService = new(groupUserService)

// AddUsers 给群组添加用户
func (*groupUserService) AddUsers(ctx context.Context, userId, groupId int64, userIds []int64) ([]int64, error) {
	group, err := repo.GroupRepo.Get(ctx, groupId)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, gerrors.ErrGroupNotExist
	}

	var existIds []int64
	var addedIds []int64

	users, err := repo.GroupUserRepo.BatchGet(groupId, userIds)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for i := range userIds {
		if _, ok := users[userIds[i]]; ok {
			existIds = append(existIds, userIds[i])
			continue
		}

		err = repo.GroupUserRepo.Save(ctx, &entity.GroupUser{
			GroupId:    groupId,
			UserId:     userIds[i],
			MemberType: int(pb.MemberType_GMT_MEMBER),
			CreateTime: now,
			UpdateTime: now,
		})
		if err != nil {
			return nil, err
		}

		addedIds = append(addedIds, userIds[i])
	}

	group.UserNum += int32(len(addedIds))
	err = repo.GroupRepo.Save(ctx, group)
	if err != nil {
		return nil, err
	}

	var addIdMap = make(map[int64]int32, len(addedIds))
	for i := range addedIds {
		addIdMap[addedIds[i]] = 0
	}

	usersResp, err := rpc.BusinessIntClient.GetUsers(ctx, &pb.GetUsersReq{UserIds: addIdMap})
	if err != nil {
		return nil, err
	}
	var members []*pb.GroupMember
	for _, v := range usersResp.Users {
		members = append(members, &pb.GroupMember{
			UserId:    v.UserId,
			Nickname:  v.Nickname,
			Sex:       v.Sex,
			AvatarUrl: v.AvatarUrl,
			UserExtra: v.Extra,
			Remarks:   "",
			Extra:     "",
		})
	}

	userResp, err := rpc.BusinessIntClient.GetUser(ctx, &pb.GetUserReq{UserId: userId})
	if err != nil {
		return nil, err
	}

	err = group.PushToGroup(ctx, groupId, pb.PushCode_PC_ADD_GROUP_MEMBERS, &pb.AddGroupMembersPush{
		OptId:   userResp.User.UserId,
		OptName: userResp.User.Nickname,
		Members: members,
	}, true)
	if err != nil {
		logger.Sugar.Error(err)
	}

	return existIds, nil
}

// DeleteUser 删除用户群组
func (*groupUserService) DeleteUser(ctx context.Context, optId, groupId, userId int64) error {
	group, err := repo.GroupRepo.Get(ctx, groupId)
	if err != nil {
		return err
	}
	if group == nil {
		return gerrors.ErrGroupNotExist
	}

	err = repo.GroupUserRepo.Delete(ctx, groupId, userId)
	if err != nil {
		return err
	}

	group.UserNum -= 1
	err = repo.GroupRepo.Save(ctx, group)
	if err != nil {
		return err
	}

	userResp, err := rpc.BusinessIntClient.GetUser(ctx, &pb.GetUserReq{UserId: optId})
	if err != nil {
		return err
	}
	err = group.PushToGroup(ctx, groupId, pb.PushCode_PC_REMOVE_GROUP_MEMBER, &pb.RemoveGroupMemberPush{
		OptId:         optId,
		OptName:       userResp.User.Nickname,
		DeletedUserId: userId,
	}, true)
	if err != nil {
		return err
	}

	return nil
}
