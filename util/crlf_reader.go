package util

import (
	"io"
)

// CRLFReader用于读取crlf的每一行数据
type CRLFReader struct {
	// 数据源
	r io.Reader
}

// ReadLine读取一行数据
func (r *CRLFReader) ReadLine() (string, error) {
	return "", nil
}

// NewCRLFReader返回新的CRLFReader
func NewCRLFReader(reader io.Reader) *CRLFReader {
	return &CRLFReader{r: reader}
}
