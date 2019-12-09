package parser

import "github.com/yamakiller/magicNet/handler/net"

// IParser Parser interface
type IParser interface {
	Decoder(c net.INetClient) (IResult, error)
	Encoder(keyPair interface{}, param IParam) []byte
}
