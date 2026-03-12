package common

import (
	"sync"
	"time"

	"github.com/dmwork-org/dmwork-lib/pkg/redis"
)

// RedisCache redis缓存
type RedisCache struct {
	conn *redis.Conn
}

// NewRedisCache 创建
func NewRedisCache(addr string, password string) *RedisCache {
	r := &RedisCache{}
	r.conn = redis.New(addr, password)
	return r
}

// Set Set
func (r *RedisCache) Set(key string, value string) error {
	return r.conn.Set(key, value)
}

// Delete 删除key
func (r *RedisCache) Delete(key string) error {
	return r.conn.Del(key)
}

// SetAndExpire 包含过期时间
func (r *RedisCache) SetAndExpire(key string, value string, expire time.Duration) error {
	return r.conn.SetAndExpire(key, value, expire)
}

// Get Get
func (r *RedisCache) Get(key string) (string, error) {
	return r.conn.GetString(key)
}

// GetRedisConn 获取redis连接
func (r *RedisCache) GetRedisConn() *redis.Conn {
	return r.conn
}

// MemoryCache 内存缓存
type MemoryCache struct {
	cacheMap map[string]string
	timers   map[string]*time.Timer
	sync.RWMutex
}

// NewMemoryCache 创建
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cacheMap: map[string]string{},
		timers:   map[string]*time.Timer{},
	}
}

// Set Set
func (m *MemoryCache) Set(key string, value string) error {
	m.Lock()
	defer m.Unlock()
	m.cacheMap[key] = value
	// Cancel any existing expiration timer for this key
	if t, ok := m.timers[key]; ok {
		t.Stop()
		delete(m.timers, key)
	}
	return nil
}

// SetAndExpire SetAndExpire
func (m *MemoryCache) SetAndExpire(key string, value string, expire time.Duration) error {
	m.Lock()
	defer m.Unlock()
	m.cacheMap[key] = value
	// Cancel any existing timer for this key to prevent stale deletion
	if t, ok := m.timers[key]; ok {
		t.Stop()
		delete(m.timers, key)
	}
	if expire > 0 {
		t := time.AfterFunc(expire, func() {
			m.Lock()
			defer m.Unlock()
			// Only delete if this timer is still the current one for the key.
			// A newer SetAndExpire may have replaced it while we waited for the lock.
			if cur, ok := m.timers[key]; ok && cur == t {
				delete(m.cacheMap, key)
				delete(m.timers, key)
			}
		})
		m.timers[key] = t
	}
	return nil
}

// Get Get
func (m *MemoryCache) Get(key string) (string, error) {
	m.RLock()
	defer m.RUnlock()
	return m.cacheMap[key], nil
}

// Delete Delete
func (m *MemoryCache) Delete(key string) error {
	m.Lock()
	defer m.Unlock()
	delete(m.cacheMap, key)
	if t, ok := m.timers[key]; ok {
		t.Stop()
		delete(m.timers, key)
	}
	return nil
}
