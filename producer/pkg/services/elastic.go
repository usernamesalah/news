package services

import (
	"context"
	"encoding/json"
	"newz/models"

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
	CheckExist(ctx context.Context, index string) error
	Get(ctx context.Context, index string, rows, page int) ([]models.News, error)
}

// elasticService hold kafka producer session
type elasticService struct {
	client *elastic.Client
}

// NewElasticService returns an initialized ElasticService implementation.
func NewElasticService(client *elastic.Client) ElasticService {
	return &elasticService{client: client}
}

// SendMessage function to send message into kafka
func (e *elasticService) CheckExist(ctx context.Context, index string) error {
	exists, err := e.client.IndexExists(index).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := e.client.CreateIndex(index).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return err
		}
	}
	return nil
}

func (e *elasticService) Get(ctx context.Context, index string, rows, page int) ([]models.News, error) {

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
		Index(index).           // search in index
		Sort("created", false). // sort by "user" field, ascending
		From(offset).Size(rows).
		Aggregation("content_id", aggr).
		Pretty(true). // pretty print request and response JSON
		Do(ctx)       // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		var cmp models.News
		_ = json.Unmarshal(*hit.Source, &cmp)

		news = append(news, cmp)
	}

	return news, nil
}
