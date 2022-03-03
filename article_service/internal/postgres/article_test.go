//go:build integration
// +build integration

package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/uptrace/bun/dbfixture"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
)

func TestList(t *testing.T) {
	db := txDB(t)
	defer db.Close()

	fixture := dbfixture.New(db, dbfixture.WithTruncateTables())
	err := fixture.Load(context.Background(), os.DirFS("./testdata"), "fixture.yaml")
	if err != nil {
		t.Fatalf("fixture.Load: %v", err)
	}

	repo := NewArticleRepository(db)
	ctx := context.Background()

	parseTime := func(s string) time.Time {
		tm, err := time.Parse(time.RFC3339, s)
		if err != nil {
			t.Fatalf("parseTime(%q): %v", s, err)
		}
		return tm
	}
	article1 := &internal.Article{
		ID:        1,
		Title:     "Redbull Hardline Results",
		Timestamp: parseTime("2006-01-03T00:00:00Z"),
		Tags: []internal.Tag{
			{ID: 2, Name: "sports"},
		},
	}
	article2 := &internal.Article{
		ID:        2,
		Title:     "Redbull Hardline - Bike Thieves!",
		Timestamp: parseTime("2006-01-02T00:00:00Z"),
		Tags: []internal.Tag{
			{ID: 1, Name: "crime"},
			{ID: 2, Name: "sports"},
			{ID: 3, Name: "business"},
		},
	}
	article3 := &internal.Article{
		ID:        3,
		Title:     "Mystery",
		Timestamp: parseTime("2006-01-01T00:00:00Z"),
	}

	cases := map[string]struct {
		want       []*internal.Article
		opts       internal.ListOptions
		beforeTest func() (afterTest func())
	}{
		"pagination 1": {
			want: []*internal.Article{article1, article2},
			opts: internal.ListOptions{Page: 1},
			beforeTest: func() (afterTest func()) {
				undo := setArticlePageSize(2)
				return undo
			},
		},
		"pagination 2": {
			want: []*internal.Article{article3},
			opts: internal.ListOptions{Page: 2},
			beforeTest: func() (afterTest func()) {
				undo := setArticlePageSize(2)
				return undo
			},
		},
		"pagination 3": {
			want: []*internal.Article{},
			opts: internal.ListOptions{Page: 3},
			beforeTest: func() (afterTest func()) {
				undo := setArticlePageSize(2)
				return undo
			},
		},
		"has any tags": {
			want: []*internal.Article{article1, article2},
			opts: internal.ListOptions{Page: 1, TagIDs: []int64{2, 3}},
		},
		"has all tags": {
			want: []*internal.Article{article2},
			opts: internal.ListOptions{Page: 1, TagIDs: []int64{1, 2, 3}, MatchAllTags: true},
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			if tc.beforeTest != nil {
				afterTest := tc.beforeTest()
				defer afterTest()
			}

			got, err := repo.List(ctx, tc.opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("(-got +want):\n%s", diff)
			}
		})
	}
}

func setArticlePageSize(n int) (undo func()) {
	old := articlePageSize
	articlePageSize = n
	return func() {
		articlePageSize = old
	}
}
