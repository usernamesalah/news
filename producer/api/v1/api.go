package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/usernamesalah/news/producer/pkg/services"
)

// API can register a set of endpoints in a router and handle
// them using the provided storage.
type API struct {
	newsService    services.NewsService
	kafkaService   services.KafkaService
	elasticService services.ElasticService

	adminUsername string
	adminPassword string
}

// NewAPI returns an initialized API type.
func NewAPI(newsService services.NewsService, kafkaService services.KafkaService, elasticService services.ElasticService, adminUsername, adminPassword string) *API {
	return &API{
		newsService:    newsService,
		kafkaService:   kafkaService,
		elasticService: elasticService,

		adminUsername: adminUsername,
		adminPassword: adminPassword,
	}
}

// Register the API's endpoints in the given router.
func (api *API) Register(g *echo.Group) {
	// News API
	g.GET("/news", api.listNews)
	g.POST("/news", api.createNews, middleware.BasicAuth(api.adminValidator))

}

func (api *API) adminValidator(username, password string, c echo.Context) (bool, error) {
	if username == api.adminUsername && password == api.adminPassword {
		return true, nil
	}
	return false, nil
}
