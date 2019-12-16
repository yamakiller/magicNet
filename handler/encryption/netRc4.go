package encryption

import "crypto/rc4"

//NetRC4Encrypt rc4
type NetRC4Encrypt struct {
	_edh *rc4.Cipher
}

//Cipher Encryptor key construction
func (slf *NetRC4Encrypt) Cipher(key []byte) (err error) {
	slf._edh, err = rc4.NewCipher(key)
	return
}

//Encrypt Encrypt data
func (slf *NetRC4Encrypt) Encrypt(dst, src []byte) error {
	slf._edh.XORKeyStream(dst, src)
	return nil
}

//Decode Decode data
func (slf *NetRC4Encrypt) Decode(dst, src []byte) error {
	slf._edh.XORKeyStream(dst, src)
	return nil
}

//Destory destory rc4 encryptor
func (slf *NetRC4Encrypt) Destory() {
	if slf._edh != nil {
		slf._edh.Reset()
	}
}
