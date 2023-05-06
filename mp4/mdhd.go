package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	// TypeMDHD 表示 mdhd 类型
	TypeMDHD Type = 1835296868
)

const (
	mdhdBoxMinContentSize = 24
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMDHD, DecodeBoxMDHD)
}

// MDHD 表示 mdhd box
// 包含这一条 track 的信息
type MDHD struct {
	fullBox
	// 创建时间
	CreateTime uint64
	// 修改时间
	ModTime uint64
	// 缩放因子
	TimeScale uint32
	// 时长
	Duration uint64
	// 播放速率
	Language uint16
	// ...
	PreDefined uint16
}

// DecodeBoxMDHD 解析 mdhd box
func DecodeBoxMDHD(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < mdhdBoxMinContentSize {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	box := new(MDHD)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.Uint24(buf[1:])
	n := 0
	if box.Version == 1 {
		// 判断
		if contentSize < 36 {
			return nil, errBoxSize
		}
		// 8
		box.CreateTime = binary.BigEndian.Uint64(buf[4:])
		// 8
		box.ModTime = binary.BigEndian.Uint64(buf[12:])
		// 4
		box.TimeScale = binary.BigEndian.Uint32(buf[20:])
		// 8
		box.Duration = binary.BigEndian.Uint64(buf[24:])
		n = 32
	} else {
		// 4
		box.CreateTime = uint64(binary.BigEndian.Uint32(buf[4:]))
		// 4
		box.ModTime = uint64(binary.BigEndian.Uint32(buf[8:]))
		// 4
		box.TimeScale = binary.BigEndian.Uint32(buf[12:])
		// 4
		box.Duration = uint64(binary.BigEndian.Uint32(buf[16:]))
		n = 20
	}
	// 2
	box.Language = binary.BigEndian.Uint16(buf[n:])
	n += 2
	// 最高位是 0
	box.Language &= 0x7FFF
	// 2
	box.PreDefined = binary.BigEndian.Uint16(buf[n:])
	n += 2
	// 返回
	return box, nil
}
