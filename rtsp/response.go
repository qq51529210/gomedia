package rtsp

import (
	"errors"
	"fmt"
	"gomedia/util"
	"io"
	"strings"
)

var (
	errResponseFormat = errors.New("error response format")
)

// Response表示请求
type Response struct {
	Version string
	Status  string
	Phrase  string
	Header  map[string]string
}

// Decode解码，但不包括body
func (r *Response) Decode(reader io.Reader) error {
	crlfReader := util.NewCRLFReader(reader)
	// 第一行
	line, err := crlfReader.ReadLine()
	if err != nil {
		return err
	}
	// version
	i := strings.IndexByte(line, ' ')
	if i < 0 {
		return errResponseFormat
	}
	r.Version = line[:i]
	line = line[i+1:]
	// stauts
	i = strings.IndexByte(line, ' ')
	if i < 0 {
		return errResponseFormat
	}
	r.Status = line[:i]
	// phrase
	r.Phrase = line[i+1:]
	// header
	for {
		line, err = crlfReader.ReadLine()
		if err != nil {
			return errResponseFormat
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
func (r *Response) Encode(writer io.Writer) error {
	// method url version
	_, err := fmt.Fprintf(writer, "%s %s %s\r\n", r.Version, r.Status, r.Phrase)
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
