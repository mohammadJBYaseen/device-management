package service

import (
	"context"
	"device-management/exception"
	"device-management/model"
	"device-management/repository"
	"fmt"

	"github.com/google/uuid"
)

type (
	DeviceService interface {
		CreateDevice(ctx context.Context, device model.Device) (model.Device, error)
		GetDevice(ctx context.Context, deviceUuid uuid.UUID) (model.Device, error)
		UpdateDevice(ctx context.Context, deviceUuid uuid.UUID, device model.Device) (model.Device, error)
		PatchDevice(ctx context.Context, deviceUuid uuid.UUID, device []model.JsonPatch) (model.Device, error)
		DeleteDevice(ctx context.Context, deviceUuid uuid.UUID) error
		SearchDevices(ctx context.Context, request model.SearchRequest) (model.Page[model.Device], error)
	}

	deviceServiceImpl struct {
		repo repository.DeviceRepository
	}
)

func NewDeviceService(repo repository.DeviceRepository) DeviceService {
	return &deviceServiceImpl{
		repo: repo,
	}
}

func (service *deviceServiceImpl) CreateDevice(ctx context.Context, device model.Device) (model.Device, error) {
	device.UUID = uuid.New()
	return service.repo.CreateDevice(ctx, device)
}

func (service *deviceServiceImpl) UpdateDevice(ctx context.Context, deviceUuid uuid.UUID, updated model.Device) (model.Device, error) {
	source, err := service.repo.GetDevice(ctx, deviceUuid)
	if err != nil {
		return source, err
	}
	return service.repo.UpdateDevice(ctx, merge(source, updated))
}

func (service *deviceServiceImpl) PatchDevice(ctx context.Context, deviceUuid uuid.UUID, jsonPatch []model.JsonPatch) (model.Device, error) {
	sourceDev, err := service.GetDevice(ctx, deviceUuid)
	if err != nil {
		return model.Device{}, err
	}
	if updatedDevice, err2 := patchMerge(&sourceDev, jsonPatch); err2 != nil {
		return model.Device{}, err2
	} else {
		return service.UpdateDevice(ctx, deviceUuid, *updatedDevice)
	}
}

func (service *deviceServiceImpl) DeleteDevice(ctx context.Context, deviceUuid uuid.UUID) error {
	return service.repo.DeleteDevice(ctx, deviceUuid)
}

func (service *deviceServiceImpl) SearchDevices(ctx context.Context, request model.SearchRequest) (model.Page[model.Device], error) {
	return service.repo.SearchDevices(ctx, request)
}

func merge(source model.Device, target model.Device) model.Device {
	source.Name = target.Name
	source.BrandName = target.BrandName
	return source
}

func (service *deviceServiceImpl) GetDevice(ctx context.Context, deviceUuid uuid.UUID) (model.Device, error) {
	return service.repo.GetDevice(ctx, deviceUuid)
}

func patchMerge(source *model.Device, target []model.JsonPatch) (*model.Device, error) {
	for _, patch := range target {
		if patch.Op != "replace" {
			return nil, exception.BadRequest{
				Message: fmt.Sprintf("Unsupported patch op: %s", patch.Op),
			}
		}
		if patch.Path == "name" {
			source.Name = patch.Value
		} else if patch.Path == "brand" {
			source.BrandName = patch.Value
		} else {
			return nil, exception.BadRequest{
				Message: fmt.Sprintf("Unsupported patch op: %s", patch.Op),
			}
		}
	}
	return source, nil
}
