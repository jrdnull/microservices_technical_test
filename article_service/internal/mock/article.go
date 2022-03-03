package mock

import (
	"context"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
)

// ArticleRepository is a test internal.ArticleRepository.
type ArticleRepository struct {
	internal.ArticleRepository

	ListFn func(context.Context, internal.ListOptions) ([]*internal.Article, error)
}

// List calls ListFn.
func (s ArticleRepository) List(ctx context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
	return s.ListFn(ctx, opts)
}
