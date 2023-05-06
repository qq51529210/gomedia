package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	TypeVMHD = 1986881636
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeVMHD, DecodeBoxVMHD)
}

// VMHD 表示 vmhd box
// 包含当前 trac k的视频描述信息
type VMHD struct {
	fullBox
	// 视频合成模式, 为 0 时拷贝原始图像,
	// 否则与 opcolor 进行合成
	GraphicsMode uint16
	// ｛red，green，blue｝
	OpColor [3]uint16
}

// DecodeBoxVMHD 解析 vmhd box
func DecodeBoxVMHD(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < 12 {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	box := new(VMHD)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.Uint24(buf[1:])
	// 2
	box.GraphicsMode = binary.BigEndian.Uint16(buf[4:])
	// 6
	box.OpColor[0] = binary.BigEndian.Uint16(buf[6:])
	box.OpColor[1] = binary.BigEndian.Uint16(buf[8:])
	box.OpColor[2] = binary.BigEndian.Uint16(buf[10:])
	// 返回
	return box, nil
}
