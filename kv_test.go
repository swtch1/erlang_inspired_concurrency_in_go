package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKVStore(t *testing.T) {
	kv := NewSafeStore()

	// create some concurrency that will fail when the race detector runs if not
	// handled properly
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			kv.Store("foo", "bar")
			kv.Store("baz", "bam")
			kv.Load("foo")
		}()
	}

	wg.Wait()

	v, ok := kv.Load("foo")
	require.True(t, ok)
	require.Equal(t, "bar", v)

	_, ok = kv.Load("not-found")
	require.False(t, ok)

	v, ok = kv.Load("baz")
	require.True(t, ok)
	require.Equal(t, "bam", v)
}
