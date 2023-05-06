package mp4

// Sample 表示某条轨道的一个数据
type Sample struct {
	// 文件的偏移
	Offset int64
	// 大小
	Size int64
}
