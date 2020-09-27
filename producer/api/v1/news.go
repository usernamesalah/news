package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/usernamesalah/news/producer/pkg/models"
)

// List news
// @Summary List news
// @Description Get the list of news
// @Tags news
// @ID list-news
// @Produce json
// @Param page query integer false "page number" default(1)
// @Param rows query integer false "rows" default(10)
// @Success 200 {array} models.News
// @Router /news [get]
func (api *API) listNews(c echo.Context) error {

	ctx := c.Request().Context()
	page, _ := strconv.Atoi(c.QueryParam("page"))
	rows, _ := strconv.Atoi(c.QueryParam("rows"))

	news, err := api.elasticService.Get(ctx, rows, page)
	if err != nil {
		return err
	}

	data := []models.News{}

	go func() {
		for i := 0; i < len(news); i++ {
			data[i], _ = api.newsService.GetNews(ctx, news[i].ID)
		}
	}()

	return c.JSON(http.StatusOK, data)
}

// Create a new news
// @Summary Create a new news
// @Description Create a new news
// @Tags news
// @ID create-new
// @Produce json
// @Param new body models.News true "Create news"
// @Success 201 {object} models.News
// @Router /news [post]
func (api *API) createNews(c echo.Context) error {

	news := new(models.News)
	if err := c.Bind(news); err != nil {
		return err
	}

	if err := c.Validate(news); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	data, _ := json.Marshal(news)

	err := api.kafkaService.SendMessage(string(data))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, news)
}
