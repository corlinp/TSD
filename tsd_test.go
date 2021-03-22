package tsd

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestTSD(t *testing.T) {
	buf := &bytes.Buffer{}

	data1 := []byte("Hello")
	data2 := []byte("world!")

	check := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	w := NewTSDWriter(buf)
	check(w.Write(1, data1))
	check(w.Write(2, data2))
	check(w.Write(1, data2))
	check(w.Write(2, data2))

	r := NewTSDReader(buf)
	id, reader, err := r.Next()
	check(err)
	if id != 1 {
		t.Fail()
	}
	data, err := ioutil.ReadAll(reader)
	check(err)
	if !bytes.Equal(data, data1) {
		t.Fail()
	}
	id, reader, err = r.Next()
	check(err)
	if id != 2 {
		t.Fail()
	}
	data, err = ioutil.ReadAll(reader)
	check(err)
	if !bytes.Equal(data, data2) {
		t.Fail()
	}
	id, reader, err = r.Next()
	check(err)
	if id != 1 {
		t.Fail()
	}
	id, reader, err = r.Next()
	check(err)
	if id != 2 {
		t.Fail()
	}
	data, err = ioutil.ReadAll(reader)
	check(err)
	if !bytes.Equal(data, data2) {
		t.Fail()
	}
}