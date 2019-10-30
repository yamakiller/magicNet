package parser

import "bytes"

// IParser Parser interface
type IParser interface {
	Decoder(keyPair interface{}, data *bytes.Buffer) (IResult, error)
	Encoder(keyPair interface{}, param IParam) []byte
}
