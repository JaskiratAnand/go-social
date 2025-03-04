package cache

import (
	"context"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func NewMockCache() Storage {
	return Storage{
		Users: &MockUserCache{},
	}
}

type MockUserCache struct {
	mock.Mock
}

func (m *MockUserCache) Get(ctx context.Context, userID uuid.UUID) (*store.Users, error) {
	args := m.Called(userID)
	return nil, args.Error(1)
}

func (m *MockUserCache) Set(ctx context.Context, user *store.Users) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserCache) Delete(ctx context.Context, userID uuid.UUID) {}
