package m3u8

// Segment 表示 m3u8 列表中的一个片段
type Segment struct {
	// 时长
	Duration float64
	//
	Title string
	// 地址
	URL string
}
