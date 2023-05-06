package mp4

import "io"

const (
	TypeMINF = 1835626086
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMINF, DecodeBoxMINF)
}

// MINF 表示 minf box
// 是一个容器
type MINF struct {
	BasicBox
}

// DecodeBoxMINF 解析 minf box
func DecodeBoxMINF(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器 box 解析子 box
	children, err := DecodeChildren(readSeeker, boxSize-headerSize)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(MINF)
	box.size = boxSize
	box._type = _type
	box.BasicBox.children = children
	// 返回
	return box, nil
}
