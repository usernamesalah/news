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

	"github.com/usernamesalah/news/producer/pkg/model"
	"github.com/usernamesalah/news/producer/pkg/services/mock"
)

type mockRequestValidator struct{}

func (m *mockRequestValidator) Validate(i interface{}) error {
	return nil
}

func TestAPI_listTeams(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/teams", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	mockTeamsService := &mocks.TeamsService{}
	mockTeamsService.On("ListTeams", mock.Anything).Return([]models.Team{}, nil)

	api := NewAPI(mockTeamsService, &mocks.PlayersService{}, "", "")
	if assert.NoError(t, api.listTeams(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[]\n", rec.Body.String())
	}
}

func TestAPI_getTeam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/teams/1", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/teams/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockTeamsService := &mocks.TeamsService{}
	mockTeamsService.On("GetTeam", mock.Anything, int64(1)).Return(models.Team{}, nil)

	api := NewAPI(mockTeamsService, &mocks.PlayersService{}, "", "")
	if assert.NoError(t, api.getTeam(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{\"id\":0,\"name\":\"\",\"description\":\"\"}\n", rec.Body.String())
	}
}

func TestAPI_createTeam(t *testing.T) {
	team := models.Team{
		ID:          1,
		Name:        "team-1",
		Description: "this is Description",
	}
	teamJSON, _ := json.Marshal(team)

	req := httptest.NewRequest(http.MethodPost, "/teams", bytes.NewReader(teamJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	e.Validator = &mockRequestValidator{}
	c := e.NewContext(req, rec)

	mockTeamsService := &mocks.TeamsService{}
	mockTeamsService.On("CreateTeam", mock.Anything, team).Return(team, nil)

	api := NewAPI(mockTeamsService, &mocks.PlayersService{}, "", "")
	if assert.NoError(t, api.createTeam(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "{\"id\":1,\"name\":\"team-1\",\"description\":\"this is Description\"}\n", rec.Body.String())
	}
}

func TestAPI_deleteTeam(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/teams/1", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/teams/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockTeamsService := &mocks.TeamsService{}
	mockTeamsService.On("DeleteTeam", mock.Anything, int64(1)).Return(nil)

	api := NewAPI(mockTeamsService, &mocks.PlayersService{}, "", "")
	if assert.NoError(t, api.deleteTeam(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Equal(t, "", rec.Body.String())
	}
}

func TestAPI_updateTeam(t *testing.T) {
	team := models.Team{
		Name:        "team-update-1",
		Description: "Description",
	}
	teamJSON, _ := json.Marshal(team)

	req := httptest.NewRequest(http.MethodPut, "/teams/1", bytes.NewReader(teamJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	e.Validator = &mockRequestValidator{}
	c := e.NewContext(req, rec)

	mockTeamsService := &mocks.TeamsService{}
	mockTeamsService.On("UpdateTeam", mock.Anything, team).Return(team, nil)

	api := NewAPI(mockTeamsService, &mocks.PlayersService{}, "", "")
	if assert.NoError(t, api.updateTeam(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "{\"id\":0,\"name\":\"team-update-1\",\"description\":\"Description\"}\n", rec.Body.String())
	}
}
