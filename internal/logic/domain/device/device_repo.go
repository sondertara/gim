package device

import (
	"context"
)

type deviceRepo struct{}

var DeviceRepo = new(deviceRepo)

// Get 获取设备
func (*deviceRepo) Get(deviceId int64) (*Device, error) {
	device, err := DeviceCache.Get(deviceId)
	if err != nil {
		return nil, err
	}

	if device != nil {
		return device, nil
	}

	device, err = DeviceDao.Get(deviceId)
	if err != nil {
		return nil, err
	}

	if device != nil {
		err = DeviceCache.Set(device)
		if err != nil {
			return nil, err
		}
	}
	return device, nil
}

// Save 保存设备信息
func (*deviceRepo) Save(device *Device) error {
	deviceId := device.Id
	err := DeviceDao.Save(device)
	if err != nil {
		return err
	}

	if deviceId != 0 {
		err = DeviceCache.Del(deviceId)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListOnlineByUserId 获取用户的所有在线设备
func (*deviceRepo) ListOnlineByUserId(ctx context.Context, userId int64) ([]Device, error) {
	devices, err := UserDeviceCache.Get(userId)
	if err != nil {
		return nil, err
	}

	if devices != nil {
		return devices, nil
	}

	devices, err = DeviceDao.ListOnlineByUserId(userId)
	if err != nil {
		return nil, err
	}

	err = UserDeviceCache.Set(userId, devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// Update 更新设备绑定用户和设备在线状态
func (*deviceRepo) Update(deviceId, userId int64, status int, connAddr string, clientAddr string) error {
	return DeviceDao.Update(deviceId, userId, status, connAddr, clientAddr)
}

// UpdateStatus 更新设备的在线状态
func (*deviceRepo) UpdateStatus(deviceId int64, status int) error {
	return DeviceDao.UpdateStatus(deviceId, status)
}

// Upgrade 升级设备
func (*deviceRepo) Upgrade(deviceId int64, systemVersion, sdkVersion string) error {
	return DeviceDao.Upgrade(deviceId, systemVersion, sdkVersion)
}

// ListOnlineByConnAddr 查询用户所有的在线设备
func (*deviceRepo) ListOnlineByConnAddr(connAddr string) ([]Device, error) {
	return DeviceDao.ListOnlineByConnAddr(connAddr)
}

// UpdateStatusByCoonAddr 更新在线状态
func (*deviceRepo) UpdateStatusByCoonAddr(connAddr string, status int) error {
	return DeviceDao.UpdateStatusByCoonAddr(connAddr, status)
}
