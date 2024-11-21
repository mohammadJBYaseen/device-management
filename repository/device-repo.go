package repository

import (
	"context"
	"device-management/exception"
	"device-management/model"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	DeviceRepository interface {
		CreateDevice(ctx context.Context, device model.Device) (model.Device, error)
		GetDevice(ctx context.Context, id uuid.UUID) (model.Device, error)
		UpdateDevice(ctx context.Context, device model.Device) (model.Device, error)
		DeleteDevice(ctx context.Context, deviceUuid uuid.UUID) error
		SearchDevices(ctx context.Context, request model.SearchRequest) (model.Page[model.Device], error)
	}

	deviceRepositoryImpl struct {
		db *gorm.DB
	}
)


func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepositoryImpl{
		db: db,
	}
}

func (repo *deviceRepositoryImpl) GetDevice(ctx context.Context, uuid uuid.UUID) (model.Device, error) {
	device := model.Device{}
	if err := repo.db.WithContext(ctx).Where("device_uuid = ?", uuid.String()).First(&device).Error; err != nil {
		return model.Device{}, err
	}
	return device, nil
}

func (repo *deviceRepositoryImpl) CreateDevice(ctx context.Context, device model.Device) (model.Device, error) {
	if err := repo.db.WithContext(ctx).Create(&device).Error; err != nil {
		return model.Device{}, err
	}
	return device, nil
}

func (repo *deviceRepositoryImpl) UpdateDevice(ctx context.Context, device model.Device) (model.Device, error) {
	if err := repo.db.WithContext(ctx).Updates(&device).Error; err != nil {
		return model.Device{}, err
	}
	return device, nil
}

func (repo *deviceRepositoryImpl) DeleteDevice(ctx context.Context, deviceUuid uuid.UUID) error {
	tx := repo.db.WithContext(ctx).Delete(&model.Device{}, "device_uuid=?", deviceUuid)
	if err := tx.Error; err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		panic(exception.NotFound{Message: "No device found with uuid " + deviceUuid.String()})
	}
	return nil
}

func (repo *deviceRepositoryImpl) SearchDevices(ctx context.Context, request model.SearchRequest) (model.Page[model.Device], error) {
	var devices []model.Device
	var err error
	var count int64
	
	dbTx := repo.db.WithContext(ctx)
	
	if request.DeviceName != "" {
    	dbTx = dbTx.Where("device_name LIKE ?", "%" + request.DeviceName + "%")
    }
    
    if request.BrandName != "" {
        dbTx = dbTx.Where("brand_name LIKE ?", "%" + request.BrandName + "%")
    }

	if err := dbTx.Model(&model.Device{}).Count(&count).Error; err != nil {
		return model.Page[model.Device]{
			PageNumber: request.PageNumber,
			PageSize:   request.PageSize,
			TotalPages: 0,
			TotalCount: count,
			Items:      devices,
			Sort:       request.Sort,
		}, err
	}
    offset := request.PageNumber * request.PageSize
	dbTx = dbTx.Model(&model.Device{}).Order(Order(mapSort(request.Sort))).Offset(offset).Limit(request.PageSize)

	if err := dbTx.Find(&devices).Error; err != nil {
		return model.Page[model.Device]{
			PageNumber: request.PageNumber,
			PageSize:   request.PageSize,
			TotalPages: 0,
			TotalCount: count,
			Items:      devices,
			Sort:       request.Sort,
		}, err
	}

	totalPage := int(math.Ceil(float64(count) / float64(request.PageSize)))

	return model.Page[model.Device]{
		PageNumber: request.PageNumber,
		PageSize:   request.PageSize,
		TotalPages: totalPage,
		TotalCount: count,
		Items:      devices,
		Sort:       request.Sort,
	}, err
}

func mapSort(sort model.Sort) model.Sort {
    sortFieldMapping := map[string] string {
        "name": "device_name",
        "created_at": "created_at",
    }
    return model.Sort{
    			SortBy:   sortFieldMapping[sort.SortBy],
    			Direction: sort.Direction,
    		}
}
