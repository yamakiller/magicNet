package parser

//Abandoned

//IParam Parser Encoder param
type IParam interface {
	SetVersion(v int32)
	SetCommand(v interface{})
	SetSerial(v uint32)
	SetWrap(v interface{})
	SetWrapLength(v int)
	SetExtern(v interface{})
}
