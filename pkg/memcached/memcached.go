package memcached

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type Config struct {
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
}

type Database struct {
	Mc *memcache.Client
}

func NewDatabase(cfg Config) *Database {
	return &Database{
		Mc: memcache.New(fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort)),
	}
}

func (db Database) Set(req *memcache.Item) error {
	if err := db.Mc.Set(req); err != nil {
		return err
	}

	return nil
}

func (db Database) Get(key string) (*memcache.Item, error) {
	it, err := db.Mc.Get(key)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func (db Database) Delete(key string) error {
	err := db.Mc.Delete(key)
	if err != nil {
		return err
	}

	return nil
}
