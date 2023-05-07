package flv

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	// TagHeaderLen表示tag头的字节
	TagHeaderLen = 11
)

// Tag表示flv的body tag头
type Tag struct {
	// 类型，音频（0x08），视频（0x09），脚本（0x12）
	Type byte
	// 数据大小
	DataSize uint32
	// 时间戳
	Timestamp uint32
}

// Decode解码
func (h *Tag) Decode(reader io.Reader) error {
	// 读取
	buf := make([]byte, TagHeaderLen)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return err
	}
	// type 8
	h.Type = buf[0]
	// data size 24
	h.DataSize = util.BigUint24(buf[1:])
	// timestamp 32
	// 先把最高位放好，直接解析4个字节
	buf[3] = buf[7]
	h.Timestamp = binary.BigEndian.Uint32(buf[3:])
	// stream id为0就算了
	// 返回
	return nil
}

// Encode编码
func (h *Tag) Encode(writer io.Writer) error {
	buf := make([]byte, TagHeaderLen)
	// type
	buf[0] = h.Type
	// data size
	util.PutUint24(buf[1:], h.DataSize)
	// timestamp
	util.PutUint24(buf[4:], h.Timestamp)
	// timestamp extend
	buf[7] = byte(h.Timestamp >> 24)
	// stream id
	// 写入
	_, err := writer.Write(buf)
	// 返回
	return err
}

// Size返回整个tag的大小
func (h *Tag) Size() uint32 {
	return TagHeaderLen + h.DataSize
}
