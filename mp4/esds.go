package mp4

import (
	"gomedia/util"
	"io"
)

const (
	// TypeESDS 表示 esds 类型
	TypeESDS Type = 1702061171
)

const (
	esdsBoxMinContentSize = 29
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeESDS, DecodeBoxESDS)
}

// ESDS 表示 esds box
type ESDS struct {
	fullBox
	Decriptor []byte
}

// DecodeBoxESDS 解析 mdat box
func DecodeBoxESDS(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 略过
	contentSize := boxSize - headerSize
	if contentSize < esdsBoxMinContentSize {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(ESDS)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.Uint24(buf[1:])
	// n
	box.Decriptor = buf[4:]
	// 返回
	return box, nil
}
