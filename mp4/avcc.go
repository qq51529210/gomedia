package mp4

import (
	"encoding/binary"
	"io"
)

const (
	TypeAVCC = 1635148611
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeAVCC, DecodeBoxAVCC)
}

// AVCC表示avcC box
type AVCC struct {
	BasicBox
	// H264
	AVCProfileIndication byte
	// H264
	ProfileCompatibility byte
	// H264
	AVCLevelIndication byte
	// H264
	LengthSizeMinusOne byte
	// SPS
	SPS [][]byte
	// PPS
	PPS [][]byte
}

// DecodeBoxAVCC解析avcC box
func DecodeBoxAVCC(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < 7 {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(AVCC)
	box.size = boxSize
	box._type = _type
	// configuration version 1
	// 1
	box.AVCProfileIndication = buf[1]
	// 1
	box.ProfileCompatibility = buf[2]
	// 1
	box.AVCLevelIndication = buf[3]
	// 1
	box.LengthSizeMinusOne = buf[4] & 0b00000011
	// number of sps 1
	numberOfSPS := buf[5] & 0b00011111
	// sps
	n := 6
	if contentSize < int64(n)+int64(numberOfSPS)*2 {
		return nil, errBoxSize
	}
	for i := byte(0); i < numberOfSPS; i++ {
		// 2
		spsLen := binary.BigEndian.Uint16(buf[n:])
		n += 2
		if contentSize < int64(n+int(spsLen)) {
			return nil, errBoxSize
		}
		// n
		sps := make([]byte, spsLen)
		n += copy(sps, buf[n:])
		box.SPS = append(box.SPS, sps)
	}
	// number of pps
	numberOfPPS := buf[n] & 0b00011111
	n++
	if contentSize < int64(n)+int64(numberOfPPS)*2 {
		return nil, errBoxSize
	}
	// pps
	for i := byte(0); i < numberOfPPS; i++ {
		// 2
		ppsLen := binary.BigEndian.Uint16(buf[n:])
		n += 2
		// n
		pps := make([]byte, ppsLen)
		n += copy(pps, buf[n:])
		box.PPS = append(box.PPS, pps)
	}
	// 剩下的
	contentSize -= int64(n)
	if contentSize > 0 {
		_, err = readSeeker.Seek(contentSize, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}
	// 返回
	return box, nil
}
