package mp4

import "io"

const (
	TypeMDAT = 1835295092
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMDAT, DecodeBoxMDAT)
}

// MDAT表示mdat box
// 可以有多个，也可以没有（当媒体数据全部为外部文件引用时）。
// 数据直接跟在box type字段后面
type MDAT struct {
	BasicBox
}

// DecodeBoxMDAT解析mdat box
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
