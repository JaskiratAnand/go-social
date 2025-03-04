package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JaskiratAnand/go-social/internal/auth"
	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/JaskiratAnand/go-social/internal/store/cache"
	"go.uber.org/zap"
)

func TestMockApplication(t *testing.T, withRedis config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	testAuth := &auth.TestAuthenticator{}

	var mockCache cache.Storage
	if withRedis.redisCfg.enabled {
		mockCache = cache.NewMockCache()
	}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCache,
		authenticator: testAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
