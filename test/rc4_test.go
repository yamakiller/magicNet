package test

import (
	"crypto/rc4"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestRC4(t *testing.T) {
	rc4str1 := []byte("xiaowangnidayede") //需要加密的字符串
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, 8998888293223)
	rc4stream, err := rc4.NewCipher(key)
	rc4stream2, err := rc4.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	plaintext := make([]byte, len(rc4str1)) //

	rc4stream.XORKeyStream(plaintext, rc4str1)

	fmt.Printf("%s\n", rc4str1)
	fmt.Printf("%s\n", plaintext)

	rc4stream2.XORKeyStream(rc4str1, plaintext)

	fmt.Printf("%s\n", rc4str1)
	fmt.Printf("%s\n", plaintext)
}
