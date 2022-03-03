package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/uptrace/bun"

	"github.com/jrdnull/microservices_technical_test/user_service/internal"
)

// UserRepository provides access to users persisted in Postgres.
type UserRepository struct {
	db *bun.DB
}

// NewUserRepository returns a new UserRepository backed by db.
func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db}
}

// Create user setting ID.
func (r *UserRepository) Create(ctx context.Context, user *internal.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

// SetNewsFeedTags sets the tag ids for the users news feed.
func (r *UserRepository) SetNewsFeedTags(ctx context.Context, id int64, tagIDs []int64) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				log.Printf("tx rollback: %v", txErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.NewDelete().Model((*UserNewsFeed)(nil)).Where("user_id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete existing tags %w", err)
	}

	newsFeed := make([]UserNewsFeed, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		newsFeed = append(newsFeed, UserNewsFeed{UserID: id, TagID: tagID})
	}
	_, err = tx.NewInsert().Model(&newsFeed).Exec(ctx)
	return err
}

// GetNewsFeedTags returns the users tag ids for their news feed.
func (r *UserRepository) GetNewsFeedTags(ctx context.Context, id int64) ([]int64, error) {
	var tagIDs []int64
	err := r.db.NewSelect().Model((*UserNewsFeed)(nil)).
		Column("tag_id").
		Where("user_id = ?", id).
		Scan(ctx, &tagIDs)
	return tagIDs, err
}

type UserNewsFeed struct {
	bun.BaseModel `bun:"user_newsfeed"`
	UserID        int64 `bun:",pk,notnull"`
	TagID         int64 `bun:",pk,notnull"`
}

func (*UserNewsFeed) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey("(user_id) REFERENCES users (id) ON DELETE CASCADE")
	return nil
}
