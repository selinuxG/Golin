package global

import "sync"

func LoadOrDefault(m *sync.Map, key, def string) string {
	if v, ok := m.Load(key); ok {
		return v.(string)
	}
	return def
}
