
package main
import proto "github.com/golang/protobuf/proto"
import (
	"encoding/binary"
	"fmt"
	"github.com/aisondhs/gotcpsrv/protocol"
	"github.com/aisondhs/gotcpsrv/protos"
	"log"
	"net"
)
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	var msgId uint16
	msgId = 1
	reqObj := new(protos.CSGetuserReq)
	reqObj.Uid =  proto.Int32(101)

	reqBytes, err := proto.Marshal(reqObj)
	checkError(err)

	var reqBuff []byte = make([]byte, 4+len(reqBytes))
	binary.BigEndian.PutUint16(reqBuff[0:2], uint16(len(reqBuff)))
	binary.BigEndian.PutUint16(reqBuff[2:4], msgId)
	copy(reqBuff[4:], reqBytes)
	// write
	conn.Write(reqBuff)

	// read
	p, err := protocol.ReadPacket(conn)
	checkError(err)
	rspObjt := new(protos.CSGetuserRsp)
	body := p.GetBody()
	proto.Unmarshal(body, rspObjt)
	fmt.Println(rspObjt)
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
