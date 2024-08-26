package storage

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const duration time.Duration = 120

type Telnum struct {
	Msisdn     string `json:"MSISDN"`
	Region     string `json:"region"`
	Abc        string `json:"abc"`
	Enabled    bool   `json:"enabled"`
	ServiceKey int    `json:"serviceKey"`
}

type Storage map[string]Telnum

type Cache struct {
	data Storage
	*sync.RWMutex
}

func MakeHash(item Telnum) string {
	//Тестовый комментарий
	str := fmt.Sprintf("%+v", item)
	h := sha256.New()
	h.Write([]byte(str))
	value := binary.LittleEndian.Uint64(h.Sum(nil))
	return strconv.FormatUint(value, 16)
}

func (c *Cache) hash(item Telnum) string {
	return MakeHash(item)
}

func (c *Cache) Create(item Telnum) string {
	c.Lock()
	defer c.Unlock()

	key := c.hash(item)
	c.data[key] = item
	return key
}

func (c *Cache) Show() Storage {
	c.Lock()
	defer c.Unlock()
	snapshot := c.data

	return snapshot
}

func (c *Cache) Get(key string) (Telnum, bool) {
	c.RLock()
	defer c.RUnlock()
	row, ok := c.data[key]

	return row, ok
}

func (c *Cache) Update(key string, item Telnum) (string, bool) {
	ok := c.Delete(key)
	if !ok {
		return "", false
	}
	return c.Create(item), true
}

func (c *Cache) Delete(key string) bool {
	c.Lock()
	defer c.Unlock()

	_, ok := c.data[key]
	if ok {
		delete(c.data, key)
	}

	return ok
}

func (c *Cache) Clear() {
	c.Lock()
	defer c.Unlock()

	for key, value := range c.data {
		if !value.Enabled {
			delete(c.data, key)
		}
	}
}

func (c *Cache) Cleanup(d time.Duration) *time.Ticker {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			c.Clear()
		}
	}()

	return ticker
}

func CreateCacheObject() Cache {
	cache := Cache{
		data:    make(Storage),
		RWMutex: &sync.RWMutex{},
	}
	cache.Cleanup(duration * time.Second)
	return cache
}
