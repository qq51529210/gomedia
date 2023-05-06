package mp4

import (
	"encoding/binary"
	"io"
)

const (
	// TypeMP4A 表示 mp4a 类型
	TypeMP4A Type = 1836069985
)

const (
	mp4aBoxMinContentSize = 28
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMP4A, DecodeBoxMP4A)
}

// MP4A 表示 mp4a box
type MP4A struct {
	BasicBox
	// ...
	DataRefernce uint16
	// 声道数量
	ChannelCount uint16
	// 采样位宽
	SampleSize uint16
	// 采样率
	SampleRate uint32
}

// DecodeBoxMP4A 解析 mp4a box
func DecodeBoxMP4A(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < mp4aBoxMinContentSize {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, mp4aBoxMinContentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(MP4A)
	box.size = boxSize
	box._type = _type
	// reserved 6
	// 2
	box.DataRefernce = binary.BigEndian.Uint16(buf[6:])
	// reserved 8
	// 2
	box.ChannelCount = binary.BigEndian.Uint16(buf[16:])
	// 2
	box.SampleSize = binary.BigEndian.Uint16(buf[18:])
	// reserved 4
	// 4
	box.SampleRate = binary.BigEndian.Uint32(buf[24:])
	//
	contentSize -= mp4aBoxMinContentSize
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
