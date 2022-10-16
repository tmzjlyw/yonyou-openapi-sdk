package cache

import (
	"fmt"
	"github.com/allegro/bigcache"
	"time"
)




var cache, _  = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))



func main() {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	cache.Set("my-unique-key", []byte("value"))
	entry, _ := cache.Get("my-unique-key")
	fmt.Println(string(entry))
}
