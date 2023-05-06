package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	TypeSTSD = 1937011556
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeSTSD, DecodeBoxSTSD)
}

// STSD 表示 stsd box
type STSD struct {
	fullBox
	//
	EntryCount uint32
}

// DecodeBoxSTSD 解析 stsd box
func DecodeBoxSTSD(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < 8 {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, 8)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(STSD)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.Uint24(buf[1:])
	// 4
	box.EntryCount = binary.BigEndian.Uint32(buf[4:])
	for i := 0; i < int(box.EntryCount); i++ {
		child, err := DecodeBox(readSeeker)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		box.children = append(box.children, child)
	}
	// 返回
	return box, nil
}
