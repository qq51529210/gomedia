package mp4

import (
	"encoding/binary"
	"io"
)

const (
	// TypeAVC1 表示 avc1 类型
	TypeAVC1 Type = 1635148593
)

const (
	avc1BoxMinContentSize = 78
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeAVC1, DecodeBoxAVC1)
}

// AVC1 表示 avc1 box
type AVC1 struct {
	BasicBox
	// ...
	DataRefernce uint16
	// 宽度
	Width uint16
	// 高度
	Height uint16
	// 水平分辨率
	HResolution uint32
	// 垂直分辨率
	VResolution uint32
	// 每个 sample 的帧数
	FrameCount uint16
	// 色深
	Depth uint16
}

// DecodeBoxAVC1 解析 avc1 box
func DecodeBoxAVC1(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < avc1BoxMinContentSize {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, avc1BoxMinContentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(AVC1)
	box.size = boxSize
	box._type = _type
	// reserved 6
	// data reference id 2
	box.DataRefernce = binary.BigEndian.Uint16(buf[6:])
	// code stream version 4
	// reserved 12
	// 2
	box.Width = binary.BigEndian.Uint16(buf[24:])
	// 2
	box.Height = binary.BigEndian.Uint16(buf[26:])
	// 4
	box.HResolution = binary.BigEndian.Uint32(buf[28:])
	// 4
	box.VResolution = binary.BigEndian.Uint32(buf[32:])
	// reserved 4
	// 2
	box.FrameCount = binary.BigEndian.Uint16(buf[40:])
	// compressorname 32
	// 2
	box.Depth = binary.BigEndian.Uint16(buf[74:])
	// reserved 2
	contentSize -= avc1BoxMinContentSize
	for contentSize > 0 {
		// child
		child, err := DecodeBox(readSeeker)
		if err != nil {
			return nil, err
		}
		box.children = append(box.children, child)
		contentSize -= child.Size()
	}
	// 返回
	return box, nil
}
