package service

import (
	"log"
	"sync"
	"time"
)

type LocalCacheService interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	GetExpireTime(key string) (*time.Time, bool)
	Del(key string)
	Has(key string) bool
}

type cacheItem struct {
	value      interface{}
	expireTime time.Time
}

type localCacheService struct {
	store sync.Map
}

func NewLocalCacheService() LocalCacheService {
	c := &localCacheService{
		store: sync.Map{},
	}
	// Start cleanup ticker
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			c.cleanUp()
		}
	}()
	return c
}

func (c *localCacheService) Set(key string, value interface{}, ttl time.Duration) {
	expire := time.Now().Add(ttl)
	c.store.Store(key, cacheItem{
		value:      value,
		expireTime: expire,
	})
}

func (c *localCacheService) GetExpireTime(key string) (*time.Time, bool) {
	val, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}

	item := val.(cacheItem)
	return &item.expireTime, true
}

func (c *localCacheService) Get(key string) (interface{}, bool) {
	val, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}

	item := val.(cacheItem)
	if time.Now().After(item.expireTime) {
		c.store.Delete(key)
		return nil, false
	}
	return item.value, true
}

func (c *localCacheService) Del(key string) {
	c.store.Delete(key)
}

func (c *localCacheService) Has(key string) bool {
	_, exists := c.Get(key)
	return exists
}

func (c *localCacheService) cleanUp() {
	log.Println("========== clean up cache ==========")
	now := time.Now()
	c.store.Range(func(key, val interface{}) bool {
		item := val.(cacheItem)
		if now.After(item.expireTime) {
			c.store.Delete(key)
		}
		return true
	})
}
