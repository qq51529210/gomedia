package mp4

import "io"

const (
	TypeSTBL = 1937007212
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeSTBL, DecodeBoxSTBL)
}

// STBL表示stbl box
// 包含了关于track中sample所有时间和位置的信息，以及sample的编解码等信息
type STBL struct {
	BasicBox
}

// DecodeBoxSTBLP解析stbl box
func DecodeBoxSTBL(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器box解析子box
	children, err := DecodeChildren(readSeeker, boxSize-headerSize)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(STBL)
	box.BasicBox.size = boxSize
	box.BasicBox._type = _type
	box.BasicBox.children = children
	// 返回
	return box, nil
}
