package flv

import (
	"io"
	"os"
	"testing"
)

func Test_FLV(t *testing.T) {
	file1, err := os.Open("1.flv")
	if err != nil {
		t.Fatal(err)
	}
	defer file1.Close()
	file2, err := os.OpenFile("2.flv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer file2.Close()
	//
	var src, dst Header
	err = src.Decode(file1)
	if err != nil {
		t.FailNow()
	}
	if !src.HasAudio || !src.HasVideo || src.DataOffset != 9 {
		t.FailNow()
	}
	err = dst.Encode(file2)
	if err != nil {
		t.FailNow()
	}
	//
	var th Tag
	prevSize := make([]byte, 4)
	for {
		_, err = io.ReadFull(file1, prevSize)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		err = th.Decode(file1)
		if err != nil {
			if err == io.EOF {
				_, err = file2.Write(prevSize)
				if err != nil {
					t.FailNow()
				}
				break
			}
		}
		//
		_, err = file2.Write(prevSize)
		if err != nil {
			t.FailNow()
		}
		err = th.Encode(file2)
		if err != nil {
			t.FailNow()
		}
		_, err = io.CopyN(file2, file1, int64(th.DataSize))
		if err != nil {
			t.FailNow()
		}
	}
}
