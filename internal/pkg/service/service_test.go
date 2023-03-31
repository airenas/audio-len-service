package service_test

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/airenas/audio-len-service/internal/pkg/service"
	"github.com/stretchr/testify/assert"
)

var data *service.Data
var saver *testSaver
var estimator *testEstimator

func initTest(t *testing.T) {
	saver = &testSaver{name: "test.wav"}
	estimator = &testEstimator{res: 2}
	data = newTestData(saver, estimator)
}

func TestNotFound(t *testing.T) {
	initTest(t)
	req, err := http.NewRequest("GET", "/any", nil)
	assert.Nil(t, err)
	resp := httptest.NewRecorder()

	service.NewRouter(data).ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestNotFound_Method(t *testing.T) {
	initTest(t)
	req, err := http.NewRequest("GET", "/duration", nil)
	assert.Nil(t, err)
	resp := httptest.NewRecorder()

	service.NewRouter(data).ServeHTTP(resp, req)
	assert.Equal(t, 405, resp.Code)
}

func TestReturns(t *testing.T) {
	initTest(t)
	req := newTestRequest("test.mp3")
	resp := httptest.NewRecorder()

	service.NewRouter(data).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `{"duration":2}`, strings.TrimSpace(resp.Body.String()))
	assert.Equal(t, "body", string(saver.data.Bytes()))
	assert.Equal(t, "test.wav", estimator.name)
}

func TestReturns_UppercaseFile(t *testing.T) {
	initTest(t)
	req := newTestRequest("test.M4A")
	resp := httptest.NewRecorder()

	service.NewRouter(data).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `{"duration":2}`, strings.TrimSpace(resp.Body.String()))
}

func TestFails_NoFile(t *testing.T) {
	initTest(t)
	req := newTestRequest("")
	resp := httptest.NewRecorder()

	service.NewRouter(data).ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}

func TestFails_Saver(t *testing.T) {
	initTest(t)
	saver.err = errors.New("olia")
	req := newTestRequest("test.wav")
	resp := httptest.NewRecorder()

	service.NewRouter(data).ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}

func TestFails_Estimator(t *testing.T) {
	initTest(t)
	estimator.err = errors.New("olia")
	req := newTestRequest("test.wav")
	resp := httptest.NewRecorder()

	service.NewRouter(data).ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}



func newTestRequest(file string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if file != "" {
		part, _ := writer.CreateFormFile("file", file)
		_, _ = io.Copy(part, strings.NewReader("body"))
	}
	writer.Close()
	req := httptest.NewRequest("POST", "/duration", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

type testSaver struct {
	name string
	err  error
	data bytes.Buffer
}

func (s *testSaver) Save(name string, reader io.Reader) (string, error) {
	io.Copy(&s.data, reader)
	return s.name, s.err
}

type testEstimator struct {
	res  float64
	err  error
	name string
}

func (s *testEstimator) Seconds(name string) (float64, error) {
	s.name = name
	return s.res, s.err
}

func newTestData(s service.FileSaver, e service.DurationEstimator) *service.Data {
	return &service.Data{Saver: s, Estimator: e}
}
