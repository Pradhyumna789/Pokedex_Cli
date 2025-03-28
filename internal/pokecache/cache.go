package pokecache

import (
  "sync"
  "time"
)

type Cache struct {
  data map[string]cacheValue
  mu sync.Mutex 
  interval time.Duration
}

type cacheValue struct {
  createdAt time.Time
  val []byte
}

func NewCache(interval time.Duration) *Cache {

  cache := &Cache {
    data: make(map[string]cacheValue),
    interval: interval,
  }
  
  go cache.reapLoop()
  return cache

}

func (c *Cache) Add(key string, val []byte) {

  c.mu.Lock()
  defer c.mu.Unlock()

  c.data[key] = cacheValue {
    createdAt: time.Now(),
    val: val,
  }
  
}

func (c *Cache) Get(key string) ([]byte, bool) {

  c.mu.Lock()
  defer c.mu.Unlock()

  cacheVal, found := c.data[key] 
  if !found || time.Since(cacheVal.createdAt) > c.interval {
    delete(c.data, key)
    return nil, false
  }

  return cacheVal.val, true

}

func (c *Cache) reapLoop() {

  ticker := time.NewTicker(c.interval) 
  defer ticker.Stop()
  for {
    <- ticker.C
    c.mu.Lock()
    for key, cacheValue := range c.data {
      if time.Since(cacheValue.createdAt) > c.interval {
        delete(c.data, key)
      } 
    }
    c.mu.Unlock()
  }
}





