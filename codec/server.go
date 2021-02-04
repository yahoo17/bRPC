package codec

import (
	"log"
	"net"
)

func StartServer(addr chan string){

	l,err := net.Listen("tcp",":0")
	if err != nil{
		log.Fatal("network error",err)
	}
	log.Println("start rpc server on",l.Addr())
	addr <- l.Addr().String()
	Accept(l)
}

