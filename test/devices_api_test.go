package test

import (
	"bytes"
	"device-management/config"
	"device-management/model"
	"device-management/router"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

var nGoRoutine int
var Server *gin.Engine
var run func(func(*testing.T)) func(*testing.T) = SetBeforeAndAfterEach(
	func(t *testing.T) {
		t.Logf("Before %s", t.Name())
		nGoRoutine = runtime.NumGoroutine()
	},
	func(t *testing.T) {
		t.Logf("After %s (goroutines: %d => %d)", t.Name(), nGoRoutine, runtime.NumGoroutine())
	},
)

func TestMain(m *testing.M) {
	log.Println("Before All")
	err := os.Setenv("PROFILE", "test")
	if err != nil {
		log.Fatalf("Error setting PROFILE env var: %v", err)
	}
	err = os.Setenv(gin.EnvGinMode, "test")
	if err != nil {
		log.Fatalf("Error setting %s env var: %v", gin.EnvGinMode, err)
	}

	e := exec.Command("pwd")
	var out bytes.Buffer
	e.Stdout = &out
	err = e.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Output: %q\n", out.String())
	if err != nil {
		log.Fatalf("Error while intialize the test %v", err)
	}

	gin.SetMode(gin.TestMode)
	Server = router.NewRouter(config.Controllers...)
	m.Run()
	log.Println("After All")
	os.Exit(0)
}

func SetBeforeAndAfterEach(beforeFunc, afterFunc func(*testing.T)) func(func(*testing.T)) func(*testing.T) {
	return func(test func(*testing.T)) func(*testing.T) {
		return func(t *testing.T) {
			if beforeFunc != nil {
				beforeFunc(t)
			}
			test(t)
			if afterFunc != nil {
				afterFunc(t)
			}
		}
	}
}

func TestCreatDeviceFailed(t *testing.T) {
	log.Println("start test")
	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "d"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestCreatDeviceSuccessfully(t *testing.T) {
	log.Println("start test")
	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "Device1", "brand":"Brand1"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusCreated, w.Code)
	device := model.Device{}
	json.Unmarshal([]byte(w.Body.String()), &device)
	assert.EqualValues(t, device.BrandName, "Brand")
	assert.EqualValues(t, device.Name, "Device")
}

func TestCreatDeviceDuplicateDeviceFailed(t *testing.T) {
	log.Println("start test")
	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "Device2", "brand":"Brand2"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	req, err = http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "Device", "brand":"Brand"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
	apiError := model.ApiError{}
	json.Unmarshal([]byte(w.Body.String()), &apiError)
	assert.EqualValues(t, apiError.Code, "400")
}
