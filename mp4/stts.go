package mp4

import (
	"encoding/binary"
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
	SampleCount uint32
	SampleDelta uint32
}

// STTS表示stts box
// 描述了sample时序的映射方法，通过它可以找到任何时间的sample。
// stts box可以包含一个压缩的表来映射时间和sample序号，
// 用其他的表来提供每个sample的长度和指针。
// 表中每个条目提供了在同一个时间偏移量里面连续的sample序号，以及samples的偏移量。
// 递增这些偏移量，就可以建立一个完整的time to sample表（时间戳到sample序号的映射表）
type STTS struct {
	BasicBox
	Entry []STTSEntry
}

// DecodeBoxSTTSP解析stts box
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
	box.BasicBox.size = boxSize
	box.BasicBox._type = _type
	// 4
	entryCount := binary.BigEndian.Uint32(buf)
	n := 4
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
