package encryption

//INetEncryption network encryption insterface
type INetEncryption interface {
	Cipher(key []byte) error
	Encrypt(dst, src []byte) error
	Decode(dst, src []byte) error
	Destory()
}
