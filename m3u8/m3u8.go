package m3u8

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// 字段名称
const (
	headerSign           = "#EXTM3U"
	headerVersion        = "#EXT-X-VERSION"
	headerTargetDuration = "#EXT-X-TARGETDURATION"
	headerMediaSequence  = "#EXT-X-MEDIA-SEQUENCE"
	headerEndList        = "#EXT-X-ENDLIST"
	headerINF            = "#EXTINF"
)

var (
	errInfo = errors.New("decode info error: miss comma")
)

// M3U8 用于编码 m3u8 文件
type M3U8 struct {
	// 版本，一般是 3
	Version string
	// 媒体片段文件的最长时长
	TargetDuration float64
	// 第一个媒体片段文件的序号
	MediaSequence int64
	// 片段
	Segment []*Segment
	// 媒体播放列表已经结束
	EndList bool
}

// Encode 编码
func (m *M3U8) Encode(writer io.Writer) error {
	// #EXTM3U
	if _, err := fmt.Fprintln(writer, headerSign); err != nil {
		return fmt.Errorf("encode header sign error: %w", err)
	}
	// #EXT-X-VERSION
	if _, err := fmt.Fprintf(writer, "%s:%s\n", headerVersion, m.Version); err != nil {
		return fmt.Errorf("encode version error: %w", err)
	}
	// #EXT-X-TARGETDURATION
	if _, err := fmt.Fprintf(writer, "%s:%.3f", headerTargetDuration, m.TargetDuration); err != nil {
		return fmt.Errorf("encode target duration error: %w", err)
	}
	// #EXT-X-MEDIA-SEQUENCE
	if _, err := fmt.Fprintf(writer, "%s:%d", headerMediaSequence, m.MediaSequence); err != nil {
		return fmt.Errorf("encode media sequence error: %w", err)
	}
	// #EXTINF
	for _, segment := range m.Segment {
		if _, err := fmt.Fprintf(writer, "%s:%.3f,%s\n%s\n", headerINF, segment.Duration, segment.Title, segment.URL); err != nil {
			return fmt.Errorf("encode inf error: %w", err)
		}
	}
	// #EXT-X-ENDLIST
	if m.EndList {
		if _, err := fmt.Fprintln(writer, headerEndList); err != nil {
			return fmt.Errorf("encode end list error: %w", err)
		}
	}
	// 返回
	return nil
}

// Decode 解码
func (m *M3U8) Decode(reader io.Reader) error {
	// 扫描器
	scanner := bufio.NewScanner(reader)
	// 扫描函数
	scanner.Split(bufio.ScanLines)
	// 循环
	for scanner.Scan() {
		// 拿到一行
		line := scanner.Text()
		// 空行
		if line == "" {
			continue
		}
		// #EXTINF
		// 這個字段出现比較多，放前面解析，减少其他的判断
		p := strings.TrimPrefix(line, headerINF)
		if p != line {
			if err := m.decodeINF(scanner, strings.TrimSpace(p[1:])); err != nil {
				return err
			}
			continue
		}
		// #EXTM3U
		p = strings.TrimPrefix(line, headerSign)
		if p != line {
			continue
		}
		// #EXT-X-VERSION
		p = strings.TrimPrefix(line, headerVersion)
		if p != line {
			m.Version = strings.TrimSpace(p[1:])
			continue
		}
		// #EXT-X-MEDIA-SEQUENCE
		p = strings.TrimPrefix(line, headerMediaSequence)
		if p != line {
			n, err := strconv.ParseInt(strings.TrimSpace(p[1:]), 10, 64)
			if err != nil {
				return fmt.Errorf("decode media sequence error: %w", err)
			}
			m.MediaSequence = n
			continue
		}
		// #EXT-X-TARGETDURATION
		p = strings.TrimPrefix(line, headerTargetDuration)
		if p != line {
			n, err := strconv.ParseFloat(strings.TrimSpace(p[1:]), 64)
			if err != nil {
				return fmt.Errorf("decode target duration error: %w", err)
			}
			m.TargetDuration = n
			continue
		}
		// #EXT-X-ENDLIST
		p = strings.TrimPrefix(line, headerEndList)
		if p != line {
			m.EndList = true
			continue
		}
	}
	// 返回
	return scanner.Err()
}

// decodeINF 拆分代码
func (m *M3U8) decodeINF(scanner *bufio.Scanner, unparseLine string) error {
	i := strings.IndexByte(unparseLine, ',')
	if i < 0 {
		return errInfo
	}
	// duration
	duration, err := strconv.ParseFloat(strings.TrimSpace(unparseLine[:i]), 64)
	if err != nil {
		return fmt.Errorf("decode info duration error: %w", err)
	}
	// title
	title := unparseLine[i+1:]
	// url
	if !scanner.Scan() {
		return fmt.Errorf("decode info url error: %w", err)
	}
	segment := new(Segment)
	segment.Duration = duration
	segment.Title = title
	segment.URL = scanner.Text()
	m.Segment = append(m.Segment, segment)
	// 返回
	return nil
}
