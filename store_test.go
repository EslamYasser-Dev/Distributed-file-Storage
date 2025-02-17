package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "testpic"
	pathKey := CASPathTransformFunc(key)
	fmt.Println(pathKey)
	ExpextedOriginalKey := "ecd4613d7061513dbcc8f28e713d862ae8ea494a"
	expectedPath := "ecd46/13d70/61513/dbcc8/f28e7/13d86/2ae8e/a494a"
	if pathKey.PathName != expectedPath {
		t.Error(t, "path transform func test failed")
	}
	if pathKey.FileName != ExpextedOriginalKey {
		t.Error(t, "path transform func test failed")
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("fooBar_%d", i)
		data := []byte("some bytes of a file is here for test")
		if err := s.Write(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("expected to NOT have key: %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}
		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Error(t, "store read test failed")
		}
		if err := s.Delete(key); err != nil {
			t.Error(err)
		}
	}
}

// helper functions
func newStore() *Store {
	return NewStore(StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	})
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
