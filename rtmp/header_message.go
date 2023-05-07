package rtmp

import (
	"encoding/binary"
	"fmt"
	"gomedia/util"
	"io"
)

const (
	messageHeaderLen  = 11
	messageHeaderLen0 = 15
	messageHeaderLen1 = 7
	messageHeaderLen2 = 3
)

// MessageHeader 用于解析消息头
type MessageHeader struct {
	FMT FMT
	// 时间戳, 3 个字节
	// 2 类型表示 delta
	// 3 类型没有
	Timestamp uint32
	// 大小, 3 个字节
	// 2,3 类型没有
	Length uint32
	// 类型, 1 个字节
	// 2,3 类型没有
	Type byte
	// 流, 4 个字节
	// 1,2,3 类型没有
	StreamID uint32
}

// Encode 编码
func (h *MessageHeader) Encode(writer io.Writer) error {
	switch h.FMT {
	case BasicHeaderFMT0:
		return h.encode0(writer)
	case BasicHeaderFMT1:
		return h.encode1(writer)
	case BasicHeaderFMT2:
		return h.encode2(writer)
	default:
		return nil
	}
}

func (h *MessageHeader) encode0(writer io.Writer) error {
	buf := make([]byte, messageHeaderLen0)
	n := 11
	// timestamp
	if h.Timestamp < 0xffffff {
		util.PutBigUint24(buf, h.Timestamp)
	} else {
		buf[0] = 0xff
		buf[1] = 0xff
		buf[2] = 0xff
		binary.BigEndian.PutUint32(buf[11:], h.Timestamp)
		n = 15
	}
	//  length
	util.PutBigUint24(buf[3:], h.Length)
	// type
	buf[6] = h.Type
	// stream id
	binary.BigEndian.PutUint32(buf[7:], h.StreamID)
	// 写入
	_, err := writer.Write(buf[:n])
	// 返回
	return err
}

func (h *MessageHeader) encode1(writer io.Writer) error {
	buf := make([]byte, messageHeaderLen1)
	// timestamp delta
	util.PutBigUint24(buf, h.Timestamp)
	//  length
	util.PutBigUint24(buf[3:], h.Length)
	// type
	buf[6] = h.Type
	// 写入
	_, err := writer.Write(buf[:7])
	// 返回
	return err
}

func (h *MessageHeader) encode2(writer io.Writer) error {
	buf := make([]byte, messageHeaderLen2)
	// timestamp delta
	util.PutBigUint24(buf, h.Timestamp)
	// 写入
	_, err := writer.Write(buf[:3])
	// 返回
	return err
}

// Decode 解码
func (h *MessageHeader) Decode(reader io.Reader) error {
	switch h.FMT {
	case BasicHeaderFMT0:
		return h.decode0(reader)
	case BasicHeaderFMT1:
		return h.decode1(reader)
	case BasicHeaderFMT2:
		return h.decode2(reader)
	default:
		return nil
	}
}

func (h *MessageHeader) decodeExtendTimestamp(reader io.Reader) error {
	if h.Timestamp == 0xffffff {
		// 再读取
		buf := make([]byte, 4)
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			return fmt.Errorf("decode fmt 0 extended timestamp error: %w", err)
		}
		h.Timestamp = binary.BigEndian.Uint32(buf)
	}
	// 返回
	return nil
}

func (h *MessageHeader) decode0(reader io.Reader) error {
	buf := make([]byte, messageHeaderLen)
	// 读取
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return fmt.Errorf("decode fmt 0 error: %w", err)
	}
	// timestamp
	h.Timestamp = util.BigUint24(buf)
	// length
	h.Length = util.BigUint24(buf[3:])
	// type
	h.Type = buf[6]
	// stream id
	h.StreamID = binary.BigEndian.Uint32(buf[7:])
	// extended timestamp
	return h.decodeExtendTimestamp(reader)
}

func (h *MessageHeader) decode1(reader io.Reader) error {
	buf := make([]byte, messageHeaderLen1)
	// 读取
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return fmt.Errorf("decode fmt 0 error: %w", err)
	}
	// timestamp
	h.Timestamp = util.BigUint24(buf)
	// length
	h.Length = util.BigUint24(buf[3:])
	// type
	h.Type = buf[6]
	// extended timestamp
	return h.decodeExtendTimestamp(reader)
}

func (h *MessageHeader) decode2(reader io.Reader) error {
	buf := make([]byte, messageHeaderLen2)
	// 读取
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return fmt.Errorf("decode fmt 0 error: %w", err)
	}
	// timestamp
	h.Timestamp = util.BigUint24(buf)
	// extended timestamp
	return h.decodeExtendTimestamp(reader)
}
