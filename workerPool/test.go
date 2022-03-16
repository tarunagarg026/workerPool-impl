package workerPool

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func executeRequest() *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestExecuteReq(t *testing.T) {
	pool := GetPoolConfig(5, 1, 10)
	var jsonStr = []byte(`{"name":"test", "delay": 1`)
	req, _ := http.NewRequest("POST", "localhost:8000/work", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest()
	parseReq(response, req, pool)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestExecuteReq_EmptyName(t *testing.T) {
	pool := GetPoolConfig(5, 1, 10)
	var jsonStr = []byte(`{"name":"", "delay": 1`)
	req, _ := http.NewRequest("POST", "localhost:8000/work", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest()
	parseReq(response, req, pool)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestExecuteReq_WrongDelay(t *testing.T) {
	pool := GetPoolConfig(5, 1, 10)
	var jsonStr = []byte(`{"name":"", "delay": 20`)
	req, _ := http.NewRequest("POST", "localhost:8000/work", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest()
	parseReq(response, req, pool)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}
