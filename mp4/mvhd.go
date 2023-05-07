package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	// TypeMVHD 表示 mvhd 类型
	TypeMVHD Type = 1836476516
)

const (
	mvhdBoxMinContentSize  = 100
	mvhdBoxMinContentSize2 = mvhdBoxMinContentSize + 12
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeMVHD, DecodeBoxMVHD)
}

// MVHD 表示 mvhd box
// 用来存放文件的总体信息
type MVHD struct {
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
	Rate uint32
	// 音量
	Volume uint16
	// 视频变换矩阵
	Matrix [9]uint32
	// ...
	PreDefined [6]uint32
	// 下一个 track id
	NextTrackID uint32
}

// DecodeBoxMVHD 解析 mvhd box
func DecodeBoxMVHD(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < mvhdBoxMinContentSize {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	box := new(MVHD)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.BigUint24(buf[1:])
	n := 0
	if box.Version == 1 {
		if contentSize < mvhdBoxMinContentSize2 {
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
	// 4
	box.Rate = binary.BigEndian.Uint32(buf[n:])
	n += 4
	// 2
	box.Volume = binary.BigEndian.Uint16(buf[n:])
	// 2+10 reserved
	n += 12
	// 36
	for i := 0; i < len(box.Matrix); i++ {
		box.Matrix[i] = binary.BigEndian.Uint32(buf[n:])
		n += 4
	}
	// 24
	for i := 0; i < len(box.PreDefined); i++ {
		box.PreDefined[i] = binary.BigEndian.Uint32(buf[n:])
		n += 4
	}
	// 4
	box.NextTrackID = binary.BigEndian.Uint32(buf[n:])
	n += 4
	// 返回
	return box, nil
}
