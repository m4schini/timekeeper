package cache

import "time"

type Cache interface {
	Set(key string, value string, ttl time.Duration)
	Get(key string) (value string, expiresAt time.Time, exists bool)
}

type InMemory struct {
	data map[string]inMemoryItem
}

type inMemoryItem struct {
	value     string
	expiresAt time.Time
}

func NewInMemory() *InMemory {
	return &InMemory{data: make(map[string]inMemoryItem)}
}

func (i *InMemory) Set(key string, value string, ttl time.Duration) {
	i.data[key] = inMemoryItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (i *InMemory) Get(key string) (value string, expiresAt time.Time, exists bool) {
	item, exists := i.data[key]
	if !exists {
		return "", expiresAt, false
	}
	if item.expiresAt.Before(time.Now()) {
		delete(i.data, key)
		return "", item.expiresAt, false
	}
	return item.value, item.expiresAt, true
}
