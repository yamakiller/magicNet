package dh64

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"math"
	"math/big"
)

const (
	//DefaultP default p
	DefaultP uint64 = 0xffffffffffffffc5
	//DefaultG default g
	DefaultG uint64 = 5
)

//KeyExchange key exchange
type KeyExchange struct {
	P uint64
	G uint64
}

func (slf *KeyExchange) mulModP(a, b uint64) uint64 {
	var m uint64
	for b > 0 {
		if b&1 > 0 {
			t := slf.P - a
			if m >= t {
				m -= t
			} else {
				m += a
			}
		}
		if a >= slf.P-a {
			a = a*2 - slf.P
		} else {
			a = a * 2
		}
		b >>= 1
	}
	return m
}

func (slf *KeyExchange) powModP(a, b uint64) uint64 {
	if b == 1 {
		return a
	}
	t := slf.powModP(a, b>>1)
	t = slf.mulModP(t, t)
	if b%2 > 0 {
		t = slf.mulModP(t, a)
	}
	return t
}

func (slf *KeyExchange) powmodp(a uint64, b uint64) uint64 {
	if a == 0 {
		panic("DH64 zero public key")
	}
	if b == 0 {
		panic("DH64 zero private key")
	}
	if a > slf.P {
		a %= slf.P
	}
	return slf.powModP(a, b)
}

//PublicKey private key to public key
func (slf *KeyExchange) PublicKey(privateKey uint64) uint64 {
	return slf.powmodp(slf.G, privateKey)
}

//KeyPair Generate public key key pair
func (slf *KeyExchange) KeyPair() (privateKey, publicKey uint64) {
	tmp, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	a := uint64(tmp.Int64())
	tmp, _ = rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	b := uint64(tmp.Int64()) + 1
	privateKey = (a << 32) | b
	publicKey = slf.PublicKey(privateKey)
	return
}

//Secret Generate Secret to uint64
func (slf *KeyExchange) Secret(privateKey, anotherPublicKey uint64) uint64 {
	return slf.powmodp(anotherPublicKey, privateKey)
}

//SecretToString Generate Secret to String
func (slf *KeyExchange) SecretToString(privateKey, anotherPublicKey uint64) string {
	secret := slf.Secret(privateKey, anotherPublicKey)
	tmpBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(tmpBytes, secret)
	return hex.EncodeToString(tmpBytes)
}

//
//1.随机生成一对64位密钥（私钥 + 公钥) myPrivateKey, myPublicKey := dh64.KeyPair()
//2.公锁发送给客户端
//3.等待客户端的公锁
//4.根据客户端的公锁 + 服务器的私锁，计算出密钥：secert := dh64.Secert(myPrivateKey, anotherPublicKey);
//5.客户端既按照此方法计算出密钥
