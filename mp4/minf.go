package mp4

import "io"

const (
	TypeMINF = 1835626086
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMINF, DecodeBoxMINF)
}

// MINF表示minf box
// 包含了所有描述该track中的媒体信息的对象，信息存储在其子box中
type MINF struct {
	BasicBox
}

// DecodeBoxMINF解析minf box
func DecodeBoxMINF(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器box解析子box
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
