package mp4

import (
	"io"
)

const (
	// TypeTRAK 表示 trak 类型
	TypeTRAK Type = 1953653099
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeTRAK, DecodeBoxTRAK)
}

// TRAK 表示 trak box
// 是一个容器 box
type TRAK struct {
	BasicBox
}

// DecodeBoxTRAK 解析 trak box
func DecodeBoxTRAK(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器 box 解析子 box
	children, err := DecodeChildren(readSeeker, boxSize-headerSize)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(TRAK)
	box.size = boxSize
	box._type = _type
	box.BasicBox.children = children
	// 返回
	return box, nil
}
