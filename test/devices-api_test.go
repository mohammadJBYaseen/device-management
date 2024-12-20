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

	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "d"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, w.Code, http.StatusBadRequest)
}

func TestCreatDeviceSuccessfully(t *testing.T) {

	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "Device1", "brand":"Brand1"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, w.Code, http.StatusCreated)
	device := model.Device{}
	json.Unmarshal([]byte(w.Body.String()), &device)
	assert.EqualValues(t, device.BrandName, "Brand1")
	assert.EqualValues(t, device.Name, "Device1")
}

func TestCreatDeviceDuplicateDeviceFailed(t *testing.T) {

	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "Device2", "brand":"Brand2"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	req, err = http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "Device2", "brand":"Brand1"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, w.Code, http.StatusConflict)
	apiError := model.ApiError{}
	json.Unmarshal([]byte(w.Body.String()), &apiError)
	assert.EqualValues(t, apiError.Code, "409")
}

func TestUpdateDeviceSuccessfully(t *testing.T) {
	log.Println("start test")
	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "Device3", "brand":"Brand3"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, w.Code, http.StatusCreated)
	device := model.Device{}
	json.Unmarshal([]byte(w.Body.String()), &device)
	assert.EqualValues(t, device.BrandName, "Brand3")

	w = httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	req, err = http.NewRequest("PUT", "/devices/"+device.UUID.String(), bytes.NewBuffer([]byte(`{"name": "Device3", "brand":"Brand4"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, w.Code, http.StatusAccepted)
	device = model.Device{}
	json.Unmarshal([]byte(w.Body.String()), &device)
	assert.EqualValues(t, device.BrandName, "Brand4")
}

func TestUpdateFailedSuccessfully(t *testing.T) {
	log.Println("start test")
	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "Device4", "brand":"Brand4"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, w.Code, http.StatusCreated)
	device := model.Device{}
	json.Unmarshal([]byte(w.Body.String()), &device)

	w = httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	req, err = http.NewRequest("PUT", "/devices/"+device.UUID.String(), bytes.NewBuffer([]byte(`{"name": "Device3"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, w.Code, http.StatusBadRequest)
}

func TestDeleteDeviceSuccessfully(t *testing.T) {
	log.Println("start test")
	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte(`{"name": "Device5", "brand":"Brand3"}`)))
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, w.Code, http.StatusCreated)
	device := model.Device{}
	json.Unmarshal([]byte(w.Body.String()), &device)
	assert.EqualValues(t, device.Name, "Device5")

	w = httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	req, err = http.NewRequest("DELETE", "/devices/"+device.UUID.String(), nil)
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}
	w = httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	assert.EqualValues(t, w.Code, http.StatusNoContent)

	w = httptest.NewRecorder()
	Server.ServeHTTP(w, req)

	req, err = http.NewRequest("GET", "/devices/"+device.UUID.String(), nil)
	if err != nil {
		t.Errorf("Error in creating request: %v", err)
	}

	w = httptest.NewRecorder()
	Server.ServeHTTP(w, req)
	assert.EqualValues(t, w.Code, http.StatusNotFound)
}
