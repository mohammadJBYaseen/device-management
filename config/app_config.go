package config

import (
	"device-management/controller"
	"device-management/repository"
	"device-management/router"
	"device-management/service"
	"errors"
	"fmt"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"gorm.io/gorm"
	"os"
	"time"
)

var ApplicationProperties AppProperties
var DB *gorm.DB
var DeviceRepo repository.DeviceRepository
var DeviceService service.DeviceService
var Controllers []router.BaseController

func init() {
	ApplicationProperties = Load()
	DB = initDBConnection()
	DeviceRepo = repository.NewDeviceRepository(DB)
	DeviceService = service.NewDeviceService(DeviceRepo)
	Controllers = []router.BaseController{controller.NewDeviceController(DeviceService)}
}

func Unload() {
	closeDBConnection()
}

type (
	AppProperties struct {
		Server   ServerProperties   `json:"server"`
		Database DatabaseProperties `json:"database"`
	}
	ServerProperties struct {
		AppName string `json:"app_name" env:"APP_NAME"`
		Port    int    `json:"port" env:"SERVER_PORT"`
	}

	DatabaseProperties struct {
		Host            string `json:"host" env:"DATABASE_HOST"`
		Port            int    `json:"port" env:"DATABASE_PORT"`
		Username        string `json:"user_name" env:"DATABASE_USERNAME"`
		Password        string `json:"password" env:"DATABASE_PASSWORD"`
		Database        string `json:"database" env:"DATABASE_NAME"`
		Config          string `json:"config"`
		Prefix          string `json:"prefix"`
		MaxIdleConn     int    `json:"max_idle_conn" env:"DATABASE_MAX_IDLE_CONN"`
		ConnMaxIdleTime string `json:"conn_max_idle_time" env:"DATABASE_CONN_MAX_IDLE_TIME"`
		MaxOpenConn     int    `json:"max_open_conn" env:"DATABASE_MAX_OPEN_CONN"`
		MinOpenConn     int    `json:"min_open_conn" env:"DATABASE_MIN_OPEN_CONN"`
		SSLMode         string `json:"ssl_mode" env:"DATABASE_SSL_MODE"`
	}
)

// SetUp Post config initialization
func (appConf *AppProperties) SetUp() error {
	return nil
}

// Load config initialization
func Load() AppProperties {
	env, exists := os.LookupEnv("PROFILE")
	if !exists {
		env = "dev"
	}
	cfg := AppProperties{}
	jsonFeeder := feeder.Json{Path: "config/config.json"}
	feederChain := config.New().AddFeeder(jsonFeeder)
	envConfigPath := fmt.Sprintf("config/config-%s.json", env)
	exists, err := fileExist(envConfigPath)
	if err != nil {
		panic(err)
	} else if exists {
		envJsonFeeder := feeder.Json{Path: envConfigPath}
		feederChain = feederChain.AddFeeder(envJsonFeeder)
	}

	feederChain = feederChain.AddFeeder(&feeder.Env{})

	if err = feederChain.AddStruct(&cfg).Feed(); err != nil {
		panic(err)
	}
	return cfg
}

func (c *DatabaseProperties) GetConnMaxIdleTime() time.Duration {

	d, err := time.ParseDuration(c.ConnMaxIdleTime)
	if err != nil {
		panic(err)
	}
	return d
}

func fileExist(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
