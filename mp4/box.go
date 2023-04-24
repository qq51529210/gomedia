package mp4

import (
	"encoding/binary"
	"errors"
	"io"
)

// 错误
var (
	errBoxSize = errors.New("error box size")
)

// 其他解析函数
var (
	decodeFuncs = make(map[Type]DecodeFunc)
)

// DecodeFunc表示解码函数
type DecodeFunc func(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error)

// Type表示box的类型
type Type uint32

func (t Type) String() string {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(t))
	return string(buf)
}

// TypeOfName返回name的Type值，4个字节
func TypeOfName(name string) Type {
	var buf [4]byte
	copy(buf[0:], name)
	return Type(binary.BigEndian.Uint32(buf[0:]))
}

// Box表示一个mp4的box接口
type Box interface {
	// SetType(_type uint32)
	// SetSize(size int64) error
	Type() Type
	Size() int64
	// 所有子box
	Children() []Box
	// 添加子box
	AddChild(box Box)
	// 获取子box
	GetChild(_type Type) Box
}

// BasicBox表示简单的box
type BasicBox struct {
	size     int64
	_type    Type
	children []Box
}

// SetType设置类型
func (b *BasicBox) SetType(_type Type) {
	b._type = _type
}

// SetSize设置大小，内部不检查参数
func (b *BasicBox) SetSize(size int64) {
	b.size = size
}

// Type实现Box接口
func (b *BasicBox) Type() Type {
	return b._type
}

// Size实现Box接口
func (b *BasicBox) Size() int64 {
	return b.size
}

// Children实现Box接口
func (b *BasicBox) Children() []Box {
	return b.children
}

// AddChild实现Box接口
func (b *BasicBox) AddChild(box Box) {
	b.children = append(b.children, box)
}

// GetChild实现Box接口
func (b *BasicBox) GetChild(_type Type) Box {
	// 查找
	for _, child := range b.children {
		// 本节点的子节点
		if child.Type() == _type {
			return child
		}
		// 子节点继续
		box := child.GetChild(_type)
		if box != nil {
			return box
		}
	}
	// 没有
	return nil
}

// 添加子定义的解析器，如果_type不正确或者重复返回错误
func AddDecodeFunc(_type Type, _func DecodeFunc) {
	// 添加
	decodeFuncs[_type] = _func
}

// DecodeBox从reader解析Box的size和type，然后调用注册的decodeFunc来解析内容
func DecodeBox(readSeeker io.ReadSeeker) (Box, error) {
	// 读取
	buf := make([]byte, 8)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	// 解析基础
	size := int64(binary.BigEndian.Uint32(buf))
	_type := Type(binary.BigEndian.Uint32(buf[4:]))
	switch {
	case size == 0:
		// 没有size，最后一个box
		return &BasicBox{
			size:  size,
			_type: _type,
		}, nil
	case size == 1:
		// 再读取8字节的长度
		_, err := io.ReadFull(readSeeker, buf)
		if err != nil {
			return nil, err
		}
		size = int64(binary.BigEndian.Uint64(buf))
		// 解析器
		decoder := decodeFuncs[_type]
		if decoder == nil {
			// 没有解析器
			return DecodeUnknownBox(readSeeker, 16, size, _type)
		}
		// 交由解析器解析
		return decoder(readSeeker, 16, size, _type)
	default:
		// 解析器
		decoder := decodeFuncs[_type]
		if decoder == nil {
			// 没有解析器
			return DecodeUnknownBox(readSeeker, 8, size, _type)
		}
		// 交由解析器解析
		return decoder(readSeeker, 8, size, _type)
	}
}

// DecodeUnknownBox不解析，直接seek到box结尾，然后返回
func DecodeUnknownBox(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	contentSize := boxSize - headerSize
	if contentSize > 0 {
		// 略过
		_, err := readSeeker.Seek(contentSize, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}
	// 返回
	return &BasicBox{
		size:  boxSize,
		_type: _type,
	}, nil
}

// DecodeChildren一般用于容器box解析子box，一直解析到contentSize为0
func DecodeChildren(readSeeker io.ReadSeeker, contentSize int64) ([]Box, error) {
	var children []Box
	// 循环解析所有的子box即可
	for contentSize > 0 {
		// 解析
		box, err := DecodeBox(readSeeker)
		if err != nil {
			if err == io.EOF {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}
		children = append(children, box)
		// 看看boxSize用完没有
		contentSize -= box.Size()
	}
	// 返回
	return children, nil
}
