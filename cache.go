package memory_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
)

var MemoryCache = InitCache()

func InitCache() *bigcache.BigCache {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
	if err != nil {
		fmt.Println(err)
	}
	return cache
}
