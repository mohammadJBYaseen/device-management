package controller

import (
	"device-management/model"
	"device-management/router"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type (
	DeviceController interface {
		CreateDevice(device model.Device) model.Device
		GetDeviceByUuid(deviceUuid uuid.UUID) model.Device
		UpdateDeviceByUuid(deviceUuid uuid.UUID, device model.Device) model.Device
		DeleteDeviceByUuid(deviceUuid uuid.UUID)
		GetDevices(deviceName string, brandName string, pageNumber int, pageSize int, sort model.Sort) model.Page[model.Device]
		Routes() []router.Route
		Group() string
	}

	deviceControllerImpl struct {
	}
)

func NewDeviceController() DeviceController {
	return &deviceControllerImpl{}
}

func (deviceController *deviceControllerImpl) CreateDevice(device model.Device) model.Device {
	//TODO implement me
	panic("implement me")
}

func (deviceController *deviceControllerImpl) GetDeviceByUuid(deviceUuid uuid.UUID) model.Device {
	//TODO implement me
	panic("implement me")
}

func (deviceController *deviceControllerImpl) UpdateDeviceByUuid(deviceUuid uuid.UUID, device model.Device) model.Device {
	//TODO implement me
	panic("implement me")
}

func (deviceController *deviceControllerImpl) DeleteDeviceByUuid(deviceUuid uuid.UUID) {
	//TODO implement me
	panic("implement me")
}

func (deviceController *deviceControllerImpl) GetDevices(deviceName string, brandName string, pageNumber int, pageSize int,
	sort model.Sort) model.Page[model.Device] {
	//TODO implement me
	panic("implement me")
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
		err = json.NewEncoder(ctx.Writer).Encode(deviceController.CreateDevice(device))
		if err != nil {
			panic(err)
		}
		ctx.Status(http.StatusCreated)
	}
}

func (deviceController *deviceControllerImpl) getDeviceByUuidHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceUuid, err := uuid.Parse(ctx.Param("device-uuid"))
		if err != nil {
			panic(err)
		}
		err = json.NewEncoder(ctx.Writer).Encode(deviceController.GetDeviceByUuid(deviceUuid))
		if err != nil {
			panic(err)
		}
		ctx.Status(http.StatusOK)
	}
}

func (deviceController *deviceControllerImpl) deleteDeviceByUuidHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceUuid, err := uuid.Parse(ctx.Param("device-uuid"))
		if err != nil {
			panic(err)
		}
		deviceController.DeleteDeviceByUuid(deviceUuid)
		ctx.String(http.StatusNoContent, "{}")
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
		err = json.NewEncoder(ctx.Writer).Encode(deviceController.UpdateDeviceByUuid(deviceUuid, device))
		if err != nil {
			panic(err)
		}
		ctx.Status(http.StatusAccepted)
	}
}

func (deviceController *deviceControllerImpl) getDevicesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func (deviceController *deviceControllerImpl) Routes() []router.Route {
	return []router.Route{
		{
			"CreateDevice",
			http.MethodPost,
			"/devices",
			"application/json",
			deviceController.getDevicesHandler(),
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
			deviceController.getDeviceByUuidHandler(),
		},
		{
			"UpdateDeviceByUuid",
			http.MethodPut,
			"/devices/:device-uuid",
			"application/json",
			deviceController.updateDeviceByUuidHandler(),
		},
		{
			"GetDevices",
			http.MethodGet,
			"/devices",
			"",
			deviceController.getDevicesHandler(),
		},
	}
}

func (deviceController *deviceControllerImpl) Group() string {
	return "devices"
}
