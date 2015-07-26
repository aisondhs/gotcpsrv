package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"github.com/Unknwon/goconfig"
	"github.com/aisondhs/gotcpsrv/service"
	"github.com/aisondhs/gotcpsrv/lib/funcmap"
	"github.com/aisondhs/gotcpsrv/lib/gametcp"
	"github.com/aisondhs/alog"
	"github.com/aisondhs/gotcpsrv/protocol"
	"github.com/aisondhs/gotcpsrv/protos"
)

type Callback struct{}

var funcs funcmap.Funcs

var HTTP_PORT string

var logdir string = "./logs"

func init() {
	// bind func map
	funcs = funcmap.NewFuncs(100)
	funcs.Bind("CSGetuserReq", service.CSGetuserReq)
	alog.Init(logdir,alog.ROTATE_BY_DAY,false)
}

func (this *Callback) OnConnect(c *gametcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	alog.Info("OnConnect:"+addr.String())
	return true
}

func (this *Callback) OnMessage(c *gametcp.Conn, p protocol.Packet) bool {
	packet := &p

	reqBytes := packet.GetBody()
	msgId := packet.GetMsgId()
	methodName := protos.GetFuncName(msgId)
	reflectData,err := funcs.Call(methodName, reqBytes)
	checkError(err)
	i := reflectData[0].Interface()
	rspBytes := i.([]byte)
	rspPacket := protocol.NewPacket(rspBytes, msgId+1, false)
	err = c.AsyncWritePacket(rspPacket, time.Second)
	checkError(err)
	return true
}

func (this *Callback) OnClose(c *gametcp.Conn) {
	alog.Info("OnClose:"+c.GetExtraData().(net.Addr).String())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	c, err := goconfig.LoadConfigFile("conf/conf.ini")
	if err != nil {
		log.Fatal(err)
	}

	HTTP_PORT, err = c.GetValue("Server", "port")
	checkError(err)

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+HTTP_PORT)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	sendChan, err := c.Int("Server", "sendChan")
	checkError(err)
	receiveChan, err := c.Int("Server", "receiveChan")
	checkError(err)

	// creates a server
	config := &gametcp.Config{
		PacketSendChanLimit:    uint32(sendChan),
		PacketReceiveChanLimit: uint32(receiveChan),
	}
	srv := gametcp.NewServer(config, &Callback{})

	// starts service
	go srv.Start(listener, time.Second*5)
	alog.Info("listening:"+listener.Addr().String())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	sig := <-chSig
	alog.Error("listening:"+sig.String())

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		alog.Error(err.Error())
	}
}
