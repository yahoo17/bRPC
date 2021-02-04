package codec

import "io"

//定义RPC的抽象接口
//必须实现Codec的方法, 分别是读头 读body, 写
type Header struct {
	ServiceMethod string
	Seq	uint64
	Error string
	//ErrorCode unint64
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody( interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(closer io.ReadWriteCloser ) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" //not implement
)

var NewCodecFuncMap map[Type]NewCodecFunc

//init函数的主要特点：
//
//init函数先于main函数自动执行，不能被其他函数调用；
//init函数没有输入参数、返回值；
//每个包可以有多个init函数；
func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec

}


