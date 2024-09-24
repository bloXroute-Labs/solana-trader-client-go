package utils

import (
	cmap "github.com/mhmtszr/concurrent-swiss-map"
)

type LockedMap[T comparable, V any] struct {
	data *cmap.CsMap[T, V]
}

func NewLockedMap[T comparable, V any]() *LockedMap[T, V] {
	return &LockedMap[T, V]{
		data: cmap.Create[T, V](),
	}
}

func (fs *LockedMap[T, V]) Get(key T) (V, bool) {
	return fs.data.Load(key)
}

// GetOrInsert inserts the entry if the key doesn't exist
func (fs *LockedMap[T, V]) GetOrInsert(key T, f func() V) (out V) {

	fs.data.SetIf(key, func(value V, found bool) (val V, set bool) {
		if found {
			out = value
			set = false
			val = value
			return
		}

		val = f()
		out = val
		set = true
		return
	})
	return
}

func (fs *LockedMap[T, V]) Set(key T, v V) {
	fs.data.Store(key, v)
}

func (fs *LockedMap[T, V]) Delete(key T) {
	fs.data.Delete(key)
}

func (fs *LockedMap[T, V]) DeleteAll() {
	fs.data.Clear()
}

func (fs *LockedMap[T, V]) DeleteWithCondition(cond func(V) bool) {
	fs.data.Range(func(key T, value V) (stop bool) {
		if cond(value) {
			fs.data.Delete(key)
		}
		return false
	})
}

func (fs *LockedMap[T, V]) ExistOrAdd(key T, value V) (exists bool) {
	fs.data.SetIf(key, func(_ V, found bool) (val V, set bool) {
		exists = found
		set = !found
		val = value
		return
	})
	return
}

func (fs *LockedMap[T, V]) Update(key T, update func(value V, exists bool) V) (val V) {
	fs.data.SetIf(key, func(valueInMap V, found bool) (valueOut V, set bool) {
		val = update(valueInMap, found)
		return val, true
	})
	return
}

func (fs *LockedMap[T, V]) Len() int {
	return fs.data.Count()
}

func (fs *LockedMap[T, V]) Keys() []T {
	var keys []T
	fs.data.Range(func(key T, _ V) (stop bool) {
		keys = append(keys, key)
		return false
	})

	return keys
}

func (fs *LockedMap[T, V]) Values() []V {
	var values []V
	fs.data.Range(func(_ T, value V) (stop bool) {
		values = append(values, value)
		return false
	})

	return values
}

func (fs *LockedMap[T, V]) Copy() map[T]V {
	m := make(map[T]V)
	fs.data.Range(func(key T, value V) (stop bool) {
		m[key] = value
		return false
	})

	return m
}
