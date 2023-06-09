package mp4

import (
	"encoding/binary"
	"gomedia/util"
	"io"
)

const (
	// TypeTKHD 表示 tkhd 类型
	TypeTKHD Type = 1953196132
)

const (
	tkhdBoxMinContentSize  = 84
	tkhdBoxMinContentSize2 = tkhdBoxMinContentSize + 12
)

func init() {
	// 注册解析器
	AddDecodeFunc(TypeTKHD, DecodeBoxTKHD)
}

// TKHD 表示 tkhd box
// 包含了该 track 信息
type TKHD struct {
	fullBox
	// 创建时间
	CreateTime uint64
	// 修改时间
	ModTime uint64
	// 当前的track id
	TrackID uint32
	// 时长
	Duration uint64
	// 分层，默认为0，值小的在上层
	Layer uint16
	// 指定 track 分组信息,
	// 默认为 0 表示该 track 未与其他 track 有群组关系
	AlternateGroup uint16
	// 音量, 如果为音频 1 表示最大音量, 否则为 0
	Volume uint16
	// 视频变换矩阵
	Matrix [9]uint32
	// 宽
	Width uint32
	// 高
	Height uint32
}

// DecodeBoxTKHD 解析 tkhd box
func DecodeBoxTKHD(readSeeker io.ReadSeeker, headerSize, boxSize int64, _type Type) (Box, error) {
	// 判断
	contentSize := boxSize - headerSize
	if contentSize < tkhdBoxMinContentSize {
		return nil, errBoxSize
	}
	// 读取
	buf := make([]byte, contentSize)
	_, err := io.ReadFull(readSeeker, buf)
	if err != nil {
		return nil, err
	}
	box := new(TKHD)
	box.size = boxSize
	box._type = _type
	// 1
	box.Version = buf[0]
	// 3
	box.Flags = util.BigUint24(buf[1:])
	n := 0
	if box.Version == 1 {
		if contentSize < tkhdBoxMinContentSize2 {
			return nil, errBoxSize
		}
		// 8
		box.CreateTime = binary.BigEndian.Uint64(buf[4:])
		// 8
		box.ModTime = binary.BigEndian.Uint64(buf[12:])
		// 4
		box.TrackID = binary.BigEndian.Uint32(buf[20:])
		// 4 reserved
		// 8
		box.Duration = binary.BigEndian.Uint64(buf[28:])
		n = 36
	} else {
		// 4
		box.CreateTime = uint64(binary.BigEndian.Uint32(buf[4:]))
		// 4
		box.ModTime = uint64(binary.BigEndian.Uint32(buf[8:]))
		// 4
		box.TrackID = binary.BigEndian.Uint32(buf[12:])
		// 4 reserved
		// 4
		box.Duration = uint64(binary.BigEndian.Uint32(buf[20:]))
		n = 24
	}
	// 4*2 reserved
	n += 8
	// 2
	box.Layer = binary.BigEndian.Uint16(buf[n:])
	n += 2
	// 2
	box.AlternateGroup = binary.BigEndian.Uint16(buf[n:])
	n += 2
	// 2
	box.Volume = binary.BigEndian.Uint16(buf[n:])
	// 2+2 reserved
	n += 4
	// 36
	for i := 0; i < len(box.Matrix); i++ {
		box.Matrix[i] = binary.BigEndian.Uint32(buf[n:])
		n += 4
	}
	// 4
	box.Width = binary.BigEndian.Uint32(buf[n:])
	n += 4
	// 4
	box.Height = binary.BigEndian.Uint32(buf[n:])
	n += 4
	// 返回
	return box, nil
}
