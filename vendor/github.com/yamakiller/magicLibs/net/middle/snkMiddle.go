package middle

import (
	"context"
	"crypto/rc4"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/yamakiller/magicLibs/encryption/dh64"
	"github.com/yamakiller/magicLibs/util"
)

//握手协议中间件
type state uint8

const (
	//SNRequest 申请SN
	snApt = 0x01
	//SNConfirm 确认SN
	snAck = 0x02
	//SNRequestSuccess 回复申请成功
	snAptSuccess = 0x11
	//SNRequestFail 回复申请失败
	snAptFail = 0x21
	//SNConfirmSuccess 回复确认成功
	snAckSuccess = 0x12
	//SNConfirmFail 回复确认失败s
	snAckFail = 0x22
)

const (
	stateIdle     = state(0)
	stateAccepted = state(1)
	stateAcked    = state(2)
)

const (
	protoLength = 9
)

var (
	//ErrUnKnownProto 未知协议
	ErrUnKnownProto = errors.New("unknown protocol")
	//ErrUnAuthorized 越权请求
	ErrUnAuthorized = errors.New("unauthorized request")
	//ErrAccidentalRelease 意外的资源释放
	ErrAccidentalRelease = errors.New("accidental release")
	//ErrHandshakeFailure 握手失败
	ErrHandshakeFailure = errors.New("Handshake failure")
)

//SpawnSnkMiddleServe 创建Snk服务端中间件
func SpawnSnkMiddleServe(seed int, p, g uint64, timeout time.Duration) *SnkMiddleServe {
	rand.Seed(time.Now().UnixNano())
	max := math.MaxUint32 / 8
	min := seed
	if min > max {
		min = max / 2
	}

	return &SnkMiddleServe{
		_swap: dh64.KeyExchange{
			P: p,
			G: g,
		},
		_handShakeTimeout: timeout,
		_swapBuffer:       make([]byte, 32),
		_sns: snks{_mask: 0x7FFFFFF,
			_max:   math.MaxUint16,
			_snint: uint32(rand.Intn(max-min) + min),
			_cap:   uint32(32),
			_ssed:  make(map[uint32]*snkData),
			_ss:    make([]*snkData, 32),
		},
	}
}

//SpawnSnkMiddleCli 创建Snk客户端中间件
func SpawnSnkMiddleCli(p, g uint64) *SnkMiddleCli {
	return &SnkMiddleCli{_swap: dh64.KeyExchange{
		P: p,
		G: g,
	}}
}

//SnkMiddleServe 服务端
type SnkMiddleServe struct {
	_swap             dh64.KeyExchange
	_handShakeTimeout time.Duration
	_swapBuffer       []byte
	_prvTime          int64
	_sns              snks
}

//Subscribe 订阅资源
func (slf *SnkMiddleServe) Subscribe(d []byte,
	c *net.UDPConn,
	a *net.UDPAddr) (interface{}, error) {

	if len(d) != protoLength {
		return nil, fmt.Errorf("protocol exception %d", len(d))
	}

	cmd := d[0]
	switch cmd {
	case snApt:
		err := slf.apt(d[1:], a)
		if err != nil {
			slf._swapBuffer[0] = snAptFail
			c.WriteToUDP(slf._swapBuffer[:1], a)
		} else {
			slf._swapBuffer[0] = snAptSuccess
			c.WriteToUDP(slf._swapBuffer[:13], a)
		}
		return nil, err
	case snAck:
		if snk, ok := slf.ack(d[1:], a); ok {
			slf._swapBuffer[0] = snAckSuccess
			c.WriteToUDP(slf._swapBuffer[:1], a)
			return snk, nil
		}

		slf._swapBuffer[0] = snAckFail
		c.WriteToUDP(slf._swapBuffer[:1], a)

		return nil, ErrHandshakeFailure
	default:
		return nil, ErrUnKnownProto
	}

}

//UnSubscribe 取消一个订阅资源
func (slf *SnkMiddleServe) UnSubscribe(d uint32) {
	slf._sns.Lock()
	defer slf._sns.Unlock()
	slf._sns.Remove(d)
}

func (slf *SnkMiddleServe) apt(packet []byte, a *net.UDPAddr) error {
	epubKey := binary.BigEndian.Uint64(packet)
	var snk *snkData
	var err error
	//1.需要这个地址，这个KEY 是否提交过
	slf._sns.Lock()
	snk = slf._sns.GetValue(epubKey, a)
	if snk == nil {
		snk, err = slf._sns.Next()
		if err != nil {
			slf._sns.Unlock()
			return err
		}
	}
	slf._sns.Unlock()

	if snk._epubKey != epubKey {
		if snk._state != stateIdle {
			return ErrUnAuthorized
		}

		prvKey, pubKey := slf._swap.KeyPair()
		snk._death = util.ToTimestamp(time.Now().
			Add(slf._handShakeTimeout))
		snk._secret = slf._swap.Secret(prvKey, epubKey)
		snk._lpubKey = pubKey
		binary.BigEndian.PutUint64(slf._swapBuffer, snk._secret)

		cipher, err := rc4.NewCipher(slf._swapBuffer[:8])
		if err != nil {
			slf._sns.Lock()
			slf._sns.Remove(snk._conv)
			slf._sns.Unlock()
			return err
		}

		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&snk._rc4)),
			unsafe.Pointer(cipher))
		snk._state = stateAccepted
	} else {
		if snk._state == stateIdle {
			return ErrUnAuthorized
		}
	}

	p := (*rc4.Cipher)(atomic.
		LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&snk._rc4))))
	if p == nil {
		slf._sns.Lock()
		slf._sns.Remove(snk._conv)
		slf._sns.Unlock()
		return ErrAccidentalRelease
	}

	binary.BigEndian.PutUint32(slf._swapBuffer[5:], snk._conv)
	p.XORKeyStream(slf._swapBuffer[1:], slf._swapBuffer[5:9])
	binary.BigEndian.PutUint64(slf._swapBuffer[5:], snk._lpubKey)
	fmt.Println("success 1")
	return nil
}

func (slf *SnkMiddleServe) ack(d []byte,
	a *net.UDPAddr) (uint32, bool) {

	snkid := binary.BigEndian.Uint32(d)
	slf._sns.Lock()
	snk := slf._sns.Get(snkid)
	if snk == nil {
		slf._sns.Unlock()
		return 0, false
	}

	slf._sns.Unlock()

	p := (*rc4.Cipher)(atomic.
		LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&snk._rc4))))

	if p == nil || snk._state == stateIdle {
		return 0, false
	}

	token := binary.BigEndian.Uint32(d[4:])
	binary.BigEndian.PutUint32(d[4:], token)
	tmpBuf := make([]byte, 4)
	p.Reset()
	p.XORKeyStream(tmpBuf, d[4:])
	token = binary.BigEndian.Uint32(tmpBuf)
	if token != snkid {

		return 0, false
	}
	snk._state = stateAcked

	return snkid, true
}

//Update 维护
func (slf *SnkMiddleServe) Update() {
	var keys []uint32
	current := util.Timestamp()
	//限制Update的频率大于100毫秒一次
	if (current - slf._prvTime) < 100 {
		return
	}
	slf._prvTime = current

	slf._sns.Lock()
	keys = slf._sns.GetKeys()
	slf._sns.Unlock()
	if keys == nil || len(keys) == 0 {
		return
	}

	var p *snkData
	for _, v := range keys {
		slf._sns.Lock()
		p = slf._sns.Get(v)
		slf._sns.Unlock()
		if p == nil {
			continue
		}

		if p._state == stateAcked {
			continue
		}

		if current > p._death {
			slf._sns.Lock()
			slf._sns.Remove(v)
			slf._sns.Unlock()
			continue
		}

	}
}

//SnkMiddleCli snk资源订阅客户端
type SnkMiddleCli struct {
	_swap     dh64.KeyExchange
	_q        chan interface{}
	_closed   bool
	_interval time.Duration
	_epubKey  uint64
	_lpubKey  uint64
	_prvKey   uint64
}

//Subscribe 订阅资源
func (slf *SnkMiddleCli) Subscribe(c *net.UDPConn,
	a *net.UDPAddr, timeout time.Duration) (interface{}, error) {
	slf._prvKey, slf._lpubKey = slf._swap.KeyPair()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	slf._q = make(chan interface{}, 1)
	slf._closed = false
	slf._interval = timeout / 10

	go func() {
		defer func() {
			close(slf._q)
		}()

		conv, err := slf.doSubscribe(c, slf._interval)
		if err != nil {
			slf._q <- err
			return
		}

		conv, err = slf.doAck(c, conv, slf._interval)
		if err != nil {
			slf._q <- err
			return
		}

		slf._q <- conv
	}()

	select {
	case msg := <-slf._q:
		if err, ok := msg.(error); ok {
			return nil, err
		}
		slf._closed = true
		return msg, nil
	case <-ctx.Done():
		slf._closed = true
		c.Close()
		return nil, errors.New("timeout")
	}
}

func (slf *SnkMiddleCli) doSubscribe(c *net.UDPConn, timeout time.Duration) (uint32, error) {

	var (
		n, chk int
		conv   uint32
		err    error
		cmd    uint8
	)

	tmpBuf := make([]byte, 32)

	for {
		tmpBuf[0] = snApt
		binary.BigEndian.PutUint64(tmpBuf[1:], slf._lpubKey)

		_, err = c.Write(tmpBuf[:9])
		if err != nil {
			return 0, err
		}

		fmt.Printf("\n")
		fmt.Println("timeout:", timeout)
		c.SetReadDeadline(time.Now().Add(timeout))
		n, _, err = c.ReadFromUDP(tmpBuf)
		if err != nil {
			fmt.Println(err)
			goto error_check
		}

		cmd = tmpBuf[0]
		if cmd == snAptFail {
			return 0, errors.New("subscription failed")
		}

		if cmd != snAptSuccess || n != 13 {
			goto error_check
		}

		slf._epubKey = binary.BigEndian.Uint64(tmpBuf[5:])
		conv = binary.BigEndian.Uint32(tmpBuf[1:])

		return conv, nil
	error_check:
		chk++
		if slf._closed {
			return 0, err
		}
	}
}

func (slf *SnkMiddleCli) doAck(c *net.UDPConn,
	conv uint32,
	timeout time.Duration) (uint32, error) {
	//to cipher
	cipher, err := slf.toCipher(slf._swap.Secret(slf._prvKey, slf._epubKey))
	if err != nil {
		return 0, err
	}

	var (
		chk           int
		n             int
		cmd           uint8
		unconv, token uint32
	)

	tmpBuf := make([]byte, 32)
	unconv = slf.toDecode(cipher, conv)

	for {
		tmpBuf[0] = snAck
		cipher.Reset()
		token = slf.toToken(cipher, unconv)
		binary.BigEndian.PutUint32(tmpBuf[1:], unconv)
		binary.BigEndian.PutUint32(tmpBuf[5:], token)

		_, err = c.Write(tmpBuf[:9])
		if err != nil {
			return 0, err
		}

		fmt.Printf("\n")
		fmt.Println("2 timeout:", timeout)
		c.SetReadDeadline(time.Now().Add(timeout))
		n, _, err = c.ReadFromUDP(tmpBuf)
		if err != nil {
			fmt.Println("2,", err)
			goto error_check
		}

		cmd = tmpBuf[0]
		if cmd == snAckFail {
			return 0, errors.New("subscription ack failed")
		}

		if cmd != snAckSuccess || n != 1 {
			goto error_check
		}

		return unconv, nil
	error_check:
		chk++
		if slf._closed {
			return 0, err
		}
	}
}

func (slf *SnkMiddleCli) toDecode(cipher *rc4.Cipher,
	conv uint32) (unconv uint32) {
	tmpBuf := make([]byte, 8)
	binary.BigEndian.PutUint32(tmpBuf[4:], conv)
	cipher.XORKeyStream(tmpBuf, tmpBuf[4:])
	unconv = binary.BigEndian.Uint32(tmpBuf)
	return
}

func (slf *SnkMiddleCli) toToken(cipher *rc4.Cipher, unconv uint32) (token uint32) {
	tmpBuf := make([]byte, 8)
	binary.BigEndian.PutUint32(tmpBuf[4:], unconv)
	cipher.XORKeyStream(tmpBuf, tmpBuf[4:])
	token = binary.BigEndian.Uint32(tmpBuf)
	return
}

func (slf *SnkMiddleCli) toCipher(secret uint64) (*rc4.Cipher, error) {
	tmpBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(tmpBuf, secret)

	return rc4.NewCipher(tmpBuf)
}

type snks struct {
	_mask  uint32
	_max   uint32 //2的幂
	_snint uint32
	_sz    int
	_cap   uint32
	_ssed  map[uint32]*snkData
	_ss    []*snkData
	sync.Mutex
}

func (slf *snks) Next() (*snkData, error) {
	var i uint32
	for {

		for i = 0; i < slf._cap; i++ {
			key := ((i + slf._snint) & slf._mask)
			if key == 0 {
				key = 1
			}
			hash := key & (slf._cap - 1)
			if slf._ss[hash] == nil {
				slf._snint = key + 1
				slf._ss[hash] = &snkData{
					_conv: key,
				}
				slf._ssed[key] = slf._ss[hash]
				slf._sz++
				return slf._ss[hash], nil
			}
		}

		newCap := slf._cap * 2
		if newCap > slf._max {
			newCap = slf._max
		}

		if newCap == slf._cap {
			return nil, errors.New("full")
		}

		slf._ss = append(slf._ss, make([]*snkData, newCap-slf._cap)...)
		for i = 0; i < slf._cap; i++ {
			if slf._ss[i] == nil {
				continue
			}

			hash := slf._ss[i]._conv & uint32(newCap-1)
			if hash == i {
				continue
			}

			tmp := slf._ss[i]
			slf._ss[hash] = tmp
			slf._ss[hash] = nil
		}
		slf._cap = newCap
	}
}

func (slf *snks) GetValue(pubKey uint64, a *net.UDPAddr) *snkData {
	for _, v := range slf._ssed {
		if v._epubKey == pubKey &&
			v._addr.IP.Equal(a.IP) &&
			v._addr.Port == a.Port {
			return v
		}
	}

	return nil
}

//Get 获取一个元素
func (slf *snks) Get(key uint32) *snkData {
	hash := key & uint32(slf._cap-1)
	if slf._ss[hash] != nil && slf._ss[hash]._conv == key {
		return slf._ss[hash]
	}
	return nil
}

func (slf *snks) GetKeys() []uint32 {
	var keys []uint32
	if slf._sz > 0 {
		keys = make([]uint32, slf._sz)

		i := 0
		for k := range slf._ssed {
			keys[i] = k
			i++
		}
	}

	return keys
}

//Remove 删除一个元素
func (slf *snks) Remove(key uint32) bool {
	hash := uint32(key) & uint32(slf._cap-1)
	if slf._ss[hash] != nil && slf._ss[hash]._conv == key {
		slf._ss[hash]._state = stateIdle
		delete(slf._ssed, key)
		slf._ss[hash] = nil
		slf._sz--
		return true
	}

	return false
}

type snkData struct {
	_conv    uint32
	_secret  uint64
	_epubKey uint64
	_lpubKey uint64
	_rc4     *rc4.Cipher
	_addr    net.UDPAddr
	_death   int64
	_state   state
}
