package mp4

import "io"

// MP4用于表示一个box结构树
type MP4 struct {
	Box []Box
}

// Decode解析所有的box并组成一棵树，如果有特别的box要自己解析，先注册自己的解析函数
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

func (m *MP4) GetBox(_type Type) Box {
	for _, box := range m.Box {
		if box.Type() == _type {
			return box
		}
	}
	return nil
}
