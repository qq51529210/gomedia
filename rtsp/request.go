package rtsp

import (
	"errors"
	"fmt"
	"gomedia/util"
	"io"
	"net/url"
	"strings"
)

var (
	errRequestFormat = errors.New("error request format")
)

// Request表示请求
type Request struct {
	Method  string
	URL     *url.URL
	Version string
	Header  map[string]string
}

// Decode解码，但不包括body
func (r *Request) Decode(reader io.Reader) error {
	crlfReader := util.NewCRLFReader(reader)
	// 第一行
	line, err := crlfReader.ReadLine()
	if err != nil {
		return err
	}
	// method
	i := strings.IndexByte(line, ' ')
	if i < 0 {
		return errRequestFormat
	}
	r.Method = line[:i]
	line = line[i+1:]
	// version
	i = strings.LastIndexByte(line, ' ')
	if i < 0 {
		return errRequestFormat
	}
	r.Version = line[i+1:]
	line = line[:i]
	r.URL, err = url.Parse(line)
	if err != nil {
		return errRequestFormat
	}
	// header
	for {
		line, err = crlfReader.ReadLine()
		if err != nil {
			return errRequestFormat
		}
		// 空行
		if line == "" {
			break
		}
		k, v, err := parseHeader(line)
		if err != nil {
			return err
		}
		r.Header[k] = v
	}
	// 返回
	return nil
}

// Encode编码，但不包括body
func (r *Request) Encode(writer io.Writer) error {
	// method url version
	_, err := fmt.Fprintf(writer, "%s %s %s\r\n", r.Method, r.URL.String(), r.Version)
	if err != nil {
		return err
	}
	// header
	for k, v := range r.Header {
		_, err = fmt.Fprintf(writer, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	// 空行
	_, err = writer.Write(crlf)
	// 返回
	return err
}
