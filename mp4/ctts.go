package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	// TypeCTTS 表示 ctts 类型
	TypeCTTS Type = 1668576371
)

const (
	cttsBoxMinContentSize = 8
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeCTTS, DecodeBoxCTTS)
}

// CTTSEntry 是 CTTS 的 Entry 字段
// 主要是用于计算pts
type CTTSEntry struct {
	SampleCount  uint32
	SampleOffset uint32
}

// CTTS 表示 ctts box
type CTTS struct {
	fullBox
	// 元素
	Entry []CTTSEntry
}

// DecodeBoxCTTS 解析 ctts box
func DecodeBoxCTTS(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < cttsBoxMinContentSize {
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
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.BigUint24(buf[1:])
	// 4
	entryCount := binary.BigEndian.Uint32(buf[4:])
	if contentSize < int64(entryCount)*8+cttsBoxMinContentSize {
		return nil, errBoxSize
	}
	// 8*entryCount
	n := cttsBoxMinContentSize
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
