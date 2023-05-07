package rtmp

import (
	"fmt"
	"io"
)

type FMT byte

// fmt 的值
const (
	BasicHeaderFMT0 FMT = 0b00000000
	BasicHeaderFMT1 FMT = 0b01000000
	BasicHeaderFMT2 FMT = 0b10000000
	BasicHeaderFMT3 FMT = 0b11000000
)

// 用于判断 csid
const (
	basicHeaderCSID0 byte = 0b00000000
	basicHeaderCSID1 byte = 0b00111111
)

const (
	basicHeaderLen = 3
)

// BasicHeader 用于解析基本消息头
type BasicHeader struct {
	// BasicHeaderFMT0 表示 Message Header 长度为 11
	// BasicHeaderFMT1 表示 Message Header 长度为 7
	// BasicHeaderFMT2 表示 Message Header 长度为 3
	// BasicHeaderFMT3 表示 Message Header 长度为 0
	FMT FMT
	// chunk stream id
	// 0 表示 Basic Header 总共要占用 2 个字节
	// 1 表示 Basic Header 总共要占用 3 个字节
	// 2 代表该 chunk 是控制信息和一些命令信息
	// 3 代表该 chunk 是客户端发出的 AMF0 命令以及服务端对该命令的应答
	// 4 代表该 chunk 是客户端发出的音频数据，用于 publish
	// 5 代表该 chunk 是服务端发出的 AMF0 命令和数据
	// 6 代表该 chunk 是服务端发出的音频数据，用于 play；或客户端发出的视频数据，用于 publish
	// 7 代表该 chunk 是服务端发出的视频数据，用于 play
	// 8 代表该 chunk 是客户端发出的 AMF0 命令，专用来发送： getStreamLength, play, publish
	CSID uint32
}

// Encode 编码
func (h *BasicHeader) Encode(writer io.Writer) error {
	buf := make([]byte, basicHeaderLen)
	buf[0] = byte(h.FMT)
	n := 0
	if h.CSID < 64 {
		// 1 字节
		buf[0] |= byte(h.CSID)
		n = 1
	} else if h.CSID < 320 {
		// 2 字节
		buf[0] |= basicHeaderCSID0
		buf[1] = byte(h.CSID - 64)
		n = 2
	} else {
		// 3 字节
		buf[0] |= basicHeaderCSID1
		csid := h.CSID - 64
		buf[1] = byte(csid)
		buf[2] = byte(csid >> 8)
		n = 3
	}
	_, err := writer.Write(buf[:n])
	// 返回
	return err
}

// Decode 解码
func (h *BasicHeader) Decode(reader io.Reader) error {
	buf := make([]byte, basicHeaderLen)
	// 第一个字节
	_, err := reader.Read(buf[:1])
	if err != nil {
		return fmt.Errorf("decode fmt and csid error: %w", err)
	}
	// fmt
	h.FMT = FMT(buf[0]) & BasicHeaderFMT3
	csid := buf[0] & basicHeaderCSID1
	switch csid {
	case basicHeaderCSID0:
		// 2 字节
		_, err := reader.Read(buf[:1])
		if err != nil {
			return fmt.Errorf("decode 2 bytes csid error: %w", err)
		}
		h.CSID = uint32(buf[0]) - 64
	case basicHeaderCSID1:
		// 3 字节
		_, err := reader.Read(buf[:2])
		if err != nil {
			return fmt.Errorf("decode 3 bytes csid error: %w", err)
		}
		h.CSID = uint32(buf[0]) + uint32(buf[1])*256 + 64
	default:
		h.CSID = uint32(csid)
	}
	// 返回
	return nil
}

// Len 返回字节数
func (h *BasicHeader) Len() int {
	if h.CSID < 64 {
		return 1
	}
	if h.CSID < 320 {
		return 2
	}
	return 3
}
