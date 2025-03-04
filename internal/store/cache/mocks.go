package cache

import (
	"context"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/google/uuid"
)

func NewMockCache() Storage {
	return Storage{
		Users: &MockUserCache{},
	}
}

type MockUserCache struct{}

func (m MockUserCache) Get(ctx context.Context, id uuid.UUID) (*store.Users, error) {
	return nil, nil
}

func (m MockUserCache) Set(ctx context.Context, user *store.Users) error {
	return nil
}

func (m MockUserCache) Delete(ctx context.Context, id uuid.UUID) {}
