package mp4

import "io"

const (
	// TypeMDAT 表示 mdat 类型
	TypeMDAT Type = 1835295092
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMDAT, DecodeBoxMDAT)
}

// MDAT 表示 mdat box
// 是真正的数据
type MDAT struct {
	BasicBox
}

// DecodeBoxMDAT 解析 mdat box
func DecodeBoxMDAT(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 略过
	contentSize := boxSize - headerSize
	if contentSize > 0 {
		_, err := readSeeker.Seek(contentSize, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}
	// 创建
	box := new(MDAT)
	box.size = boxSize
	box._type = _type
	// 返回
	return box, nil
}
