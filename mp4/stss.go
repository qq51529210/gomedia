package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	TypeSTSS = 1937011571
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeSTSS, DecodeBoxSTSS)
}

// STSSEntry 是 STSS 的 Entry字段
type STSSEntry struct {
	SampleCount uint32
	SampleDelta uint32
}

// STSS 表示 stts box
// 确定media中的关键帧,
// 如果此表不存在,说明每一 个sample 都是一个关键帧
type STSS struct {
	fullBox
	SampleNumber []uint32
}

// DecodeBoxSTSS 解析 stts box
func DecodeBoxSTSS(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
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
	box := new(STSS)
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
	// 4*entryCount
	box.SampleNumber = make([]uint32, entryCount)
	for i := 0; i < int(entryCount); i++ {
		box.SampleNumber[i] = binary.BigEndian.Uint32(buf[n:])
		n += 4
	}
	// 返回
	return box, nil
}
