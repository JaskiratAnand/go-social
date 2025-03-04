package main

import (
	"net/http"
	"testing"

	"github.com/JaskiratAnand/go-social/internal/store/cache"
	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {
	withRedis := config{
		redisCfg: redisConfig{
			enabled: true,
		},
	}

	app := TestMockApplication(t, withRedis)
	mux := app.mount()

	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})

	t.Run("should hit the cache first and if user not exist then set cache", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserCache)

		mockCacheStore.On("Get", "uuid-1").Return(nil, nil)
		mockCacheStore.On("Get", "uuid-2").Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.AssertNumberOfCalls(t, "Get", 2)
		mockCacheStore.Calls = nil
	})

	t.Run("should NOT hit cache if not enabled", func(t *testing.T) {
		withRedis := config{
			redisCfg: redisConfig{
				enabled: true,
			},
		}

		app := TestMockApplication(t, withRedis)
		mux := app.mount()

		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserCache)

		mockCacheStore.On("Get", "uuid-1").Return(nil, nil)
		mockCacheStore.On("Get", "uuid-2").Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.AssertNotCalled(t, "Get")
		mockCacheStore.Calls = nil
	})

}
