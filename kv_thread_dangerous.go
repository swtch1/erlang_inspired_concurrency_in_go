package main

// dangerousStore stores and loads key / value pairs.  It is not thread safe.
type dangerousStore struct {
	data map[string]string
}

func newDangerousStore() *dangerousStore {
	return &dangerousStore{
		data: make(map[string]string),
	}
}

func (s *dangerousStore) Store(k, v string) {
	s.data[k] = v
}

func (s *dangerousStore) Load(k string) (string, bool) {
	v, ok := s.data[k]
	return v, ok
}
