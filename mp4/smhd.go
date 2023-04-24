package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	TypeSMHD = 1936549988
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeSMHD, DecodeBoxSMHD)
}

// SMHD表示smhd box
// 包含当前track的音频描述信息，如编码格式等信息
type SMHD struct {
	fullBox
	// 立体声平衡，[8.8] 格式值，一般为0，
	// -1.0表示全部左声道，1.0表示全部右声道
	Balance uint16
}

// DecodeBoxSMHD解析smhd box
func DecodeBoxSMHD(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
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
	box := new(SMHD)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.Uint24(buf[1:])
	// 2
	box.Balance = binary.BigEndian.Uint16(buf[4:])
	// 2 reserved
	// 返回
	return box, nil
}
