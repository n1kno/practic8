package store

import "sync"

// Record – одна запись о пользователе
type Record struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Profession string `json:"profession"`
}

var (
	mu      sync.RWMutex
	records []Record
)

// Add добавляет запись в хранилище
func Add(r Record) {
	mu.Lock()
	defer mu.Unlock()
	records = append(records, r)
}

// GetAll возвращает копию всех записей
func GetAll() []Record {
	mu.RLock()
	defer mu.RUnlock()
	// Возвращаем копию, чтобы избежать гонок при изменении среза извне
	cp := make([]Record, len(records))
	copy(cp, records)
	return cp
}
