// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: agreement.proto

package actor

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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
	return fileDescriptor_92312855cad1e50f, []int{0}
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

func (m *PID) GetID() uint32 {
	if m != nil {
		return m.ID
	}
	return 0
}

// user messages
type Kill struct {
}

func (m *Kill) Reset()      { *m = Kill{} }
func (*Kill) ProtoMessage() {}
func (*Kill) Descriptor() ([]byte, []int) {
	return fileDescriptor_92312855cad1e50f, []int{1}
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
	return fileDescriptor_92312855cad1e50f, []int{2}
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
	return fileDescriptor_92312855cad1e50f, []int{3}
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
	return fileDescriptor_92312855cad1e50f, []int{4}
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
	return fileDescriptor_92312855cad1e50f, []int{5}
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
	return fileDescriptor_92312855cad1e50f, []int{6}
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
	return fileDescriptor_92312855cad1e50f, []int{7}
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
	return fileDescriptor_92312855cad1e50f, []int{8}
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
	return fileDescriptor_92312855cad1e50f, []int{9}
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

func init() { proto.RegisterFile("agreement.proto", fileDescriptor_92312855cad1e50f) }

var fileDescriptor_92312855cad1e50f = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x90, 0xbb, 0x4e, 0x02, 0x41,
	0x14, 0x86, 0x67, 0xb8, 0x7b, 0x8c, 0xa8, 0x5b, 0x11, 0x63, 0x0e, 0x64, 0x62, 0x41, 0xc3, 0x92,
	0xa8, 0x95, 0xa5, 0xd9, 0x66, 0x63, 0x43, 0x56, 0x8c, 0xb1, 0x32, 0xcb, 0xee, 0xb8, 0x4c, 0xc2,
	0xee, 0x90, 0x61, 0x90, 0x96, 0x47, 0xf0, 0x15, 0xec, 0x7c, 0x14, 0x4b, 0x4a, 0x0a, 0x0b, 0x19,
	0x1a, 0x4b, 0x1e, 0xc1, 0xec, 0xe0, 0xa5, 0xb1, 0xb1, 0xfb, 0xce, 0xf9, 0xff, 0x6f, 0x26, 0x39,
	0xb0, 0x1f, 0x26, 0x8a, 0xf3, 0x94, 0x67, 0xda, 0x1d, 0x2b, 0xa9, 0xa5, 0x53, 0x0e, 0x23, 0x2d,
	0xd5, 0x51, 0x27, 0x11, 0x7a, 0x38, 0x1d, 0xb8, 0x91, 0x4c, 0xbb, 0x89, 0x4c, 0x64, 0xd7, 0xa6,
	0x83, 0xe9, 0x83, 0x9d, 0xec, 0x60, 0x69, 0x6b, 0xb1, 0x26, 0x14, 0x7b, 0xbe, 0xe7, 0xd4, 0xa1,
	0xe0, 0x7b, 0x0d, 0xda, 0xa2, 0xed, 0xbd, 0xa0, 0xe0, 0x7b, 0x17, 0xb5, 0xcd, 0x73, 0x93, 0xcc,
	0xdf, 0x5a, 0x84, 0x55, 0xa0, 0x74, 0x25, 0x46, 0x23, 0xd6, 0x81, 0xf2, 0x6d, 0xa8, 0xa3, 0xa1,
	0x73, 0x02, 0xd5, 0x59, 0x0e, 0x5c, 0xd9, 0xfe, 0xee, 0x29, 0xb8, 0xf6, 0x67, 0xb7, 0xe7, 0x7b,
	0xc1, 0x77, 0xc4, 0xba, 0x50, 0xbd, 0xc9, 0x66, 0xff, 0x10, 0xee, 0x00, 0xfa, 0x5c, 0xa5, 0x22,
	0x0b, 0x35, 0x8f, 0x9d, 0x63, 0x28, 0xce, 0x86, 0xf2, 0x8f, 0x7e, 0xbe, 0x76, 0x3a, 0xe0, 0x84,
	0x71, 0xac, 0xf8, 0x64, 0x72, 0xaf, 0x7f, 0x9c, 0x46, 0xa1, 0x45, 0xdb, 0xb5, 0xe0, 0xf0, 0x2b,
	0xf9, 0x7d, 0x8c, 0x1d, 0x40, 0x3d, 0xe0, 0x11, 0x17, 0x8f, 0xbc, 0x2f, 0x52, 0x2e, 0xa7, 0x9a,
	0x01, 0xd4, 0xae, 0xb5, 0x1c, 0x8f, 0x45, 0x96, 0xb0, 0x1d, 0xa8, 0x5a, 0xe6, 0xf1, 0x16, 0x43,
	0x95, 0x3b, 0x15, 0x28, 0xe5, 0xdb, 0xcb, 0xf3, 0xc5, 0x0a, 0xc9, 0x72, 0x85, 0x64, 0xb3, 0x42,
	0x32, 0x37, 0x48, 0x5f, 0x0c, 0xd2, 0x57, 0x83, 0x74, 0x61, 0x90, 0xbe, 0x1b, 0xa4, 0x1f, 0x06,
	0xc9, 0xc6, 0x20, 0x7d, 0x5a, 0x23, 0x59, 0xac, 0x91, 0x2c, 0xd7, 0x48, 0x06, 0x15, 0x7b, 0xdc,
	0xb3, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1c, 0x9d, 0x8a, 0x3a, 0xa5, 0x01, 0x00, 0x00,
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
	if this.ID != that1.ID {
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
	if m.ID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintAgreement(dAtA, i, uint64(m.ID))
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
		i = encodeVarintAgreement(dAtA, i, uint64(m.Watcher.Size()))
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
		i = encodeVarintAgreement(dAtA, i, uint64(m.Watcher.Size()))
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
		i = encodeVarintAgreement(dAtA, i, uint64(m.Who.Size()))
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

func encodeVarintAgreement(dAtA []byte, offset int, v uint64) int {
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
	if m.ID != 0 {
		n += 1 + sovAgreement(uint64(m.ID))
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
		n += 1 + l + sovAgreement(uint64(l))
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
		n += 1 + l + sovAgreement(uint64(l))
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
		n += 1 + l + sovAgreement(uint64(l))
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

func sovAgreement(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAgreement(x uint64) (n int) {
	return sovAgreement(uint64((x << 1) ^ uint64((int64(x) >> 63))))
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
func valueToStringAgreement(v interface{}) string {
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
				return ErrIntOverflowAgreement
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
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAgreement
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAgreement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgreement
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAgreement
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
				return ErrIntOverflowAgreement
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
			skippy, err := skipAgreement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgreement
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAgreement
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
				return ErrIntOverflowAgreement
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
					return ErrIntOverflowAgreement
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
				return ErrInvalidLengthAgreement
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAgreement
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
			skippy, err := skipAgreement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgreement
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAgreement
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
				return ErrIntOverflowAgreement
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
					return ErrIntOverflowAgreement
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
				return ErrInvalidLengthAgreement
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAgreement
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
			skippy, err := skipAgreement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgreement
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAgreement
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
				return ErrIntOverflowAgreement
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
					return ErrIntOverflowAgreement
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
				return ErrInvalidLengthAgreement
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAgreement
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
					return ErrIntOverflowAgreement
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
			skippy, err := skipAgreement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgreement
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAgreement
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
				return ErrIntOverflowAgreement
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
			skippy, err := skipAgreement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgreement
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAgreement
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
				return ErrIntOverflowAgreement
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
			skippy, err := skipAgreement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgreement
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAgreement
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
				return ErrIntOverflowAgreement
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
			skippy, err := skipAgreement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgreement
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAgreement
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
				return ErrIntOverflowAgreement
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
			skippy, err := skipAgreement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgreement
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAgreement
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
				return ErrIntOverflowAgreement
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
			skippy, err := skipAgreement(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgreement
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAgreement
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
func skipAgreement(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAgreement
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
					return 0, ErrIntOverflowAgreement
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
					return 0, ErrIntOverflowAgreement
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
				return 0, ErrInvalidLengthAgreement
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthAgreement
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowAgreement
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
				next, err := skipAgreement(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthAgreement
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
	ErrInvalidLengthAgreement = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAgreement   = fmt.Errorf("proto: integer overflow")
)
