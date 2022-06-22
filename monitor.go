package gomsg

import (
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
	ants "github.com/panjf2000/ants/v2"
)

type MsgFn = func(key string, val interface{})

type monitors struct {
	size int
	msg  cmap.ConcurrentMap[[]MsgFn]
	pool *ants.Pool
}

type Monitor interface {
	WaitFor(key string, fn MsgFn)
	Send(key string, data interface{})
}

func New(size int) *monitors {
	pool, _ := ants.NewPool(size)
	cmap.SHARD_COUNT = size
	m := monitors{
		size: size,
		msg:  cmap.New[[]MsgFn](),
		pool: pool,
	}
	return &m
}

func (m *monitors) WaitFor(key string, fn MsgFn) {
	s := []MsgFn{fn}
	v, exist := m.msg.Get(key)
	if exist {
		v = append(s, v...)
	}
	m.msg.Set(key, s)
}

func (m *monitors) Send(key string, data interface{}) {
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

func (m *monitors) AsyncSend(key string, data interface{}) {
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
