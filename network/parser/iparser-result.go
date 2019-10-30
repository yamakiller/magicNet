package parser

//IResult Parser Decoder result
type IResult interface {
	GetVersion() int32
	GetCommand() interface{}
	GetSerial() uint32
	GetWrap() interface{}
	GetExtern() interface{}
}
