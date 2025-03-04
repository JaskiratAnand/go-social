package cache

import (
	"context"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(context.Context, uuid.UUID) (*store.Users, error)
		Set(context.Context, *store.Users) error
		Delete(context.Context, uuid.UUID)
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}
