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

// DecodeFunc 表示解码函数
type DecodeFunc func(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error)

// Type 表示 box 的类型
type Type uint32

func (t Type) String() string {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(t))
	return string(buf)
}

// TypeOfName 返回 name 的 Type
func TypeOfName(name string) Type {
	var buf [4]byte
	copy(buf[0:], name)
	return Type(binary.BigEndian.Uint32(buf[0:]))
}

// Box 表示一个 mp4 的 box 接口
type Box interface {
	// 类型
	Type() Type
	// 大小
	Size() int64
	// 获取子 box
	Children() []Box
	// 添加子 box
	AddChild(box Box)
	// 获取指定类型的子 box
	GetChild(_type Type) []Box
}

// BasicBox 表示简单的 box
type BasicBox struct {
	size     int64
	_type    Type
	children []Box
}

// Type 实现 Box 接口
func (b *BasicBox) Type() Type {
	return b._type
}

// Size 实现 Box 接口
func (b *BasicBox) Size() int64 {
	return b.size
}

// Children 实现 Box 接口
func (b *BasicBox) Children() []Box {
	return b.children
}

// AddChild 实现 Box 接口
func (b *BasicBox) AddChild(box Box) {
	b.children = append(b.children, box)
}

// GetChild 实现 Box 接口
func (b *BasicBox) GetChild(_type Type) (boxs []Box) {
	// 查找
	for _, child := range b.children {
		// 本节点的子节点
		if child.Type() == _type {
			boxs = append(boxs, child)
		}
	}
	// 没有
	return
}

// fullBox 表示 full box
type fullBox struct {
	BasicBox
	// 版本
	Version uint8
	// ...
	Flags uint32
}

// 添加子定义的解析器,
// _func 必须 Read 或者 Seek 相应的字节,
// 否则下一个box解析size会出错
func AddDecodeFunc(_type Type, _func DecodeFunc) {
	// 添加
	decodeFuncs[_type] = _func
}

// DecodeBox 从 reader 解析 Box 的 size 和 type,
// 然后调用已经注册的 DecodeFunc 来解析内容,
// 如果没有找到相应的 type , 那么使用 DecodeUnknownBox
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
		// 没有 size 最后一个 box
		return &BasicBox{
			size:  size,
			_type: _type,
		}, nil
	case size == 1:
		// 再读取 8 字节的长度
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

// DecodeUnknownBox 不解析，直接 seek 到 box 结尾，然后返回
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

// DecodeChildren 一般用于容器 box 解析子 box , 一直解析到 contentSize 为 0
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
