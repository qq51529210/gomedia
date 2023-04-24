package mp4

import (
	"io"
)

const (
	TypeTRAK = 1953653099
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeTRAK, DecodeBoxTRAK)
}

// TRAK表示trak box
// 子box包含了该track的媒体数据引用和描述（hint track除外）。
// 一个MP4文件中的媒体可以包含多个track，且至少有一个track，
// 这些track之间彼此独立，有自己的时间和空间信息。
// trak box必须包含一个tkhd box和一个mdia box
type TRAK struct {
	BasicBox
}

// DecodeBoxTRAK解析trak box
func DecodeBoxTRAK(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 容器box解析子box
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
