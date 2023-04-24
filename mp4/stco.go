package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	TypeSTCO = 1937007471
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeSTCO, DecodeBoxSTCO)
}

// STCO表示stco box
// 定义了每个chunk的偏移
type STCO struct {
	fullBox
	// 每个chunk的偏移
	ChunkOffset []uint32
}

// DecodeBoxSTCO解析stco box
func DecodeBoxSTCO(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
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
	box := new(STCO)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.Uint24(buf[1:])
	// 4
	entryCount := binary.BigEndian.Uint32(buf[4:])
	n := 8
	if contentSize < int64(entryCount)*4+8 {
		return nil, errBoxSize
	}
	box.ChunkOffset = make([]uint32, entryCount)
	for i := 0; i < len(box.ChunkOffset); i++ {
		box.ChunkOffset[i] = binary.BigEndian.Uint32(buf[n:])
		n += 4
	}
	// 返回
	return box, nil
}
