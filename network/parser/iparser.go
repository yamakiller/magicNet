package parser

import "bytes"

// IParser Parser interface
type IParser interface {
	Analysis(data *bytes.Buffer) (string, uint64, []byte, error)
	Assemble(version int32, handle uint64, agreeName string, data []byte, length int32) []byte
}
