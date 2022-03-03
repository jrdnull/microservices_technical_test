//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/jrdnull/microservices_technical_test/user_service/internal"
)

func TestCreate(t *testing.T) {
	db := txDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := internal.User{Email: "test@example.org"}
	if err := repo.Create(ctx, &user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID == 0 {
		t.Fatalf("user.ID not set")
	}

	var gotUser internal.User
	err := db.NewSelect().
		Model(&gotUser).
		Where("id = ?", user.ID).
		Scan(ctx)
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if diff := cmp.Diff(gotUser, user); diff != "" {
		t.Errorf("(-got +want):\n%s", diff)
	}
}

func TestSetNewsFeedTags(t *testing.T) {
	db := txDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := internal.User{Email: "test@example.org"}
	if err := repo.Create(ctx, &user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	tags := []int64{1, 2, 3}
	if err := repo.SetNewsFeedTags(ctx, user.ID, tags); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count, err := db.NewSelect().
		Model((*UserNewsFeed)(nil)).
		Where("user_id = ?", user.ID).
		Count(ctx)
	if err != nil {
		t.Fatalf("count: %v", err)
	}

	if count != len(tags) {
		t.Errorf("count tags, got: %v; want %v", count, len(tags))
	}

	t.Run("resetting tags", func(t *testing.T) {
		tags := []int64{4, 5}
		if err := repo.SetNewsFeedTags(ctx, user.ID, tags); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		count, err := db.NewSelect().
			Model((*UserNewsFeed)(nil)).
			Where("user_id = ?", user.ID).
			Count(ctx)
		if err != nil {
			t.Fatalf("count: %v", err)
		}

		if count != len(tags) {
			t.Errorf("count tags, got: %v; want %v", count, len(tags))
		}
	})
}

func TestGetNewsFeedTags(t *testing.T) {
	db := txDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := internal.User{Email: "test@example.org"}
	if err := repo.Create(ctx, &user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	tags := []int64{1, 2, 3}
	if err := repo.SetNewsFeedTags(ctx, user.ID, tags); err != nil {
		t.Fatalf("set tags: %v", err)
	}

	gotTags, err := repo.GetNewsFeedTags(ctx, user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff(gotTags, tags); diff != "" {
		t.Errorf("(-got +want):\n%s", diff)
	}
}
