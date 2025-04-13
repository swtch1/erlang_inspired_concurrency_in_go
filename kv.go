package main

type operation int

const (
	opUnknown operation = iota
	opStore
	opLoad
)

type kvEvent struct {
	op operation
	k  string
	v  string
	ok bool
}

// KVStore is a thread safe wrapper of kvStore that uses an Erlang-inspired
// model for allowing concurrent operations by passing messages to a synchronous
// process.  This is much like the actor model, but a contrived implementation
// which is intended to more closely model (my limited understanding of)
// Erlang's behind-the-scenes handling of behaviors.
//
// Ideas inspired from
// https://stevana.github.io/erlangs_not_about_lightweight_processes_and_message_passing.html
type KVStore struct {
	inner *kvStore
	in    chan kvEvent
	out   chan kvEvent
}

func NewKVStore() *KVStore {
	kvstore := &KVStore{
		inner: newkvStore(),
		in:    make(chan kvEvent),
		out:   make(chan kvEvent),
	}

	// The goroutine handles all operations synchronously by passing messages to
	// the inner kvStore.
	go func() {
		for event := range kvstore.in {
			switch event.op {
			case opUnknown:
				panic("we don't know")
			case opStore:
				kvstore.inner.Store(event.k, event.v)
				kvstore.out <- kvEvent{}
			case opLoad:
				v, ok := kvstore.inner.Load(event.k)
				kvstore.out <- kvEvent{
					v:  v,
					ok: ok,
				}
			}
		}
	}()

	return kvstore
}

func (s *KVStore) Store(k, v string) {
	s.in <- kvEvent{
		op: opStore,
		k:  k,
		v:  v,
	}
	<-s.out
}

func (s *KVStore) Load(k string) (string, bool) {
	s.in <- kvEvent{
		op: opLoad,
		k:  k,
	}
	res := <-s.out
	return res.v, res.ok
}

// kvStore stores and loads key / value pairs.  It is not thread safe.
type kvStore struct {
	data map[string]string
}

func newkvStore() *kvStore {
	return &kvStore{
		data: make(map[string]string),
	}
}

func (s *kvStore) Store(k, v string) {
	s.data[k] = v
}

func (s *kvStore) Load(k string) (string, bool) {
	v, ok := s.data[k]
	return v, ok
}
