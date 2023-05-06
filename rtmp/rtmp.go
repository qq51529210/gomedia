package rtmp

// package main
//  import (
// 	"bufio"
// 	"bytes"
// 	"encoding/binary"
// 	"fmt"
// 	"net"
// )
//  type RtmpHeader struct {
// 	Format       int
// 	ChunkStream  int
// 	Timestamp    uint32
// 	MsgLength    uint32
// 	MsgTypeId    uint8
// 	MsgStreamId  uint32
// 	ExtendedTime uint32
// }
//  type RtmpSession struct {
// 	conn         net.Conn
// 	readBuffer   *bufio.Reader
// 	chunkSize    int
// 	prevHeaders  map[int]*RtmpHeader
// 	prevMessages map[int][]byte
// }
//  func NewRtmpSession(conn net.Conn) *RtmpSession {
// 	return &RtmpSession{
// 		conn:         conn,
// 		readBuffer:   bufio.NewReader(conn),
// 		chunkSize:    128,
// 		prevHeaders:  make(map[int]*RtmpHeader),
// 		prevMessages: make(map[int][]byte),
// 	}
// }
//  func (s *RtmpSession) ReadPacket() ([]byte, error) {
// 	var packet []byte
//  	header := &RtmpHeader{}
// 	chunkStream := 0
//  	for {
// 		dataBuffer := make([]byte, s.chunkSize)
//  		if header.Format == 0 {
// 			binary.Read(s.readBuffer, binary.BigEndian, &header.ChunkStream)
// 			header.ChunkStream = header.ChunkStream & 0x3f
//  			binary.Read(s.readBuffer, binary.BigEndian, &header.Timestamp)
//  			binary.Read(s.readBuffer, binary.BigEndian, &header.MsgLength)
// 			binary.Read(s.readBuffer, binary.BigEndian, &header.MsgTypeId)
//  			binary.Read(s.readBuffer, binary.BigEndian, &header.MsgStreamId)
// 		} else {
// 			prevHeader := s.prevHeaders[chunkStream]
//  			if header.Format == 1 {
// 				binary.Read(s.readBuffer, binary.BigEndian, &header.Timestamp)
// 				header.Timestamp += prevHeader.Timestamp
//  				binary.Read(s.readBuffer, binary.BigEndian, &header.MsgLength)
// 				header.MsgLength = prevHeader.MsgLength
//  				binary.Read(s.readBuffer, binary.BigEndian, &header.MsgTypeId)
//  				header.MsgStreamId = prevHeader.MsgStreamId
// 			} else if header.Format == 2 {
// 				binary.Read(s.readBuffer, binary.BigEndian, &header.Timestamp)
// 				header.Timestamp += prevHeader.Timestamp
//  				header.MsgLength = prevHeader.MsgLength
//  				header.MsgTypeId = prevHeader.MsgTypeId
//  				header.MsgStreamId = prevHeader.MsgStreamId
// 			} else {
// 				header.Timestamp = prevHeader.Timestamp
//  				header.MsgLength = prevHeader.MsgLength
//  				header.MsgTypeId = prevHeader.MsgTypeId
//  				binary.Read(s.readBuffer, binary.BigEndian, &header.MsgStreamId)
// 			}
// 		}
//  		if header.Timestamp == 0xffffff {
// 			binary.Read(s.readBuffer, binary.BigEndian, &header.ExtendedTime)
// 		}
//  		dataBuffer = dataBuffer[:header.MsgLength]
// 		s.readBuffer.Read(dataBuffer)
//  		if header.Format == 0 {
// 			s.prevHeaders[chunkStream] = header
// 			packet = dataBuffer
// 		} else {
// 			prevPacket := s.prevMessages[chunkStream]
// 			packet = make([]byte, len(prevPacket)+int(header.MsgLength))
// 			copy(packet, prevPacket)
// 			copy(packet[len(prevPacket):], dataBuffer)
// 			s.prevMessages[chunkStream] = packet
// 		}
//  		if header.MsgLength != uint32(len(packet)) {
// 			continue
// 		}
//  		break
// 	}
//  	return packet, nil
// }
//  func main() {
// 	conn, err := net.Dial("tcp", "localhost:1935")
// 	if err != nil {
// 		fmt.Println("连接失败:", err)
// 		return
// 	}
// 	defer conn.Close()
//  	rtmpSession := NewRtmpSession(conn)
//  	for {
// 		data, err := rtmpSession.ReadPacket()
// 		if err != nil {
// 			fmt.Println("读取数据失败:", err)
// 			return
// 		}
//  		// 处理接收到的数据
// 		fmt.Println("接收到的RTMP数据:", bytes.NewBuffer(data).String())
// 	}
// }
