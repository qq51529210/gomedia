package rtsp

import (
	"errors"
	"strings"
)

var (
	errHeaderFormat = errors.New("error header format")
)

// parseHeader解析key:value的行，返回key,value
func parseHeader(line string) (string, string, error) {
	i := strings.IndexByte(line, ':')
	if i < 0 {
		return "", "", errHeaderFormat
	}
	return strings.TrimSpace(line[:i]), strings.TrimSpace(line[i+1:]), nil
}
