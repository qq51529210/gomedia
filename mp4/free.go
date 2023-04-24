package mp4

import (
	"io"
)

const (
	TypeFREE = 1718773093
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeFREE, DecodeBoxFREE)
}

// FREE表示free box
// 内容是无关紧要的，可以被忽略。
// 该box被删除后，不会对播放产生任何影响
type FREE struct {
	BasicBox
}

// DecodeBoxFREE解析free box
func DecodeBoxFREE(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize > 0 {
		// 忽略
		_, err := readSeeker.Seek(contentSize, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}
	// 解析
	box := new(FREE)
	box.size = boxSize
	box._type = _type
	// 返回
	return box, nil
}
