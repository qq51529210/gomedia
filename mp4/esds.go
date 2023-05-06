package mp4

//  import (
//     "bytes"
//     "encoding/binary"
//     "fmt"
//     "os"
// )
//  type Descriptor struct {
//     Tag byte
//     Data []byte
// }
//  type EsdsBox struct {
//     Version byte
//     Flags [3]byte
//     Descriptors []Descriptor
// }
//  func readBox(r *bytes.Reader) *Box {
//     size := readUint32(r)
//     boxType := readBoxType(r)
//     if size == 1 {
//         size = readUint64(r)
//     } else if size == 0 {
//         // the box extends to the end of the file
//         size = uint32(r.Len() + 8)
//     }
//     data := make([]byte, size-8)
//     _, _ = r.Read(data)
//     return &Box{
//         Size: size,
//         Type: boxType,
//         Data: data,
//     }
// }
//  func readBoxType(r *bytes.Reader) string {
//     b := make([]byte, 4)
//     _, _ = r.Read(b)
//     return string(b)
// }
//  func readUint32(r *bytes.Reader) uint32 {
//     var u uint32
//     _ = binary.Read(r, binary.BigEndian, &u)
//     return u
// }
//  func readEsdsBox(box *Box) *EsdsBox {
//     r := bytes.NewReader(box.Data)
//     esds := &EsdsBox{}
//     esds.Version = readByte(r)
//     _ = readUint24(r) // flags
//     for r.Len() > 0 {
//         tag, length, err := readDescriptorTagAndLength(r)
//         if err != nil {
//             fmt.Println(err)
//             break
//         }
//         data := make([]byte, length)
//         _, _ = r.Read(data)
//         esds.Descriptors = append(esds.Descriptors, Descriptor{
//             Tag:  tag,
//             Data: data,
//         })
//     }
//     return esds
// }
//  func readDescriptorTagAndLength(r *bytes.Reader) (byte, int, error) {
//     tag := readByte(r)
//     var length int
//     b := readByte(r)
//     if b&0x80 != 0 {
//         // long form
//         lengthBytes := make([]byte, b&0x7f)
//         _, _ = r.Read(lengthBytes)
//         for _, b := range lengthBytes {
//             length = length<<7 | int(b&0x7f)
//         }
//     } else {
//         // short form
//         length = int(b & 0x7f)
//     }
//     if length > r.Len() {
//         return 0, 0, fmt.Errorf("invalid descriptor length")
//     }
//     return tag, length, nil
// }
//  func readByte(r *bytes.Reader) byte {
//     b, _ := r.ReadByte()
//     return b
// }
//  func main() {
//     // Open the MP4 file
//     f, err := os.Open("audio.mp4")
//     if err != nil {
//         fmt.Println(err)
//         return
//     }
//     defer f.Close()
//      // Find the esds box
//     var esdsBox *Box
//     for {
//         box := readBox(bytes.NewReader(nil))
//         if box == nil {
//             break
//         }
//         if box.Type == "esds" {
//             esdsBox = box
//             break
//         }
//     }
//      // Parse the esds box
//     esds := readEsdsBox(esdsBox)
//     fmt.Printf("%+v\n", esds)
// }
