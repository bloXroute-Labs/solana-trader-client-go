package utils

import (
	"sync"
)

type RequestID struct {
	id   uint64
	lock *sync.Mutex
}

func NewRequestID() RequestID {
	return RequestID{
		id:   1,
		lock: &sync.Mutex{},
	}
}

func (r *RequestID) Current() uint64 {
	return r.id
}

func (r *RequestID) Next() uint64 {
	r.lock.Lock()
	defer r.lock.Unlock()

	val := r.id
	r.id++
	return val
}
