package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/usernamesalah/news/producer/pkg/model"
)

// NewsService service interface.
type NewsService interface {
	ListNews(ctx context.Context) ([]models.News, error)
	GetNews(ctx context.Context, id int64) (models.News, error)
	CreateNews(ctx context.Context, news models.News) (models.News, error)
	DeleteNews(ctx context.Context, id int64) error
	UpdateNews(ctx context.Context, news models.News) (models.News, error)
}

type newsService struct {
	db *sqlx.DB
}

// NewNewsService returns an initialized NewsService implementation.
func NewNewsService(db *sqlx.DB) NewsService {
	return &newsService{db: db}
}

func (s *newsService) ListNews(ctx context.Context) ([]models.News, error) {
	query := `
		SELECT
			id
			, author
			, body
			, created
		FROM news`

	var news []models.News
	if err := s.db.SelectContext(ctx, &news, query); err != nil {
		return nil, fmt.Errorf("get the list of news: %s", err)
	}

	return news, nil
}

func (s *newsService) GetNews(ctx context.Context, id int64) (models.News, error) {
	query := `
		SELECT
			id
			, author
			, body
			, created
		FROM news
		WHERE id = $1`

	var news models.News
	if err := s.db.GetContext(ctx, &news, query, id); err != nil {
		return models.News{}, fmt.Errorf("get an news: %s", err)
	}

	return news, nil
}

func (s *newsService) CreateNews(ctx context.Context, news models.News) (models.News, error) {
	query := "INSERT INTO news (author, body) VALUES ($1, $2) RETURNING id"

	var id int64
	if err := s.db.QueryRowxContext(ctx, query, news.Author, news.Body).Scan(&id); err != nil {
		return models.News{}, fmt.Errorf("insert new news: %s", err)
	}

	newNews, err := s.GetNews(ctx, id)
	if err != nil {
		return models.News{}, fmt.Errorf("get new news: %s", err)
	}

	return newNews, nil
}

func (s *newsService) DeleteNews(ctx context.Context, id int64) error {
	query := `DELETE FROM news WHERE id = $1`

	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("delete an news: %s", err)
	}

	return nil
}

func (s *newsService) UpdateNews(ctx context.Context, news models.News) (models.News, error) {
	query := `UPDATE news SET author=$1, body=$2  Where id=$3`

	if _, err := s.db.ExecContext(ctx, query, news.Author, news.Body, news.ID); err != nil {
		return models.News{}, fmt.Errorf("Update news: %s", err)
	}

	News, err := s.GetNews(ctx, news.ID)
	if err != nil {
		return models.News{}, fmt.Errorf("get news: %s", err)
	}

	return News, nil
}
