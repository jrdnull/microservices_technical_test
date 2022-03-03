package internal_test

import (
	"context"
	"testing"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
	"github.com/jrdnull/microservices_technical_test/article_service/internal/mock"
)

func TestArticleServiceList(t *testing.T) {
	t.Run("defaults page size", func(t *testing.T) {
		srv := internal.NewArticleService(mock.ArticleRepository{
			ListFn: func(ctx context.Context, options internal.ListOptions) ([]*internal.Article, error) {
				if options.Page != 1 {
					t.Errorf("called with options.Page: %d; want: 1", options.Page)
				}
				return nil, nil
			},
		})

		_, err := srv.List(context.Background(), internal.ListOptions{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
