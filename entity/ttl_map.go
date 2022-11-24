package entity

import (
	"sync"
	"time"
)

type TTLMap struct {
	TTL time.Duration

	data sync.Map
}

type expireEntry struct {
	ExpiresAt time.Time
	Value     interface{}
}

func (t *TTLMap) Store(key string, val interface{}) {
	t.data.Store(key, expireEntry{
		ExpiresAt: time.Now().Add(t.TTL),
		Value:     val,
	})
}

func (t *TTLMap) Load(key string) (val interface{}) {
	entry, ok := t.data.Load(key)
	if !ok {
		return nil
	}

	expireEntry := entry.(expireEntry)
	if expireEntry.ExpiresAt.After(time.Now()) {
		return nil
	}

	return expireEntry.Value
}

func NewTTLMap(ttl time.Duration) (m TTLMap) {
	m.TTL = ttl

	go func() {
		for now := range time.Tick(time.Second) {
			m.data.Range(func(k, v interface{}) bool {
				if v.(expireEntry).ExpiresAt.After(now) {
					m.data.Delete(k)
				}
				return true
			})
		}
	}()

	return
}
