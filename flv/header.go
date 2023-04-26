package flv

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	// FLVVersion表示版本
	FLVVersion = 1
)

var (
	errFormat = errors.New("error format")
)

// Header表示flv文件头
type Header struct {
	// 是否有视频
	HasVideo bool
	// 是否有音频
	HasAudio bool
	// 数据偏移
	DataOffset uint32
}

// Decode解码
func (h *Header) Decode(reader io.Reader) error {
	// 读取
	buf := make([]byte, 9)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return err
	}
	// FLV 24
	if buf[0] != 'F' && buf[1] != 'L' && buf[2] != 'V' {
		return errFormat
	}
	// version 8
	if buf[3] != FLVVersion {
		return errFormat
	}
	// flag 8
	switch buf[4] {
	case 0b00000100:
		h.HasAudio = true
	case 0b00000001:
		h.HasVideo = true
	case 0b00000101:
		h.HasAudio = true
		h.HasVideo = true
	default:
		return errFormat
	}
	// offset 32
	h.DataOffset = binary.BigEndian.Uint32(buf[5:])
	// 返回
	return nil
}

// Encode编码
func (h *Header) Encode(writer io.Writer) error {
	buf := make([]byte, 9)
	// FLV 24
	buf[0] = 'F'
	buf[1] = 'L'
	buf[2] = 'V'
	// version 8
	buf[3] = FLVVersion
	// flag 8
	if h.HasAudio {
		buf[4] |= 0b00000100
	}
	if h.HasVideo {
		buf[4] |= 0b00000001
	}
	// offset 32
	binary.BigEndian.PutUint32(buf[5:], h.DataOffset)
	// 写入
	_, err := writer.Write(buf)
	// 返回
	return err
}
