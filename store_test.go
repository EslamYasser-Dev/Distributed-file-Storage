package main

import (
	"bytes"
	"testing"
)

func TestStroreDeleteKey(t *testing.T) {

}
func TestPathTransformFunc(t *testing.T) {
	key := "testpic"
	pathKey := CASPathTransformFunc(key)
	ExpextedOriginalKey := "bb73aaafa1596e5425dc514a361ad4ef658f2758"
	expectedPath := "bb73a/aafa1/596e5/425dc/514a3/61ad4/ef658/f2758"
	if pathKey.PathName != expectedPath {
		t.Error(t, "path transform func test failed")
	}
	if pathKey.FileName != ExpextedOriginalKey {
		t.Error(t, "path transform func test failed")
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	data := bytes.NewReader([]byte("hello worldasdasdd f"))
	if err := s.writeStream("upload", data); err != nil {
		t.Error(err)
	}

}
