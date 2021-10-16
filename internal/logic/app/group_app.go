package app

import (
	"context"
	"gim/internal/logic/domain/group/repo"
	"gim/internal/logic/domain/group/service"
	"gim/pkg/pb"
	"time"
)

type groupApp struct{}

var GroupApp = new(groupApp)

func (*groupApp) CreateGroup(ctx context.Context, userId int64, in *pb.CreateGroupReq) (int64, error) {
	return service.GroupService.Create(ctx, userId, in)
}

func (*groupApp) GetGroup(ctx context.Context, groupId int64) (*pb.Group, error) {
	group, err := repo.GroupRepo.Get(ctx, groupId)
	if err != nil {
		return nil, err
	}

	return group.ToProto(), nil
}

func (*groupApp) GetUserGroups(ctx context.Context, userId int64) ([]*pb.Group, error) {
	groups, err := repo.GroupUserRepo.ListByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	pbGroups := make([]*pb.Group, len(groups))
	for i := range groups {
		pbGroups[i] = groups[i].ToProto()
	}
	return pbGroups, nil
}

// Update 更新群组
func (*groupApp) Update(ctx context.Context, userId int64, update *pb.UpdateGroupReq) error {
	group, err := repo.GroupRepo.Get(ctx, update.GroupId)
	if err != nil {
		return err
	}

	err = group.Update(ctx, userId, update)
	if err != nil {
		return err
	}

	err = repo.GroupRepo.Save(ctx, group)
	if err != nil {
		return err
	}
	return nil
}

func (*groupApp) AddGroupMembers(ctx context.Context, userId, groupId int64, userIds []int64) ([]int64, error) {
	return service.GroupUserService.AddUsers(ctx, userId, groupId, userIds)
}

// UpdateUser 更新群组用户
func (*groupApp) UpdateUser(ctx context.Context, in *pb.UpdateGroupMemberReq) error {
	user, err := repo.GroupUserRepo.Get(ctx, in.GroupId, in.UserId)
	if err != nil {
		return nil
	}
	if user == nil {
		return nil
	}

	user.MemberType = int(in.MemberType)
	user.Remarks = in.Remarks
	user.Extra = in.Extra
	user.UpdateTime = time.Now()

	err = repo.GroupUserRepo.Save(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (*groupApp) DeleteMember(ctx context.Context, groupId int64, userId int64, optId int64) error {
	return service.GroupUserService.DeleteUser(ctx, groupId, userId, optId)
}

// GetGroupMembers 获取群组成员
func (*groupApp) GetGroupMembers(ctx context.Context, groupId int64) ([]*pb.GroupMember, error) {
	return service.GroupService.GetUsers(ctx, groupId)
}

// SendMessage 获取群组成员
func (*groupApp) SendMessage(ctx context.Context, sender *pb.Sender, req *pb.SendMessageReq) (int64, error) {
	group, err := repo.GroupRepo.Get(ctx, req.ReceiverId)
	if err != nil {
		return 0, err
	}

	return group.SendToGroup(ctx, sender, req)
}
