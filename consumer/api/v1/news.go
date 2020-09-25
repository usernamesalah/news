package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/usernamesalah/news/producer/pkg/model"
)

// List news
// @Summary List news
// @Description Get the list of news
// @Tags news
// @ID list-news
// @Produce json
// @Success 200 {array} models.News
// @Router /news [get]
func (api *API) listNews(c echo.Context) error {
	ctx := c.Request().Context()

	news, err := api.newsService.ListNews(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, news)
}

// Get an news
// @Summary Get an news
// @Description Get an news by id
// @Tags news
// @ID get-new
// @Produce json
// @Param id path int true "News ID"
// @Success 200 {object} models.News
// @Router /news/{id} [get]
func (api *API) getNews(c echo.Context) error {
	ctx := c.Request().Context()

	idString := c.Param("id")
	id, _ := strconv.ParseInt(idString, 10, 64)

	new, err := api.newsService.GetNews(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, new)
}

// Create a new new
// @Summary Create a new news
// @Description Create a new news
// @Tags news
// @ID create-new
// @Produce json
// @Param new body models.News true "Create new"
// @Success 201 {object} models.News
// @Router /news [post]
func (api *API) createNews(c echo.Context) error {

	new := new(models.News)
	if err := c.Bind(new); err != nil {
		return err
	}

	if err := c.Validate(new); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	data, _ := json.Marshal(new)
	err := api.kafkaService.SendMessage("news_created", string(data))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, new)
}
