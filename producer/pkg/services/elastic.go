package services

import (
	"context"
	"encoding/json"

	"github.com/usernamesalah/news/producer/pkg/models"

	"github.com/olivere/elastic"
)

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"news":{
			"properties":{
				"id":{
					"type":"number"
				},
				"author":{
					"type":"keyword"
				},
				"body":{
					"type":"keyword"
				},
				"created":{
					"type":"date"
				}
			}
		}
	}
}`

// ElasticService service interface.
type ElasticService interface {
	CheckExist(ctx context.Context) error
	Get(ctx context.Context, rows, page int) ([]models.News, error)
	Create(ctx context.Context, news models.News) (models.News, error)
	GetByID(ctx context.Context, id int64) (models.News, error)
}

// elasticService hold kafka producer session
type elasticService struct {
	client *elastic.Client
	index  string
}

// NewElasticService returns an initialized ElasticService implementation.
func NewElasticService(client *elastic.Client, index string) ElasticService {
	return &elasticService{client: client, index: index}
}

// CheckExist index Elastic
func (e *elasticService) CheckExist(ctx context.Context) error {
	exists, err := e.client.IndexExists(e.index).Do(ctx)
	if err != nil {
		// Handle error
		return err
	}
	if !exists {
		// Create a new index.
		createIndex, err := e.client.CreateIndex(e.index).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return err
		}
	}

	return nil
}

func (e *elasticService) Get(ctx context.Context, rows, page int) ([]models.News, error) {

	if err := e.CheckExist(ctx); err != nil {
		return []models.News{}, err
	}

	var news []models.News
	if rows <= 0 {
		rows = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * rows

	aggr := elastic.NewTermsAggregation().Field("id").Size(rows)
	searchResult, err := e.client.Search().
		Index(e.index).         // search in index
		Sort("created", false). // sort by "user" field, ascending
		From(offset).Size(rows).
		Aggregation("content_id", aggr).
		Pretty(true). // pretty print request and response JSON
		Do(ctx)       // execute
	if err != nil {
		return []models.News{}, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var cmp models.News
		_ = json.Unmarshal(*hit.Source, &cmp)

		news = append(news, cmp)
	}

	return news, nil
}

func (e *elasticService) Create(ctx context.Context, news models.News) (models.News, error) {

	if err := e.CheckExist(ctx); err != nil {
		return models.News{}, err
	}

	_, err := e.client.Index().
		Index(e.index).
		BodyJson(news).
		Do(ctx)
	if err != nil {
		return models.News{}, err
	}

	return news, nil
}

func (e *elasticService) GetByID(ctx context.Context, id int64) (models.News, error) {

	if err := e.CheckExist(ctx); err != nil {
		return models.News{}, err
	}

	var news models.News

	data, err := e.client.Get().
		Index(e.index).
		Id(string(id)).
		Do(ctx)
	if err != nil {
		return models.News{}, err
	}

	_ = json.Unmarshal(*data.Source, &news)

	return news, nil
}
