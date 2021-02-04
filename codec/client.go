package codec

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func StartClient(addr chan string){
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close()}()

	time.Sleep(time.Second)

	_ = json.NewEncoder(conn).Encode(DefaultOption)
	cc := NewGobCodec(conn)

	for i := 0; i < 5; i++{
		h := & Header{
			ServiceMethod: "Foo.Sum",
			Seq: uint64(i),
		}

		_ = cc.Write(h, fmt.Sprintf("brpc req:%d",h.Seq))
		_ = cc.ReadHeader(h)

		var reply string

		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}

