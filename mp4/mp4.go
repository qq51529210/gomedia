package mp4

import "io"

// MP4 用于表示一个 mp4 文件
// 其实就是第一层的 box
type MP4 struct {
	Box []Box
}

// Decode 解析所有的 box
func (m *MP4) Decode(readSeeker io.ReadSeeker) error {
	// 循环解析即可
	for {
		box, err := DecodeBox(readSeeker)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		m.Box = append(m.Box, box)
	}
}

// GetBox 返回指定 _type 的子 box
func (m *MP4) GetBox(_type Type) Box {
	for _, box := range m.Box {
		if box.Type() == _type {
			return box
		}
	}
	return nil
}
