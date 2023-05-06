package mp4

import (
	"encoding/binary"
	"io"
)

const (
	// TypeFTYP 表示 ftyp 类型
	TypeFTYP Type = 1718909296
)

const (
	ftypBoxMinContentSize = 12
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeFTYP, DecodeBoxFTYP)
}

// FTYP 表示 ftyp box
// 只有一个并且只能被包含在文件顶层,
// 同时应该出现在文件的最开始的位置
type FTYP struct {
	BasicBox
	MajorBrand       uint32
	MinorVersion     uint32
	CompatibleBrands []byte
}

// DecodeBoxFTYP 解析 ftyp box
func DecodeBoxFTYP(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < ftypBoxMinContentSize {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	// 解析
	box := new(FTYP)
	box.size = boxSize
	box._type = _type
	// 4
	box.MajorBrand = binary.BigEndian.Uint32(buf)
	// 4
	box.MinorVersion = binary.BigEndian.Uint32(buf[4:])
	// 4
	box.CompatibleBrands = buf[8:]
	// 返回
	return box, nil
}
