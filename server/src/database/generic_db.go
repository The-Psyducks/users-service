package database

import (
	"sort"
	"time"
)

type Timestamped interface {
	GetCreatedAt() time.Time
}

type genericDB[T Timestamped] struct {
	data map[string]T
}

func (db genericDB[T]) Create(data T, key string) (T, error) {
	db.data[key] = data
	return data, nil
}

func (db genericDB[T]) Get(key string) (T, error) {
	item, found := db.data[key]

	if !found {
		return item, ErrKeyNotFound
	}

	return item, nil
}

func (db genericDB[T]) Delete(key string) error {
	_, found := db.data[key]

	if !found {
		return ErrKeyNotFound
	}

	delete(db.data, key)
	return nil
}

func (db genericDB[T]) GetAll() ([]T, error) {
	allItems := make([]T, 0, len(db.data))

	for _, value := range db.data {
		allItems = append(allItems, value)
	}

	return allItems, nil
}

func (db genericDB[T]) GetAllInOrder() ([]T, error) {
	allItems := make([]T, 0, len(db.data))

	for _, value := range db.data {
		allItems = append(allItems, value)
	}

	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].GetCreatedAt().After(allItems[j].GetCreatedAt())
	})

	return allItems, nil
}

func CreateDB[T Timestamped]() genericDB[T] {
	data := make(map[string]T, 1)
	return genericDB[T]{data}
}