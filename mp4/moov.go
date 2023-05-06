package mp4

import "io"

const (
	// TypeMOOV 表示 moov 类型
	TypeMOOV Type = 1836019574
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMOOV, DecodeBoxMOOV)
}

// MOOV 表示 moov box
// 有且只有一个并且包含在文件层
type MOOV struct {
	BasicBox
}

// DecodeBoxMOOV解析moov box
func DecodeBoxMOOV(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器 box 解析子 box
	children, err := DecodeChildren(readSeeker, boxSize-headerSize)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(MOOV)
	box.size = boxSize
	box._type = _type
	box.BasicBox.children = children
	// 返回
	return box, nil
}
