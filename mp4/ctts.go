package mp4

import (
	"encoding/binary"
	"io"
)

const (
	TypeCTTS = 1718909296
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeCTTS, DecodeBoxCTTS)
}

// CTTSEntry是CTTS的Entry字段
type CTTSEntry struct {
	SampleCount  uint32
	SampleOffset uint32
}

// CTTS表示ctts box
type CTTS struct {
	BasicBox
	Entry []CTTSEntry
}

// DecodeBoxCTTSP解析ctts box
func DecodeBoxCTTS(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < 4 {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	// 解析
	box := new(CTTS)
	box.BasicBox.size = boxSize
	box.BasicBox._type = _type
	// 4
	entryCount := binary.BigEndian.Uint32(buf)
	n := 4
	if contentSize < int64(entryCount*8+4) {
		return nil, errBoxSize
	}
	// 8*entryCount
	box.Entry = make([]CTTSEntry, entryCount)
	for i := 0; i < int(entryCount); i++ {
		box.Entry[i].SampleCount = binary.BigEndian.Uint32(buf[n:])
		n += 4
		box.Entry[i].SampleOffset = binary.BigEndian.Uint32(buf[n:])
		n += 4
	}
	// 返回
	return box, nil
}
