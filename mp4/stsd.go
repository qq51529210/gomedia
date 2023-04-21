package mp4

import "io"

const (
	TypeSTSD = 1937011556
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeSTSD, DecodeBoxSTSD)
}

// STSD表示stsd box
type STSD struct {
	BasicBox
}

// DecodeBoxSTSDP解析stsd box
func DecodeBoxSTSD(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器box解析子box
	children, err := DecodeChildren(readSeeker, boxSize-headerSize)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(STSD)
	box.BasicBox.size = boxSize
	box.BasicBox._type = _type
	box.BasicBox.children = children
	// 返回
	return box, nil
}
