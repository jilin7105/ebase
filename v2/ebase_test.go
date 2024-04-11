package ebase

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jilin7105/ebase/config"
	"gorm.io/gorm"
	"testing"
)

func TestInitEbase(t *testing.T) {
	cfg := config.Config{
		Databases: []config.DbConfig{
			{
				Name:     "testDB",
				Host:     "127.0.0.1",
				Port:     3306,
				Username: "root",
				Password: "root",
				Dbname:   "test",
			},
		},
		Redis: []config.RedisConfig{
			{
				Name:     "testRedis",
				Host:     "127.0.0.1",
				Port:     6379,
				Password: "",
				DB:       0,
			},
		},
	}

	var eb = Eb{
		cxt:    context.Background(),
		Config: cfg,
		Redis:  map[string]*redis.Client{},
		DBs:    map[string]*gorm.DB{},
	}
	eb.initMysql()
	eb.initRedis()
	if len(eb.DBs) != 1 {
		t.Errorf("Expected 1 DB, got %v", len(eb.DBs))
	}

	if len(eb.Redis) != 1 {
		t.Errorf("Expected 1 Redis client, got %v", len(eb.Redis))
	}
}
