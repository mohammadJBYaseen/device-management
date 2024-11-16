package controller

import (
	"context"
	"device-management/exception"
	"device-management/model"
	"device-management/router"
	"device-management/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"slices"
	"strconv"
)

var (
	supportedSortBy        = []string{"creation_time", "name"}
	supportedSortDirection = []string{"ASC", "DESC"}
)

type (
	DeviceController interface {
		CreateDevice(ctx context.Context, device model.Device) model.Device
		GetDeviceByUuid(ctx context.Context, deviceUuid uuid.UUID) model.Device
		UpdateDeviceByUuid(ctx context.Context, deviceUuid uuid.UUID, device model.Device) model.Device
		PatchDeviceByUuid(ctx context.Context, deviceUuid uuid.UUID, device []model.JsonPatch) model.Device
		DeleteDeviceByUuid(ctx context.Context, deviceUuid uuid.UUID)
		GetDevices(ctx context.Context, deviceName string, brandName string, pageNumber int, pageSize int, sort model.Sort) model.Page[model.Device]
		Routes() []router.Route
		Group() string
	}

	deviceControllerImpl struct {
		deviceService service.DeviceService
	}
)

func NewDeviceController(deviceService service.DeviceService) DeviceController {
	return &deviceControllerImpl{
		deviceService: deviceService,
	}
}

func (deviceController *deviceControllerImpl) CreateDevice(ctx context.Context, device model.Device) model.Device {
	creatDevice, err := deviceController.deviceService.CreatDevice(ctx, device)
	if err != nil {
		panic(err)
	}
	return creatDevice
}

func (deviceController *deviceControllerImpl) GetDeviceByUuid(ctx context.Context, deviceUuid uuid.UUID) model.Device {
	device, err := deviceController.deviceService.GetDevice(ctx, deviceUuid)
	if err != nil {
		panic(err)
	}
	return device
}

func (deviceController *deviceControllerImpl) UpdateDeviceByUuid(ctx context.Context, deviceUuid uuid.UUID, device model.Device) model.Device {
	updatedDevice, err := deviceController.deviceService.UpdateDevice(ctx, deviceUuid, device)
	if err != nil {
		panic(err)
	}
	return updatedDevice
}

func (deviceController *deviceControllerImpl) DeleteDeviceByUuid(ctx context.Context, deviceUuid uuid.UUID) {
	err := deviceController.deviceService.DeleteDevice(ctx, deviceUuid)
	if err != nil {
		panic(err)
	}
}

func (deviceController *deviceControllerImpl) GetDevices(ctx context.Context, deviceName string, brandName string, pageNumber int, pageSize int,
	sort model.Sort) model.Page[model.Device] {
	devicesPage, err := deviceController.deviceService.SearchDevices(ctx, model.SearchRequest{
		DeviceName: deviceName,
		BrandName:  brandName,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Sort:       sort,
	})
	if err != nil {
		panic(err)
	}
	return devicesPage
}

func (deviceController *deviceControllerImpl) PatchDeviceByUuid(ctx context.Context, deviceUuid uuid.UUID, jsonPath []model.JsonPatch) model.Device {
	updatedDevice, err := deviceController.deviceService.PatchDevice(ctx, deviceUuid, jsonPath)
	if err != nil {
		panic(err)
	}
	return updatedDevice
}

// handler

func (deviceController *deviceControllerImpl) createDeviceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.ContentType()
		var device model.Device
		err := ctx.Bind(&device)
		if err != nil {
			panic(err)
		}
		ctx.JSON(http.StatusCreated, deviceController.CreateDevice(ctx, device))
		ctx.Next()
	}
}

func (deviceController *deviceControllerImpl) getDeviceByUuidHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceUuid, err := uuid.Parse(ctx.Param("device-uuid"))
		if err != nil {
			panic(err)
		}
		ctx.JSON(http.StatusOK, deviceController.GetDeviceByUuid(ctx, deviceUuid))
		ctx.Next()
	}
}

func (deviceController *deviceControllerImpl) deleteDeviceByUuidHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceUuid, err := uuid.Parse(ctx.Param("device-uuid"))
		if err != nil {
			panic(err)
		}
		deviceController.DeleteDeviceByUuid(ctx, deviceUuid)
		ctx.JSON(http.StatusNoContent, "{}")
		ctx.Next()
	}
}

func (deviceController *deviceControllerImpl) updateDeviceByUuidHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceUuid, err := uuid.Parse(ctx.Param("device-uuid"))
		if err != nil {
			panic(err)
		}
		var device model.Device
		err = ctx.Bind(&device)
		if err != nil {
			panic(err)
		}
		ctx.JSON(http.StatusAccepted, deviceController.UpdateDeviceByUuid(ctx, deviceUuid, device))
		ctx.Next()
	}
}

func (deviceController *deviceControllerImpl) patchDeviceByUuidHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceUuid, err := uuid.Parse(ctx.Param("device-uuid"))
		if err != nil {
			panic(err)
		}
		var patchRequest []model.JsonPatch
		err = ctx.Bind(&patchRequest)
		if err != nil {
			panic(err)
		}
		ctx.JSON(http.StatusAccepted, deviceController.PatchDeviceByUuid(ctx, deviceUuid, patchRequest))
		ctx.Next()
	}
}

func (deviceController *deviceControllerImpl) getDevicesHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceName := ctx.Query("device_name")
		brandName := ctx.Query("brand_name")
		sortBy := ctx.DefaultQuery("sort_by", "created_at")
		if !slices.Contains(supportedSortBy, sortBy) {
			panic(exception.BadRequest{Message: fmt.Sprintf("unsupported sort_by, value must be one of: [created_at, name], provided: %s", sortBy)})
		}
		sortDirection := ctx.DefaultQuery("sort_dir", "DESC")
		if !slices.Contains(supportedSortDirection, sortDirection) {
			panic(exception.BadRequest{Message: fmt.Sprintf("unsupported sort_dir, value must be one of: [ASC, DESC], provided: %s", sortDirection)})
		}
		pageNumber, err := strconv.Atoi(ctx.DefaultQuery("page_number", "0"))
		if err != nil {
			panic(err)
		}
		pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
		if err != nil {
			panic(err)
		}
		ctx.JSON(http.StatusOK, deviceController.GetDevices(ctx, deviceName, brandName, pageNumber, pageSize, model.Sort{
			SortBy:    sortBy,
			Direction: sortDirection,
		}))
		ctx.Next()
	}
}

func (deviceController *deviceControllerImpl) Routes() []router.Route {
	return []router.Route{
		{
			"CreateDevice",
			http.MethodPost,
			"/devices",
			"application/json",
			deviceController.createDeviceHandler(),
		},
		{
			"GetDeviceByUuid",
			http.MethodGet,
			"/devices/:device-uuid",
			"",
			deviceController.getDeviceByUuidHandler(),
		},
		{
			"DeleteDeviceByUuid",
			http.MethodDelete,
			"/devices/:device-uuid",
			"",
			deviceController.deleteDeviceByUuidHandler(),
		},
		{
			"UpdateDeviceByUuid",
			http.MethodPut,
			"/devices/:device-uuid",
			"application/json",
			deviceController.updateDeviceByUuidHandler(),
		},
		{
			"PatchDeviceByUuid",
			http.MethodPatch,
			"/devices/:device-uuid",
			"application/json",
			deviceController.patchDeviceByUuidHandler(),
		},
		{
			"GetDevices",
			http.MethodGet,
			"/devices",
			"application/json",
			deviceController.getDevicesHandler(),
		},
	}
}

func (deviceController *deviceControllerImpl) Group() string {
	return "devices"
}
