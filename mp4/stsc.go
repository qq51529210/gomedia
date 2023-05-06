package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	TypeSTSC = 1937011555
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeSTSC, DecodeBoxSTSC)
}

// STSCEntry 是 STS C的 Entry 字段
type STSCEntry struct {
	// chunk 开始的偏移
	FirstChunk uint32
	// chunk 有多少个 sample
	SamplePerChunk uint32
	// sample 的描述, 默认设置为1
	SampleDescriptionIndex uint32
}

// STSC 表示 stsc box
// 表明每一个 chunk 中有多少个 sample
type STSC struct {
	fullBox
	// 元素
	Entry []STSCEntry
}

// DecodeBoxSTSC 解析 stsc box
func DecodeBoxSTSC(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < 8 {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	// 解析
	box := new(STSC)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.Uint24(buf[1:])
	// 4
	entryCount := binary.BigEndian.Uint32(buf[4:])
	n := 8
	if contentSize < int64(entryCount)*12+8 {
		return nil, errBoxSize
	}
	// 12*entryCount
	box.Entry = make([]STSCEntry, entryCount)
	for i := 0; i < int(entryCount); i++ {
		box.Entry[i].FirstChunk = binary.BigEndian.Uint32(buf[n:])
		n += 4
		box.Entry[i].SamplePerChunk = binary.BigEndian.Uint32(buf[n:])
		n += 4
		box.Entry[i].SampleDescriptionIndex = binary.BigEndian.Uint32(buf[n:])
		n += 4
	}
	// 返回
	return box, nil
}
