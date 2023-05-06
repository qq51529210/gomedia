package mp4

import "io"

const (
	TypeSTBL = 1937007212
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeSTBL, DecodeBoxSTBL)
}

// STBL 表示 stbl box
// 是一个容器
type STBL struct {
	BasicBox
}

// DecodeBoxSTBL 解析 stbl box
func DecodeBoxSTBL(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器 box 解析子 box
	children, err := DecodeChildren(readSeeker, boxSize-headerSize)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(STBL)
	box.size = boxSize
	box._type = _type
	box.BasicBox.children = children
	// 返回
	return box, nil
}
