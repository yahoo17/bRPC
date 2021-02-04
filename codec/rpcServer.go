package codec

import (

	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int	//magic number标记这是一个brpc 请求
	CodecType Type // 客户端可能会选择不同的Codec来解析
}
var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:  GobType,
}

type Server struct {

}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server)Accept(lis net.Listener) {
	for{
		conn,err := lis.Accept()
		if err != nil{
			log.Println("rpc server: accept error",err)
			return
		}
		go server.ServeConn(conn)

	}

}
func (server *Server)ServeConn(conn io.ReadWriteCloser){
	defer func() { _ = conn.Close()}()

	var opt Option

	if err := json.NewDecoder(conn).Decode(&opt); err != nil{
		log.Println("rpc server:option error",err)
		return

	}

	if opt.MagicNumber != MagicNumber {
		log.Println("rpc server: invalid magic number",opt.MagicNumber)
		return
	}

	f := NewCodecFuncMap[opt.CodecType]

	if f == nil{
		log.Println("rpc server: invalid codec type",opt.CodecType)
		return
	}
	server.serveCodec( f(conn) )
}

var invalidRequest = struct {}{}

type request struct {
	h *Header
	argv reflect.Value
	replyv reflect.Value
}
func (server *Server)readRequestHeader(cc Codec)(*Header, error){
	var h Header
	if err := cc.ReadHeader(&h); err != nil{
		if err != io.EOF && err != io.ErrUnexpectedEOF{
			log.Println("rpc server: read header error",err)
		}
		return nil, err
	}
	return &h, nil

}

func (server *Server) readRequest(cc Codec)(*request, error)  {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}
	// 我们还不知道request的类型
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil{
		log.Println("rpc server: read argv err",err)
	}

	return  req,nil

}

func (server *Server) sendResponse(cc Codec, h *Header,
	body interface{}, sending * sync.Mutex) {

	sending.Lock()
	defer sending.Unlock()

	if err := cc.Write(h, body); err != nil{
		log.Println("rpc serer: write response error",err)
	}

}
func (server *Server) handleRequest(cc Codec, req *request,
	sending *sync.Mutex, wg * sync.WaitGroup){
	//
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("bRpc resp: %s",req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}
func (server *Server)serveCodec(cc Codec)  {
	//读取请求
	//处理请求
	//回复请求
	sending := new (sync.Mutex)
	wg := new(sync.WaitGroup)
	for{
		req, err := server.readRequest(cc)
		if err != nil{
			if req == nil{
				break;
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()

}

func Accept(lis net.Listener){
	DefaultServer.Accept(lis)
}