//go:build integration
// +build integration

package endpoints

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-redis/redismock/v9"
	"github.com/gorilla/mux"
	"github.com/maxgarvey/cq_server/admin"
	"github.com/maxgarvey/cq_server/handlers"
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
)

func setupIntegrationRouter() (*mux.Router, redismock.ClientMock, *admin.MockAdmin, *rabbitmq.FakeRabbitmq) {
	logger := slog.Default()
	db, mock := redismock.NewClientMock()
	mockedRedis := &redis.Redis{Client: *db}
	mockAdmin := admin.InitMock()
	fakeRabbitmq := &rabbitmq.FakeRabbitmq{}
	testClock := clock.NewMock()
	tokenFunc := func() string { return "test-token" }

	router := mux.NewRouter()
	router.HandleFunc("/get/{requestType}/{id}", Get(mockedRedis, logger)).Methods("GET")
	router.HandleFunc("/ask/{requestType}", Ask(testClock, fakeRabbitmq, mockedRedis, tokenFunc, logger)).Methods("POST")
	router.HandleFunc("/update/{requestType}/{id}", Update(testClock, mockedRedis, logger)).Methods("POST")
	router.HandleFunc("/health", Health(logger)).Methods("GET")
	router.HandleFunc("/admin/login", AdminLogin(&mockAdmin, *logger)).Methods("POST")
	router.HandleFunc("/admin/get/{requestType}/{id}", AdminGet(&mockAdmin, *mockedRedis, *logger)).Methods("GET")
	router.HandleFunc("/admin/ask/{requestType}", AdminAsk(&mockAdmin, testClock, fakeRabbitmq, mockedRedis, tokenFunc, logger)).Methods("POST")
	return router, mock, &mockAdmin, fakeRabbitmq
}

func TestIntegration_Get(t *testing.T) {
	router, mock, _, _ := setupIntegrationRouter()
	recorder := httptest.NewRecorder()
	// Expect Redis Get to be called with key "NOOP:test"
	record := `{"body":"test","id":"test","request_type":"NOOP","status":"DONE","timestamp":123}`
	mock.ExpectGet("NOOP:test").SetVal(record)
	request, _ := http.NewRequest("GET", "/get/noop/test", nil)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recorder.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestIntegration_Health(t *testing.T) {
	router, _, _, _ := setupIntegrationRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recorder.Code)
	}
	if recorder.Body.String() != "healthy" {
		t.Errorf("expected body 'healthy', got '%s'", recorder.Body.String())
	}
}

func TestIntegration_Ask(t *testing.T) {
	router, mock, _, fakeRabbitmq := setupIntegrationRouter()
	recorder := httptest.NewRecorder()
	body := `{"body":"test body"}`
	// Build expected record JSON
	expectedRecord := `{"body":"{\"body\":\"test body\"}","id":"test-token","request_type":0,"status":0,"timestamp":0}`
	mock.ExpectSet("NOOP:test-token", []byte(expectedRecord), 0).SetVal("OK")
	request, _ := http.NewRequest("POST", "/ask/noop", strings.NewReader(body))
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recorder.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
	// Verify RabbitMQ Publish called
	found := false
	for _, msg := range fakeRabbitmq.PublishedMessages {
		if msg == expectedRecord {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected RabbitMQ Publish to be called with correct record")
	}
}

func TestIntegration_Update(t *testing.T) {
	router, mock, _, _ := setupIntegrationRouter()
	recorder := httptest.NewRecorder()
	body := `{"status":"DONE"}`
	// Build expected record JSON for Get and Set
	getRecord := `{"body":"test","id":"test","request_type":0,"status":0,"timestamp":123}`
	setRecord := `{"body":"test","id":"test","request_type":0,"status":1,"timestamp":123}`
	mock.ExpectGet("NOOP:test").SetVal(getRecord)
	mock.ExpectSet("NOOP:test", []byte(setRecord), 0).SetVal("OK")
	request, _ := http.NewRequest("POST", "/update/noop/test", strings.NewReader(body))
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recorder.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestIntegration_AdminLogin(t *testing.T) {
	router, _, mockAdmin, _ := setupIntegrationRouter()
	recorder := httptest.NewRecorder()
	body := `{"username":"admin","password":"password"}`
	request, _ := http.NewRequest("POST", "/admin/login", strings.NewReader(body))
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recorder.Code)
	}
	if len(mockAdmin.LoginCalls) == 0 {
		t.Errorf("expected Login to be called on mockAdmin")
	}
}

func TestIntegration_AdminGet(t *testing.T) {
	router, mock, mockAdmin, _ := setupIntegrationRouter()
	recorder := httptest.NewRecorder()
	// Expect ValidateSession to be called and Redis Get
	getRecord := `{"body":"test","id":"test","request_type":0,"status":0,"timestamp":123}`
	mock.ExpectGet("NOOP:test").SetVal(getRecord)
	request, _ := http.NewRequest("GET", "/admin/get/noop/test", nil)
	request.Header.Set("SESSION", "session")
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recorder.Code)
	}
	if len(mockAdmin.ValidateSessionCalls) == 0 {
		t.Errorf("expected ValidateSession to be called on mockAdmin")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
}

func TestIntegration_AdminAsk(t *testing.T) {
	router, mock, mockAdmin, fakeRabbitmq := setupIntegrationRouter()
	recorder := httptest.NewRecorder()
	body := `{"body":"admin ask body"}`
	expectedRecord := `{"body":"{\"body\":\"admin ask body\"}","id":"test-token","request_type":0,"status":0,"timestamp":0}`
	mock.ExpectSet("NOOP:test-token", []byte(expectedRecord), 0).SetVal("OK")
	request, _ := http.NewRequest("POST", "/admin/ask/noop", strings.NewReader(body))
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recorder.Code)
	}
	if len(mockAdmin.ValidateSessionCalls) == 0 {
		t.Errorf("expected ValidateSession to be called on mockAdmin")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("redis expectations not met: %v", err)
	}
	found := false
	for _, msg := range fakeRabbitmq.PublishedMessages {
		if msg == expectedRecord {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected RabbitMQ Publish to be called with correct record")
	}
}

func TestIntegration_Download(t *testing.T) {
	router, mock, _, _ := setupIntegrationRouter()
	recorder := httptest.NewRecorder()

	// Setup a test HTTP server to serve a file
	testContent := "integration download test"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testContent))
	}))
	defer ts.Close()

	// Prepare download body as JSON using structs
	downloadBody := handlers.DownloadBody{
		Source:      ts.URL,
		Destination: "integration_download.txt",
	}
	downloadBodyBytes, _ := json.Marshal(downloadBody)

	expectedRecord := `{"body":"{\"source\":\"` + ts.URL + `\",\"destination\":\"integration_download.txt\"}","id":"test-token","request_type":2,"status":0,"timestamp":0}`
	mock.ExpectSet("RequestType(2):test-token", []byte(expectedRecord), 0).SetVal("OK")

	// POST to /ask/download
	request, _ := http.NewRequest(
		"POST", "/ask/download", strings.NewReader(string(downloadBodyBytes)))
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recorder.Code)
	}

	expectedRecord = `{"body":"{\"source\":\"` + ts.URL + `\",\"destination\":\"integration_download.txt\"}","id":"test-token","request_type":2,"status":0,"timestamp":0}`
	mock.ExpectGet("RequestType(2):test-token").SetVal(expectedRecord)

	// Poll for DONE status in the record
	var done bool
	var lastBody string
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		rec := httptest.NewRecorder()
		getReq, _ := http.NewRequest("GET", "/get/download/test-token", nil)
		router.ServeHTTP(rec, getReq)
		lastBody = rec.Body.String()
		if strings.Contains(lastBody, `"status":"DONE"`) || strings.Contains(lastBody, `"status":1`) {
			done = true
			break
		}
	}
	if !done {
		t.Errorf("expected record to be in DONE status, got: %s", lastBody)
	}

	// Check that the file was downloaded
	data, err := os.ReadFile("integration_download.txt")
	if err != nil {
		t.Fatalf("expected file to be created, got error: %v", err)
	}
	if string(data) != testContent {
		t.Errorf("expected file content %q, got %q", testContent, string(data))
	}

	// Clean up
}
