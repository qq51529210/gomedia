package mp4

import (
	"io"
)

const (
	TypeMDIA = 1835297121
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMDIA, DecodeBoxMDIA)
}

// MDIA 表示 mdia box
// 是一个容器
type MDIA struct {
	BasicBox
}

// DecodeBoxMDIA解析mdia box
func DecodeBoxMDIA(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器 box 解析子 box
	children, err := DecodeChildren(readSeeker, boxSize-headerSize)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(MDIA)
	box.size = boxSize
	box._type = _type
	box.BasicBox.children = children
	// 返回
	return box, nil
}
