package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/beme/abide"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
	"github.com/jrdnull/microservices_technical_test/article_service/internal/mock"
)

func TestListArticles(t *testing.T) {
	testUser := &internal.User{NewsFeedTags: []int64{1, 2, 3}}
	userSrv := internal.NewUserService(mock.UserRepository{
		GetFn: func(_ context.Context, id int64) (*internal.User, error) {
			return testUser, nil
		},
	})

	article1 := &internal.Article{
		ID:        1,
		Title:     "1",
		Timestamp: time.Time{},
		Tags: []internal.Tag{
			{ID: 1, Name: "1"},
			{ID: 1, Name: "2"},
		},
	}
	article2 := &internal.Article{ID: 1,
		Title:     "2",
		Timestamp: time.Time{},
		Tags: []internal.Tag{
			{ID: 1, Name: "1"},
		},
	}

	doRequest := func(values url.Values, listFn func(context.Context, internal.ListOptions) ([]*internal.Article, error)) *http.Response {
		articleSrv := internal.NewArticleService(mock.ArticleRepository{
			ListFn: listFn,
		})

		e := NewServer(articleSrv, userSrv, Config{ArticleTagFilterInputs: 2})
		req := httptest.NewRequest(http.MethodGet, "/v1/articles?"+values.Encode(), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		return rec.Result()
	}
	t.Run("unauthenticated", func(t *testing.T) {
		resp := doRequest(url.Values{}, nil)
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
	t.Run("repo error", func(t *testing.T) {
		resp := doRequest(url.Values{
			"auth_key": {"user1"},
		}, func(_ context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
			return nil, errors.New("don't expose internal errors")
		})
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
	t.Run("panic", func(t *testing.T) {
		resp := doRequest(url.Values{
			"auth_key": {"user1"},
		}, func(_ context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
			panic("oops")
		})
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
	t.Run("unfiltered", func(t *testing.T) {
		resp := doRequest(url.Values{
			"auth_key": {"user1"},
		}, func(_ context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
			return []*internal.Article{article1, article2}, nil
		})
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
	t.Run("all_tags, incorrect amount", func(t *testing.T) {
		resp := doRequest(url.Values{
			"auth_key": {"user1"},
			"all_tags": {"1,2,3"},
		}, func(_ context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
			return nil, nil
		})
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
	t.Run("all_tags, invaid", func(t *testing.T) {
		resp := doRequest(url.Values{
			"auth_key": {"user1"},
			"all_tags": {"foo,bar,baz"},
		}, func(_ context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
			return nil, nil
		})
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
	t.Run("all_tags", func(t *testing.T) {
		resp := doRequest(url.Values{
			"auth_key": {"user1"},
			"all_tags": {"1, 2"},
		}, func(_ context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
			if !opts.MatchAllTags {
				t.Errorf("opts.MatchAllTags => false; want true")
			}
			if want := []int64{1, 2}; !reflect.DeepEqual(opts.TagIDs, want) {
				t.Errorf("opts.TagIDs => %v; want %v", opts.TagIDs, want)
			}
			return []*internal.Article{article1}, nil
		})
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
	t.Run("all_tags", func(t *testing.T) {
		resp := doRequest(url.Values{
			"auth_key": {"user1"},
			"any_tags": {"1, 2, 3"},
		}, func(_ context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
			if opts.MatchAllTags {
				t.Errorf("opts.MatchAllTags => true; want false")
			}
			if want := []int64{1, 2, 3}; !reflect.DeepEqual(opts.TagIDs, want) {
				t.Errorf("opts.TagIDs => %v; want %v", opts.TagIDs, want)
			}
			return []*internal.Article{article1, article2}, nil
		})
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
}

func TestGetUserFeed(t *testing.T) {
	testUser := &internal.User{NewsFeedTags: []int64{1, 2, 3}}
	userSrv := internal.NewUserService(mock.UserRepository{
		GetFn: func(_ context.Context, id int64) (*internal.User, error) {
			return testUser, nil
		},
	})

	article1 := &internal.Article{
		ID:        1,
		Title:     "1",
		Timestamp: time.Time{},
		Tags: []internal.Tag{
			{ID: 1, Name: "1"},
			{ID: 1, Name: "2"},
		},
	}
	article2 := &internal.Article{ID: 1,
		Title:     "2",
		Timestamp: time.Time{},
		Tags: []internal.Tag{
			{ID: 1, Name: "1"},
		},
	}

	doRequest := func(values url.Values, listFn func(context.Context, internal.ListOptions) ([]*internal.Article, error)) *http.Response {
		articleSrv := internal.NewArticleService(mock.ArticleRepository{
			ListFn: listFn,
		})

		e := NewServer(articleSrv, userSrv, Config{ArticleTagFilterInputs: 2})
		req := httptest.NewRequest(http.MethodGet, "/v1/articles/feed?"+values.Encode(), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		return rec.Result()
	}
	t.Run("unauthenticated", func(t *testing.T) {
		resp := doRequest(url.Values{}, nil)
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
	t.Run("repo error", func(t *testing.T) {
		resp := doRequest(url.Values{
			"auth_key": {"user1"},
		}, func(_ context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
			return nil, errors.New("don't expose internal errors")
		})
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
	t.Run("panic", func(t *testing.T) {
		resp := doRequest(url.Values{
			"auth_key": {"user1"},
		}, func(_ context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
			panic("oops")
		})
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
	t.Run("filtered by user tags", func(t *testing.T) {
		resp := doRequest(url.Values{
			"auth_key": {"user1"},
		}, func(_ context.Context, opts internal.ListOptions) ([]*internal.Article, error) {
			if opts.MatchAllTags {
				t.Errorf("opts.MatchAllTags => true; want false")
			}
			if want := testUser.NewsFeedTags; !reflect.DeepEqual(opts.TagIDs, want) {
				t.Errorf("opts.TagIDs => %v; want %v", opts.TagIDs, want)
			}
			return []*internal.Article{article1, article2}, nil
		})
		abide.AssertHTTPResponse(t, t.Name(), resp)
	})
}
