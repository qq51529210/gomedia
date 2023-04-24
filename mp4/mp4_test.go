package mp4

import (
	"os"
	"testing"
)

func TestDecodeBox(t *testing.T) {
	file, err := os.Open("1.mp4")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	//
	mp4 := new(MP4)
	err = mp4.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	//
	testGetBox(t, mp4)
}

func testGetBox(t *testing.T, mp4 *MP4) {
	var box Box
	box = mp4.GetBox(TypeFTYP)
	if _, ok := box.(*FTYP); !ok {
		t.FailNow()
	}
	box = mp4.GetBox(TypeMOOV)
	if _, ok := box.(*MOOV); !ok {
		t.FailNow()
	}
}
