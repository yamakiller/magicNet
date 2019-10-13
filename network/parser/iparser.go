package parser

import "bytes"

// IParser Parser interface
type IParser interface {
	Analysis(keyPair interface{}, data *bytes.Buffer) (string, uint64, uint32, []byte, error)
	Assemble(keyPair interface{}, version int32, handle uint64, Serial uint32, agreeName string, data []byte, length int32) []byte
}
