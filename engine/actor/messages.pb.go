// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: messages.proto

package actor

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

func (m *PID) Reset()      { *m = PID{} }
func (*PID) ProtoMessage() {}
func (*PID) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{0}
}
func (m *PID) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PID.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PID.Merge(m, src)
}
func (m *PID) XXX_Size() int {
	return m.Size()
}
func (m *PID) XXX_DiscardUnknown() {
	xxx_messageInfo_PID.DiscardUnknown(m)
}

var xxx_messageInfo_PID proto.InternalMessageInfo

func (m *PID) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

// user messages
type Kill struct {
}

func (m *Kill) Reset()      { *m = Kill{} }
func (*Kill) ProtoMessage() {}
func (*Kill) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{1}
}
func (m *Kill) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Kill) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Kill.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Kill) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Kill.Merge(m, src)
}
func (m *Kill) XXX_Size() int {
	return m.Size()
}
func (m *Kill) XXX_DiscardUnknown() {
	xxx_messageInfo_Kill.DiscardUnknown(m)
}

var xxx_messageInfo_Kill proto.InternalMessageInfo

// system messages
type Watch struct {
	Watcher *PID `protobuf:"bytes,1,opt,name=watcher,proto3" json:"watcher,omitempty"`
}

func (m *Watch) Reset()      { *m = Watch{} }
func (*Watch) ProtoMessage() {}
func (*Watch) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{2}
}
func (m *Watch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Watch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Watch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Watch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Watch.Merge(m, src)
}
func (m *Watch) XXX_Size() int {
	return m.Size()
}
func (m *Watch) XXX_DiscardUnknown() {
	xxx_messageInfo_Watch.DiscardUnknown(m)
}

var xxx_messageInfo_Watch proto.InternalMessageInfo

func (m *Watch) GetWatcher() *PID {
	if m != nil {
		return m.Watcher
	}
	return nil
}

type Unwatch struct {
	Watcher *PID `protobuf:"bytes,1,opt,name=watcher,proto3" json:"watcher,omitempty"`
}

func (m *Unwatch) Reset()      { *m = Unwatch{} }
func (*Unwatch) ProtoMessage() {}
func (*Unwatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{3}
}
func (m *Unwatch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Unwatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Unwatch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Unwatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Unwatch.Merge(m, src)
}
func (m *Unwatch) XXX_Size() int {
	return m.Size()
}
func (m *Unwatch) XXX_DiscardUnknown() {
	xxx_messageInfo_Unwatch.DiscardUnknown(m)
}

var xxx_messageInfo_Unwatch proto.InternalMessageInfo

func (m *Unwatch) GetWatcher() *PID {
	if m != nil {
		return m.Watcher
	}
	return nil
}

type Terminated struct {
	Who               *PID `protobuf:"bytes,1,opt,name=who,proto3" json:"who,omitempty"`
	AddressTerminated bool `protobuf:"varint,2,opt,name=address_terminated,json=addressTerminated,proto3" json:"address_terminated,omitempty"`
}

func (m *Terminated) Reset()      { *m = Terminated{} }
func (*Terminated) ProtoMessage() {}
func (*Terminated) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{4}
}
func (m *Terminated) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Terminated) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Terminated.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Terminated) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Terminated.Merge(m, src)
}
func (m *Terminated) XXX_Size() int {
	return m.Size()
}
func (m *Terminated) XXX_DiscardUnknown() {
	xxx_messageInfo_Terminated.DiscardUnknown(m)
}

var xxx_messageInfo_Terminated proto.InternalMessageInfo

func (m *Terminated) GetWho() *PID {
	if m != nil {
		return m.Who
	}
	return nil
}

func (m *Terminated) GetAddressTerminated() bool {
	if m != nil {
		return m.AddressTerminated
	}
	return false
}

type ReceiveTimeout struct {
}

func (m *ReceiveTimeout) Reset()      { *m = ReceiveTimeout{} }
func (*ReceiveTimeout) ProtoMessage() {}
func (*ReceiveTimeout) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{5}
}
func (m *ReceiveTimeout) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReceiveTimeout) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReceiveTimeout.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReceiveTimeout) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReceiveTimeout.Merge(m, src)
}
func (m *ReceiveTimeout) XXX_Size() int {
	return m.Size()
}
func (m *ReceiveTimeout) XXX_DiscardUnknown() {
	xxx_messageInfo_ReceiveTimeout.DiscardUnknown(m)
}

var xxx_messageInfo_ReceiveTimeout proto.InternalMessageInfo

type Stopping struct {
}

func (m *Stopping) Reset()      { *m = Stopping{} }
func (*Stopping) ProtoMessage() {}
func (*Stopping) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{6}
}
func (m *Stopping) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Stopping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Stopping.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Stopping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stopping.Merge(m, src)
}
func (m *Stopping) XXX_Size() int {
	return m.Size()
}
func (m *Stopping) XXX_DiscardUnknown() {
	xxx_messageInfo_Stopping.DiscardUnknown(m)
}

var xxx_messageInfo_Stopping proto.InternalMessageInfo

type Stopped struct {
}

func (m *Stopped) Reset()      { *m = Stopped{} }
func (*Stopped) ProtoMessage() {}
func (*Stopped) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{7}
}
func (m *Stopped) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Stopped) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Stopped.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Stopped) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stopped.Merge(m, src)
}
func (m *Stopped) XXX_Size() int {
	return m.Size()
}
func (m *Stopped) XXX_DiscardUnknown() {
	xxx_messageInfo_Stopped.DiscardUnknown(m)
}

var xxx_messageInfo_Stopped proto.InternalMessageInfo

type Started struct {
}

func (m *Started) Reset()      { *m = Started{} }
func (*Started) ProtoMessage() {}
func (*Started) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{8}
}
func (m *Started) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Started) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Started.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Started) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Started.Merge(m, src)
}
func (m *Started) XXX_Size() int {
	return m.Size()
}
func (m *Started) XXX_DiscardUnknown() {
	xxx_messageInfo_Started.DiscardUnknown(m)
}

var xxx_messageInfo_Started proto.InternalMessageInfo

type Stop struct {
}

func (m *Stop) Reset()      { *m = Stop{} }
func (*Stop) ProtoMessage() {}
func (*Stop) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{9}
}
func (m *Stop) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Stop) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Stop.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Stop) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stop.Merge(m, src)
}
func (m *Stop) XXX_Size() int {
	return m.Size()
}
func (m *Stop) XXX_DiscardUnknown() {
	xxx_messageInfo_Stop.DiscardUnknown(m)
}

var xxx_messageInfo_Stop proto.InternalMessageInfo

func init() {
	proto.RegisterType((*PID)(nil), "actor.PID")
	proto.RegisterType((*Kill)(nil), "actor.Kill")
	proto.RegisterType((*Watch)(nil), "actor.Watch")
	proto.RegisterType((*Unwatch)(nil), "actor.Unwatch")
	proto.RegisterType((*Terminated)(nil), "actor.Terminated")
	proto.RegisterType((*ReceiveTimeout)(nil), "actor.ReceiveTimeout")
	proto.RegisterType((*Stopping)(nil), "actor.Stopping")
	proto.RegisterType((*Stopped)(nil), "actor.Stopped")
	proto.RegisterType((*Started)(nil), "actor.Started")
	proto.RegisterType((*Stop)(nil), "actor.Stop")
}

func init() { proto.RegisterFile("messages.proto", fileDescriptor_4dc296cbfe5ffcd5) }

var fileDescriptor_4dc296cbfe5ffcd5 = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x90, 0xbd, 0x4e, 0x32, 0x41,
	0x14, 0x86, 0x67, 0xf8, 0xff, 0xce, 0x17, 0x89, 0x6e, 0x45, 0x8c, 0x39, 0x90, 0x89, 0x05, 0x0d,
	0x4b, 0xa2, 0x56, 0x96, 0x86, 0x86, 0xd8, 0x90, 0x15, 0x63, 0xac, 0xcc, 0xb2, 0x33, 0x2e, 0x93,
	0xb0, 0x0c, 0x99, 0x1d, 0xa4, 0xe5, 0x12, 0xbc, 0x05, 0x3b, 0x2f, 0xc5, 0x92, 0x92, 0xc2, 0x42,
	0x86, 0xc6, 0x72, 0x2f, 0xc1, 0xec, 0xe0, 0x4f, 0x63, 0x63, 0xf7, 0x9c, 0xf3, 0xbe, 0xcf, 0x4c,
	0x72, 0xa0, 0x9e, 0x88, 0x34, 0x0d, 0x63, 0x91, 0xfa, 0x33, 0xad, 0x8c, 0xf2, 0xca, 0x61, 0x64,
	0x94, 0x3e, 0xec, 0xc4, 0xd2, 0x8c, 0xe7, 0x23, 0x3f, 0x52, 0x49, 0x37, 0x56, 0xb1, 0xea, 0xba,
	0x74, 0x34, 0xbf, 0x77, 0x93, 0x1b, 0x1c, 0xed, 0x2c, 0xd6, 0x84, 0xe2, 0xa0, 0xdf, 0xf3, 0xea,
	0x50, 0x90, 0xbc, 0x41, 0x5b, 0xb4, 0xbd, 0x17, 0x14, 0x24, 0x3f, 0xaf, 0x65, 0x4f, 0x4d, 0xb2,
	0x7c, 0x6d, 0x11, 0x56, 0x81, 0xd2, 0xa5, 0x9c, 0x4c, 0x58, 0x07, 0xca, 0x37, 0xa1, 0x89, 0xc6,
	0xde, 0x31, 0x54, 0x17, 0x39, 0x08, 0xed, 0xfa, 0xff, 0x4f, 0xc0, 0x77, 0x3f, 0xfb, 0x83, 0x7e,
	0x2f, 0xf8, 0x8a, 0x58, 0x17, 0xaa, 0xd7, 0xd3, 0xc5, 0x1f, 0x84, 0x5b, 0x80, 0xa1, 0xd0, 0x89,
	0x9c, 0x86, 0x46, 0x70, 0xef, 0x08, 0x8a, 0x8b, 0xb1, 0xfa, 0xa5, 0x9f, 0xaf, 0xbd, 0x0e, 0x78,
	0x21, 0xe7, 0x5a, 0xa4, 0xe9, 0x9d, 0xf9, 0x76, 0x1a, 0x85, 0x16, 0x6d, 0xd7, 0x82, 0x83, 0xcf,
	0xe4, 0xe7, 0x31, 0xb6, 0x0f, 0xf5, 0x40, 0x44, 0x42, 0x3e, 0x88, 0xa1, 0x4c, 0x84, 0x9a, 0x1b,
	0x06, 0x50, 0xbb, 0x32, 0x6a, 0x36, 0x93, 0xd3, 0x98, 0xfd, 0x83, 0xaa, 0x63, 0xc1, 0x77, 0x18,
	0xea, 0xdc, 0xa9, 0x40, 0x29, 0xdf, 0x5e, 0x9c, 0xad, 0x36, 0x48, 0xd6, 0x1b, 0x24, 0xd9, 0x06,
	0xc9, 0xd2, 0x22, 0x7d, 0xb6, 0x48, 0x5f, 0x2c, 0xd2, 0x95, 0x45, 0xfa, 0x66, 0x91, 0xbe, 0x5b,
	0x24, 0x99, 0x45, 0xfa, 0xb8, 0x45, 0xb2, 0xda, 0x22, 0x59, 0x6f, 0x91, 0x8c, 0x2a, 0xee, 0xb8,
	0xa7, 0x1f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x5c, 0x03, 0x0f, 0xff, 0xa4, 0x01, 0x00, 0x00,
}

func (this *PID) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PID)
	if !ok {
		that2, ok := that.(PID)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Id != that1.Id {
		return false
	}
	return true
}
func (this *Kill) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Kill)
	if !ok {
		that2, ok := that.(Kill)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (this *Watch) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Watch)
	if !ok {
		that2, ok := that.(Watch)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Watcher.Equal(that1.Watcher) {
		return false
	}
	return true
}
func (this *Unwatch) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Unwatch)
	if !ok {
		that2, ok := that.(Unwatch)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Watcher.Equal(that1.Watcher) {
		return false
	}
	return true
}
func (this *Terminated) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Terminated)
	if !ok {
		that2, ok := that.(Terminated)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Who.Equal(that1.Who) {
		return false
	}
	if this.AddressTerminated != that1.AddressTerminated {
		return false
	}
	return true
}
func (this *ReceiveTimeout) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ReceiveTimeout)
	if !ok {
		that2, ok := that.(ReceiveTimeout)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (this *Stopping) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Stopping)
	if !ok {
		that2, ok := that.(Stopping)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (this *Stopped) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Stopped)
	if !ok {
		that2, ok := that.(Stopped)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (this *Started) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Started)
	if !ok {
		that2, ok := that.(Started)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (this *Stop) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Stop)
	if !ok {
		that2, ok := that.(Stop)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (m *PID) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PID) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintMessages(dAtA, i, uint64(m.Id))
	}
	return i, nil
}

func (m *Kill) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Kill) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *Watch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Watch) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Watcher != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMessages(dAtA, i, uint64(m.Watcher.Size()))
		n1, err1 := m.Watcher.MarshalTo(dAtA[i:])
		if err1 != nil {
			return 0, err1
		}
		i += n1
	}
	return i, nil
}

func (m *Unwatch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Unwatch) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Watcher != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMessages(dAtA, i, uint64(m.Watcher.Size()))
		n2, err2 := m.Watcher.MarshalTo(dAtA[i:])
		if err2 != nil {
			return 0, err2
		}
		i += n2
	}
	return i, nil
}

func (m *Terminated) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Terminated) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Who != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMessages(dAtA, i, uint64(m.Who.Size()))
		n3, err3 := m.Who.MarshalTo(dAtA[i:])
		if err3 != nil {
			return 0, err3
		}
		i += n3
	}
	if m.AddressTerminated {
		dAtA[i] = 0x10
		i++
		if m.AddressTerminated {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func (m *ReceiveTimeout) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReceiveTimeout) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *Stopping) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Stopping) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *Stopped) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Stopped) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *Started) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Started) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func (m *Stop) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Stop) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func encodeVarintMessages(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *PID) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovMessages(uint64(m.Id))
	}
	return n
}

func (m *Kill) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *Watch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Watcher != nil {
		l = m.Watcher.Size()
		n += 1 + l + sovMessages(uint64(l))
	}
	return n
}

func (m *Unwatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Watcher != nil {
		l = m.Watcher.Size()
		n += 1 + l + sovMessages(uint64(l))
	}
	return n
}

func (m *Terminated) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Who != nil {
		l = m.Who.Size()
		n += 1 + l + sovMessages(uint64(l))
	}
	if m.AddressTerminated {
		n += 2
	}
	return n
}

func (m *ReceiveTimeout) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *Stopping) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *Stopped) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *Started) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *Stop) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovMessages(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMessages(x uint64) (n int) {
	return sovMessages(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Kill) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Kill{`,
		`}`,
	}, "")
	return s
}
func (this *Watch) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Watch{`,
		`Watcher:` + strings.Replace(fmt.Sprintf("%v", this.Watcher), "PID", "PID", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Unwatch) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Unwatch{`,
		`Watcher:` + strings.Replace(fmt.Sprintf("%v", this.Watcher), "PID", "PID", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Terminated) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Terminated{`,
		`Who:` + strings.Replace(fmt.Sprintf("%v", this.Who), "PID", "PID", 1) + `,`,
		`AddressTerminated:` + fmt.Sprintf("%v", this.AddressTerminated) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ReceiveTimeout) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ReceiveTimeout{`,
		`}`,
	}, "")
	return s
}
func (this *Stopping) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Stopping{`,
		`}`,
	}, "")
	return s
}
func (this *Stopped) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Stopped{`,
		`}`,
	}, "")
	return s
}
func (this *Started) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Started{`,
		`}`,
	}, "")
	return s
}
func (this *Stop) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Stop{`,
		`}`,
	}, "")
	return s
}
func valueToStringMessages(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *PID) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PID: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PID: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Kill) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Kill: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Kill: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Watch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Watch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Watch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Watcher", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Watcher == nil {
				m.Watcher = &PID{}
			}
			if err := m.Watcher.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Unwatch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Unwatch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Unwatch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Watcher", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Watcher == nil {
				m.Watcher = &PID{}
			}
			if err := m.Watcher.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Terminated) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Terminated: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Terminated: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Who", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthMessages
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMessages
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Who == nil {
				m.Who = &PID{}
			}
			if err := m.Who.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddressTerminated", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.AddressTerminated = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ReceiveTimeout) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ReceiveTimeout: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReceiveTimeout: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Stopping) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Stopping: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Stopping: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Stopped) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Stopped: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Stopped: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Started) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Started: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Started: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Stop) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Stop: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Stop: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipMessages(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessages
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMessages(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMessages
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessages
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMessages
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthMessages
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowMessages
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipMessages(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthMessages
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthMessages = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMessages   = fmt.Errorf("proto: integer overflow")
)
