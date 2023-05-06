package mp4

import (
	"fmt"
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

func TestType(t *testing.T) {
	fmt.Println(uint32(TypeOfName("esds")))

	if TypeAVC1 != TypeOfName("avc1") {
		t.FailNow()
	}
	if TypeAVCC != TypeOfName("avcC") {
		t.FailNow()
	}
	if TypeCO64 != TypeOfName("co64") {
		t.FailNow()
	}
	if TypeCTTS != TypeOfName("ctts") {
		t.FailNow()
	}
	if TypeFTYP != TypeOfName("ftyp") {
		t.FailNow()
	}
	if TypeHDLR != TypeOfName("hdlr") {
		t.FailNow()
	}
	if TypeMDAT != TypeOfName("mdat") {
		t.FailNow()
	}
	if TypeMDHD != TypeOfName("mdhd") {
		t.FailNow()
	}
	if TypeMDIA != TypeOfName("mdia") {
		t.FailNow()
	}
	if TypeMINF != TypeOfName("minf") {
		t.FailNow()
	}
	if TypeMOOV != TypeOfName("moov") {
		t.FailNow()
	}
	if TypeMP4A != TypeOfName("mp4a") {
		t.FailNow()
	}
	if TypeMVHD != TypeOfName("mvhd") {
		t.FailNow()
	}
	if TypeSMHD != TypeOfName("smhd") {
		t.FailNow()
	}
	if TypeSTBL != TypeOfName("stbl") {
		t.FailNow()
	}
	if TypeSTCO != TypeOfName("stco") {
		t.FailNow()
	}
	if TypeSTSC != TypeOfName("stsc") {
		t.FailNow()
	}
	if TypeSTSD != TypeOfName("stsd") {
		t.FailNow()
	}
	if TypeSTSS != TypeOfName("stss") {
		t.FailNow()
	}
	if TypeSTSZ != TypeOfName("stsz") {
		t.FailNow()
	}
	if TypeSTTS != TypeOfName("stts") {
		t.FailNow()
	}
	if TypeTKHD != TypeOfName("tkhd") {
		t.FailNow()
	}
	if TypeTRAK != TypeOfName("trak") {
		t.FailNow()
	}
	if TypeVMHD != TypeOfName("vmhd") {
		t.FailNow()
	}
}
