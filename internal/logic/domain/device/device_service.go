package device

import (
	"context"
	"gim/pkg/pb"
	"gim/pkg/rpc"
)

const (
	DeviceOnline  = 1
	DeviceOffline = 0
)

type deviceService struct{}

var DeviceService = new(deviceService)

// Register 注册设备
func (*deviceService) Register(ctx context.Context, device *Device) error {
	err := DeviceDao.Save(device)
	if err != nil {
		return err
	}

	return nil
}

// SignIn 长连接登录
func (*deviceService) SignIn(ctx context.Context, userId, deviceId int64, token string, connAddr string, clientAddr string) error {
	_, err := rpc.BusinessIntClient.Auth(ctx, &pb.AuthReq{UserId: userId, DeviceId: deviceId, Token: token})
	if err != nil {
		return err
	}

	// 标记用户在设备上登录
	device, err := DeviceRepo.Get(deviceId)
	if err != nil {
		return err
	}
	if device == nil {
		return nil
	}

	device.Online(userId, connAddr, clientAddr)

	err = DeviceRepo.Save(device)
	if err != nil {
		return err
	}
	return nil
}

// Auth 权限验证
func (*deviceService) Auth(ctx context.Context, userId, deviceId int64, token string) error {
	_, err := rpc.BusinessIntClient.Auth(ctx, &pb.AuthReq{UserId: userId, DeviceId: deviceId, Token: token})
	if err != nil {
		return err
	}
	return nil
}

func (*deviceService) ListOnlineByUserId(ctx context.Context, userId int64) ([]*pb.Device, error) {
	devices, err := DeviceRepo.ListOnlineByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	pbDevices := make([]*pb.Device, len(devices))
	for i := range devices {
		pbDevices[i] = devices[i].ToProto()
	}
	return pbDevices, nil
}

// ServerStop 设备离线
func (*deviceService) ServerStop(ctx context.Context, connAddr string) error {
	devices, err := DeviceDao.ListOnlineByConnAddr(connAddr)
	if err != nil {
		return err
	}

	err = DeviceDao.UpdateStatusByCoonAddr(connAddr, DeviceOffLine)
	if err != nil {
		return err
	}

	var userIds = make([]int64, 0, len(devices))
	for i := range devices {
		userIds = append(userIds, devices[i].UserId)
	}

	// todo 多实例Redis会有问题
	err = UserDeviceCache.Del(userIds...)
	if err != nil {
		return err
	}
	return nil
}
