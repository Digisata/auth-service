package repository

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/memcached"
)

type CacheRepository struct {
	memcachedDB *memcached.Database
}

func NewCacheRepository(memcachedDB *memcached.Database) *CacheRepository {
	return &CacheRepository{
		memcachedDB: memcachedDB,
	}
}

func (cr CacheRepository) Set(req domain.CacheItem) error {
	item := &memcache.Item{
		Key:        req.Key,
		Value:      []byte(req.Value),
		Expiration: int32(req.Exp),
	}

	if err := cr.memcachedDB.Set(item); err != nil {
		return err
	}

	return nil
}

func (cr CacheRepository) Get(key string) (domain.CacheItem, error) {
	var item domain.CacheItem

	it, err := cr.memcachedDB.Get(key)
	if err != nil {
		return item, err
	}

	item = domain.CacheItem{
		Key:   it.Key,
		Value: string(it.Value),
		Exp:   int(it.Expiration),
	}

	return item, nil
}

func (cr CacheRepository) Delete(key string) error {
	err := cr.memcachedDB.Delete(key)
	if err != nil {
		return err
	}

	return nil
}
