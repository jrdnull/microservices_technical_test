package internal

import (
	"context"
	"time"
)

// ArticleRepository provides access to articles.
type ArticleRepository interface {
	List(ctx context.Context, opts ListOptions) ([]*Article, error)
}

// ArticleService provides methods for working with articles.
type ArticleService struct {
	r ArticleRepository
}

// NewArticleService returns a new ArticleService backed by r.
func NewArticleService(r ArticleRepository) *ArticleService {
	return &ArticleService{r}
}

// ListOptions provides optional parameters to the List method.
type ListOptions struct {
	TagIDs       []int64
	MatchAllTags bool
	Page         int
}

// List articles paginated matching opts.
func (s *ArticleService) List(
	ctx context.Context, opts ListOptions,
) ([]*Article, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}

	return s.r.List(ctx, opts)
}

// Article represents a news article.
type Article struct {
	ID        int64     `bun:",pk,autoincrement" json:"id"`
	Title     string    `json:"title"`
	Timestamp time.Time `json:"timestamp"`
	Tags      []Tag     `bun:"m2m:article_tags" json:"tags"`
}
