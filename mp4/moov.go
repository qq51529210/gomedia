package mp4

import "io"

const (
	TypeMOOV = 1836019574
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMOOV, DecodeBoxMOOV)
}

// MOOV表示moov box
// 用来存放媒体的metadata信息，其内容信息由子box诠释。
// 该box有且只有一个并且包含在文件层，
// 一般情况下moov box会紧随ftyp box出现，但也有放在文件末尾的。
type MOOV struct {
	BasicBox
}

// DecodeBoxMOOVP解析moov box
func DecodeBoxMOOV(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器box解析子box
	children, err := DecodeChildren(readSeeker, boxSize-headerSize)
	if err != nil {
		return nil, err
	}
	// 创建
	box := new(MOOV)
	box.BasicBox.size = boxSize
	box.BasicBox._type = _type
	box.BasicBox.children = children
	// 返回
	return box, nil
}
