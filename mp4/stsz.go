package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	TypeSTSZ = 1937011578
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeSTSZ, DecodeBoxSTSZ)
}

// STSZ表示stsz box
// 定义了每个sample的大小，包含了媒体中全部sample的数目和一张给出每个sample大小的表
type STSZ struct {
	fullBox
	// 如果是0就使用Size数组
	SampleSize uint32
	// sample的数量
	SampleCount uint32
	// 每一个sample的size
	SizeList []uint32
}

// DecodeBoxSTSZ解析stsz box
func DecodeBoxSTSZ(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < 12 {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	// 解析
	box := new(STSZ)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.Uint24(buf[1:])
	// 4
	box.SampleSize = binary.BigEndian.Uint32(buf[4:])
	// 4
	box.SampleCount = binary.BigEndian.Uint32(buf[8:])
	//
	n := 12
	box.SizeList = make([]uint32, box.SampleCount)
	if box.SampleSize == 0 {
		if contentSize < int64(box.SampleCount)*4+int64(n) {
			return nil, errBoxSize
		}
		for i := 0; i < len(box.SizeList); i++ {
			box.SizeList[i] = binary.BigEndian.Uint32(buf[n:])
			n += 4
		}
	} else {
		for i := 0; i < len(box.SizeList); i++ {
			box.SizeList[i] = box.SampleSize
		}
	}
	// 返回
	return box, nil
}
