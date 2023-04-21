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

// MDIA表示mdia box
// 包含类整个track的媒体信息，比如媒体类型和sample信息
type MDIA struct {
	BasicBox
}

// DecodeBoxMDIA解析mdia box
func DecodeBoxMDIA(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器box解析子box
	children, err := DecodeChildren(readSeeker, boxSize-headerSize)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(MDIA)
	box.BasicBox.size = boxSize
	box.BasicBox._type = _type
	box.BasicBox.children = children
	// 返回
	return box, nil
}
