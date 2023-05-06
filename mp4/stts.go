package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	TypeSTTS = 1937011827
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeSTTS, DecodeBoxSTTS)
}

// STTSEntry是STTS的Entry字段
type STTSEntry struct {
	// 个数
	SampleCount uint32
	// 这些 sample 的时间
	SampleDelta uint32
}

// STTS 表示 stts box
// 描述了sample时序的映射方法,
// 通过它可以找到任何时间的sample
type STTS struct {
	fullBox
	// ...
	Entry []STTSEntry
}

// DecodeBoxSTTS解析stts box
func DecodeBoxSTTS(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
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
	box := new(STTS)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.Uint24(buf[1:])
	// 4
	entryCount := binary.BigEndian.Uint32(buf[4:])
	n := 8
	if contentSize < int64(entryCount*8+4) {
		return nil, errBoxSize
	}
	// 8*entryCount
	box.Entry = make([]STTSEntry, entryCount)
	for i := 0; i < int(entryCount); i++ {
		box.Entry[i].SampleCount = binary.BigEndian.Uint32(buf[n:])
		n += 4
		box.Entry[i].SampleDelta = binary.BigEndian.Uint32(buf[n:])
		n += 4
	}
	// 返回
	return box, nil
}
