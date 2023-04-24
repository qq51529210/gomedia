package mp4

import (
	"encoding/binary"
	"io"
)

const (
	TypeMP4A = 1836069985
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMP4A, DecodeBoxMP4A)
}

// MP4A表示mp4a box
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

// DecodeBoxMP4A解析mp4a box
func DecodeBoxMP4A(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < 26 {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, 26)
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
	// pre_defined 2
	// reserved 2
	// 4
	box.SampleRate = binary.BigEndian.Uint32(buf[22:])
	//
	contentSize -= 26
	if contentSize > 0 {
		// todo解析
		_, err = readSeeker.Seek(contentSize, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}
	// 返回
	return box, nil
}
