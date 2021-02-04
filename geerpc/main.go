package main

import "bRPC/codec"

func main()  {
	addr := make(chan string)

	go codec.StartServer(addr)
	codec.StartClient(addr)

}
