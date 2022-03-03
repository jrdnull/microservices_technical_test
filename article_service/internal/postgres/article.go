package postgres

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
)

var articlePageSize = 20

// ArticleRepository provides access to users persisted in Postgres.
type ArticleRepository struct {
	db *bun.DB
}

// NewArticleRepository returns a new ArticleRepository backed by db.
func NewArticleRepository(db *bun.DB) *ArticleRepository {
	return &ArticleRepository{db}
}

// List articles paginated matching opts.
func (r *ArticleRepository) List(
	ctx context.Context, opts internal.ListOptions,
) ([]*internal.Article, error) {
	articles := []*internal.Article{}
	q := r.db.NewSelect().Model(&articles).Relation("Tags")
	if len(opts.TagIDs) > 0 {
		// we have to join again to exclude so we don't lose other tags
		q.Join("JOIN article_tags AS f ON f.article_id = article.id AND f.tag_id IN (?)", bun.In(opts.TagIDs))
		q.GroupExpr("article.id")
		if opts.MatchAllTags {
			q.Having("COUNT(f.tag_id) = ?", len(opts.TagIDs))
		}
	}
	q.OrderExpr("article.timestamp DESC")
	q.Offset((opts.Page - 1) * articlePageSize).Limit(articlePageSize)

	return articles, q.Scan(ctx)
}

// ArticleTag represents the m2m relationship between articles and tags.
type ArticleTag struct {
	ArticleID int64             `bun:",pk,notnull"`
	Article   *internal.Article `bun:"rel:belongs-to"`
	TagID     int64             `bun:",pk,notnull"`
	Tag       *internal.Tag     `bun:"rel:belongs-to"`
}

func (*ArticleTag) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(article_id) REFERENCES articles (id) ON DELETE CASCADE")
	query.ForeignKey("(tag_id) REFERENCES tags (id) ON DELETE CASCADE")
	return nil
}
