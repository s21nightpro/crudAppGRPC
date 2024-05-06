package cache

import (
	"database/sql"
	_go "github.com/s21nightpro/crudAppGRPC/internal/grpc/user"
	"sync"
	"time"
)

type server struct {
	_go.UnimplementedUserServiceServer
	users map[string]*_go.User
	cache *Cache
	db    *sql.DB
	mu    sync.Mutex
}
type Cache struct {
	items map[string]Item
	mu    sync.Mutex
}

type Item struct {
	Value      interface{}
	Expiration int64
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]Item),
	}
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiration := time.Now().Add(duration).UnixNano()
	c.items[key] = Item{Value: value, Expiration: expiration}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	if time.Now().UnixNano() > item.Expiration {
		delete(c.items, key)
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}
