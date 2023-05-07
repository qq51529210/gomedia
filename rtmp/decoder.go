package rtmp

// import (
// 	"encoding/binary"
// 	"fmt"
// 	"net"
// )

// func main() {
// 	conn, err := net.Dial("tcp", "server.example.com:1935")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close()

// 	// 读取 RTMP 消息头
// 	head := &RtmpMsgHeader{}
// 	err = binary.Read(conn, binary.BigEndian, head)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 根据消息类型处理不同的消息内容
// 	switch head.MsgType {
// 	case 1:
// 		// 处理 SetChunkSize 消息
// 		var chunkSize uint32
// 		err = binary.Read(conn, binary.BigEndian, &chunkSize)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("SetChunkSize: %d\n", chunkSize)
// 	case 4:
// 		// 处理 Acknowledgement 消息
// 		var ack uint32
// 		err = binary.Read(conn, binary.BigEndian, &ack)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("Acknowledgement: %d\n", ack)
// 	default:
// 		// 其他消息
// 		buf := make([]byte, head.MsgLength)
// 		err = binary.Read(conn, binary.BigEndian, buf)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("Unknown message type %d with length %d\n", head.MsgType, head.MsgLength)
// 	}
// }
