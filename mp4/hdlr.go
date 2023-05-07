package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	// TypeHDLR 表示 hdlr 类型
	TypeHDLR Type = 1751411826
)

// 数据类型
const (
	HDLRTypeVide Type = 1986618469
	HDLRTypeSoun Type = 1936684398
	HDLRTypeHint Type = 1751740020
	HDLRTypeMeta Type = 1835365473
	HDLRTypeAuxv Type = 1635088502
)

const (
	hdlrBoxMinContentSize = 24
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeHDLR, DecodeBoxHDLR)
}

// HDLR 表示 hdlr box
// 主要用于判断是什么类型的 track
type HDLR struct {
	fullBox
	// ...
	PreDefined uint32
	// 类型
	HandlerType uint32
	// 名称
	Name []byte
}

// DecodeBoxHDLR 解析 hdlr box
func DecodeBoxHDLR(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < hdlrBoxMinContentSize {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	box := new(HDLR)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.BigUint24(buf[1:])
	// 4
	box.PreDefined = binary.BigEndian.Uint32(buf[4:])
	// 4
	box.HandlerType = binary.BigEndian.Uint32(buf[8:])
	// Reserved 4*3
	// box.Reserved[0] = binary.BigEndian.Uint32(buf[12:])
	// box.Reserved[1] = binary.BigEndian.Uint32(buf[16:])
	// box.Reserved[2] = binary.BigEndian.Uint32(buf[20:])
	// n
	box.Name = append(box.Name, buf[24:]...)
	// 返回
	return box, nil
}
