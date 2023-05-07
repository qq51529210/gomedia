package util

import (
	"bytes"
	"testing"
)

func Test_CRLFReader(t *testing.T) {
	r := NewCRLFReader(bytes.NewBufferString("123\n456\r\n789\r\n\r\n"))
	line, err := r.ReadLine()
	if err != nil {
		t.Fatal(err)
	}
	if line != "123\n456" {
		t.FailNow()
	}
	line, err = r.ReadLine()
	if err != nil {
		t.Fatal(err)
	}
	if line != "789" {
		t.FailNow()
	}
	line, err = r.ReadLine()
	if err != nil {
		t.Fatal(err)
	}
	if line != "" {
		t.FailNow()
	}
}
