package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestStroreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	data := bytes.NewReader([]byte("hello deleted f"))
	if err := s.writeStream("deletepic", data); err != nil {
		t.Error(err)
	}

	if err := s.Delete("deletepic"); err != nil {
		t.Error(err)
	}
}
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
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	key := "testpic"
	s := NewStore(opts)
	data := []byte("hello worldasdasdd f")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	f, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}
	b, _ := io.ReadAll(f)
	if string(b) != string(data) {
		t.Error(t, "store read test failed")
	}
	s.Delete(key)
}
