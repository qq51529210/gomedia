package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	TypeCO64 = 1668232756
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeCO64, DecodeBoxCO64)
}

// CO64 表示 co64 box
// 定义了每个 chunk 的偏移
type CO64 struct {
	BasicBox
	// 版本
	Version uint8
	// ...
	Flags uint32
	// 每个 chunk 的偏移
	ChunkOffset []uint64
}

// DecodeBoxCO64 解析 co64 box
func DecodeBoxCO64(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
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
	box := new(CO64)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.Uint24(buf[1:])
	// 4
	entryCount := binary.BigEndian.Uint32(buf[4:])
	n := 8
	if contentSize < int64(entryCount)*8+8 {
		return nil, errBoxSize
	}
	for i := 0; i < int(entryCount); i++ {
		box.ChunkOffset[i] = binary.BigEndian.Uint64(buf[n:])
		n += 8
	}
	// 返回
	return box, nil
}
