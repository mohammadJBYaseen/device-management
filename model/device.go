package model

import "time"

type Device struct {
	Uuid string `json:"uuid,omitempty"`

	Name string `json:"name"`

	Brand string `json:"brand,omitempty"`

	CreationTime time.Time `json:"creation_time,omitempty"`
}
