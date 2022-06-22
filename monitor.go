package gomsg

import (
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
	ants "github.com/panjf2000/ants/v2"
)

type T any

type MsgFn[T any] func(key string, val T)

type monitors[T any] struct {
	size int
	msg  cmap.ConcurrentMap[[]MsgFn[T]]
	pool *ants.Pool
}

type Monitor[T any] interface {
	WaitFor(key string, fn MsgFn[T])
	Send(key string, data T)
}

func New[T any](size int) Monitor[T] {
	pool, _ := ants.NewPool(size)
	cmap.SHARD_COUNT = size
	m := monitors[T]{
		size: size,
		msg:  cmap.New[[]MsgFn[T]](),
		pool: pool,
	}
	return &m
}

func (m *monitors[T]) WaitFor(key string, fn MsgFn[T]) {
	m.msg.Upsert(key, []MsgFn[T]{fn}, func(exist bool, v, nv []MsgFn[T]) []MsgFn[T] {
		if exist {
			v = append(v, nv...)
		}
		return nv
	})
}

func (m *monitors[T]) Send(key string, data T) {
	v, exist := m.msg.Get(key)
	if !exist {
		return
	}
	wg := sync.WaitGroup{}
	for i := range v {
		wg.Add(1)
		_ = m.pool.Submit(func() {
			defer wg.Done()
			v[i](key, data)
		})
	}
	wg.Wait()
}

func (m *monitors[T]) AsyncSend(key string, data T) {
	v, exist := m.msg.Get(key)
	if !exist {
		return
	}
	for i := range v {
		_ = m.pool.Submit(func() {
			v[i](key, data)
		})
	}
}
