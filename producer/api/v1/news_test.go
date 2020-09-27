package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/usernamesalah/news/producer/pkg/models"
	"github.com/usernamesalah/news/producer/pkg/services/mocks"
)

type mockRequestValidator struct{}

func (m *mockRequestValidator) Validate(i interface{}) error {
	return nil
}

func TestAPI_listNews(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/news", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	mockNewsService := &mocks.NewsService{}
	mockNewsService.On("ListNews", mock.Anything).Return([]models.News{}, nil)

	api := NewAPI(mockNewsService, &mocks.KafkaService{}, &mocks.ElasticService{}, "", "")
	if assert.NoError(t, api.listNews(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[]\n", rec.Body.String())
	}
}

func TestAPI_createNews(t *testing.T) {
	news := models.News{
		Author: "news",
		Body:   "this is Description",
	}
	newsJSON, _ := json.Marshal(news)

	req := httptest.NewRequest(http.MethodPost, "/news", bytes.NewReader(newsJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	e.Validator = &mockRequestValidator{}
	c := e.NewContext(req, rec)

	mockNewsService := &mocks.NewsService{}
	mockNewsService.On("CreateNews", mock.Anything, news).Return(news, nil)

	api := NewAPI(mockNewsService, &mocks.KafkaService{}, &mocks.ElasticService{}, "", "")
	if assert.NoError(t, api.createNews(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "{\"id\":1,\"author\":\"news\",\"body\":\"this is Description\"}\n", rec.Body.String())
	}
}
