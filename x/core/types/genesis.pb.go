// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/core/v1/genesis.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	types "github.com/troykessler/hyperlane-cosmos/x/core/01_interchain_security/types"
	types1 "github.com/troykessler/hyperlane-cosmos/x/core/02_post_dispatch/types"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	// ism_genesis
	IsmGenesis *types.GenesisState `protobuf:"bytes,1,opt,name=ism_genesis,json=ismGenesis,proto3" json:"ism_genesis,omitempty"`
	// post_dispatch_genesis
	PostDispatchGenesis  *types1.GenesisState `protobuf:"bytes,2,opt,name=post_dispatch_genesis,json=postDispatchGenesis,proto3" json:"post_dispatch_genesis,omitempty"`
	Mailboxes            []*Mailbox           `protobuf:"bytes,3,rep,name=mailboxes,proto3" json:"mailboxes,omitempty"`
	Messages             []*MailboxMessage    `protobuf:"bytes,4,rep,name=messages,proto3" json:"messages,omitempty"`
	IsmSequence          uint64               `protobuf:"varint,5,opt,name=ism_sequence,json=ismSequence,proto3" json:"ism_sequence,omitempty"`
	PostDispatchSequence uint64               `protobuf:"varint,6,opt,name=post_dispatch_sequence,json=postDispatchSequence,proto3" json:"post_dispatch_sequence,omitempty"`
	AppSequence          uint64               `protobuf:"varint,7,opt,name=app_sequence,json=appSequence,proto3" json:"app_sequence,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_9329350a78ea2d1f, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetIsmGenesis() *types.GenesisState {
	if m != nil {
		return m.IsmGenesis
	}
	return nil
}

func (m *GenesisState) GetPostDispatchGenesis() *types1.GenesisState {
	if m != nil {
		return m.PostDispatchGenesis
	}
	return nil
}

func (m *GenesisState) GetMailboxes() []*Mailbox {
	if m != nil {
		return m.Mailboxes
	}
	return nil
}

func (m *GenesisState) GetMessages() []*MailboxMessage {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *GenesisState) GetIsmSequence() uint64 {
	if m != nil {
		return m.IsmSequence
	}
	return 0
}

func (m *GenesisState) GetPostDispatchSequence() uint64 {
	if m != nil {
		return m.PostDispatchSequence
	}
	return 0
}

func (m *GenesisState) GetAppSequence() uint64 {
	if m != nil {
		return m.AppSequence
	}
	return 0
}

// Mailbox message for genesis state
type MailboxMessage struct {
	MailboxId uint64 `protobuf:"varint,1,opt,name=mailbox_id,json=mailboxId,proto3" json:"mailbox_id,omitempty"`
	MessageId []byte `protobuf:"bytes,2,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
}

func (m *MailboxMessage) Reset()         { *m = MailboxMessage{} }
func (m *MailboxMessage) String() string { return proto.CompactTextString(m) }
func (*MailboxMessage) ProtoMessage()    {}
func (*MailboxMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_9329350a78ea2d1f, []int{1}
}
func (m *MailboxMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MailboxMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MailboxMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MailboxMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MailboxMessage.Merge(m, src)
}
func (m *MailboxMessage) XXX_Size() int {
	return m.Size()
}
func (m *MailboxMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_MailboxMessage.DiscardUnknown(m)
}

var xxx_messageInfo_MailboxMessage proto.InternalMessageInfo

func (m *MailboxMessage) GetMailboxId() uint64 {
	if m != nil {
		return m.MailboxId
	}
	return 0
}

func (m *MailboxMessage) GetMessageId() []byte {
	if m != nil {
		return m.MessageId
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "hyperlane.core.v1.GenesisState")
	proto.RegisterType((*MailboxMessage)(nil), "hyperlane.core.v1.MailboxMessage")
}

func init() { proto.RegisterFile("hyperlane/core/v1/genesis.proto", fileDescriptor_9329350a78ea2d1f) }

var fileDescriptor_9329350a78ea2d1f = []byte{
	// 411 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xbf, 0x8e, 0xda, 0x40,
	0x10, 0xc6, 0x71, 0x20, 0x97, 0xdc, 0x62, 0x45, 0x8a, 0xf3, 0x47, 0x96, 0xa5, 0x73, 0xe0, 0x2a,
	0x1a, 0xd6, 0xe2, 0xb8, 0x22, 0x4d, 0x9a, 0x28, 0x52, 0x44, 0x01, 0x85, 0x49, 0x95, 0xc6, 0x5a,
	0xec, 0x11, 0x5e, 0x05, 0x7b, 0x37, 0x9e, 0x05, 0xe1, 0xb7, 0xc8, 0x03, 0xe5, 0x01, 0x52, 0x52,
	0xa6, 0x8c, 0xe0, 0x45, 0x22, 0xff, 0xc1, 0x60, 0x3b, 0x69, 0x67, 0x7e, 0xdf, 0xb7, 0x33, 0xdf,
	0x0e, 0x79, 0x17, 0xa6, 0x12, 0x92, 0x0d, 0x8b, 0xc1, 0xf1, 0x45, 0x02, 0xce, 0x6e, 0xe2, 0xac,
	0x21, 0x06, 0xe4, 0x48, 0x65, 0x22, 0x94, 0x30, 0x5e, 0x56, 0x00, 0xcd, 0x00, 0xba, 0x9b, 0x58,
	0x77, 0x6d, 0x8d, 0x4a, 0x25, 0x94, 0x0a, 0x6b, 0xda, 0x68, 0xf3, 0x58, 0x41, 0xe2, 0x87, 0x8c,
	0xc7, 0x1e, 0x82, 0xbf, 0x4d, 0xb8, 0x4a, 0x5b, 0xcf, 0x58, 0xe3, 0x86, 0x48, 0x0a, 0x54, 0x5e,
	0xc0, 0x51, 0x32, 0xe5, 0x87, 0x2d, 0xfc, 0xfe, 0x67, 0x97, 0xe8, 0x9f, 0x8b, 0xca, 0x52, 0x31,
	0x05, 0xc6, 0x17, 0xd2, 0xe7, 0x18, 0x79, 0x25, 0x65, 0x6a, 0x03, 0x6d, 0xd4, 0x7f, 0x98, 0xd2,
	0xc6, 0xf0, 0xff, 0x18, 0x85, 0xee, 0x26, 0xf4, 0xda, 0xc9, 0x25, 0x1c, 0xa3, 0xb2, 0x60, 0x30,
	0xf2, 0xa6, 0x36, 0x48, 0xe5, 0xff, 0x24, 0xf7, 0x1f, 0x37, 0xfd, 0x6b, 0x70, 0xcb, 0xf9, 0x55,
	0xd6, 0xfe, 0x54, 0x76, 0xcf, 0x4f, 0xbc, 0x27, 0xb7, 0x11, 0xe3, 0x9b, 0x95, 0xd8, 0x03, 0x9a,
	0xdd, 0x41, 0x77, 0xd4, 0x7f, 0xb0, 0x68, 0x2b, 0x73, 0x3a, 0x2f, 0x18, 0xf7, 0x02, 0x1b, 0x1f,
	0xc8, 0xf3, 0x08, 0x10, 0xd9, 0x1a, 0xd0, 0xec, 0xe5, 0xc2, 0xe1, 0xff, 0x85, 0xf3, 0x82, 0x74,
	0x2b, 0x89, 0x31, 0x24, 0x7a, 0x96, 0x18, 0xc2, 0xf7, 0x2d, 0xc4, 0x3e, 0x98, 0x4f, 0x07, 0xda,
	0xa8, 0xe7, 0x66, 0x29, 0x2e, 0xcb, 0x92, 0xf1, 0x48, 0xde, 0xd6, 0xd7, 0xaf, 0xe0, 0x9b, 0x1c,
	0x7e, 0x7d, 0xbd, 0x50, 0xa5, 0x1a, 0x12, 0x9d, 0x49, 0x79, 0x61, 0x9f, 0x15, 0xc6, 0x4c, 0xca,
	0x33, 0x72, 0xbf, 0x20, 0x2f, 0xea, 0x73, 0x19, 0x77, 0x84, 0x94, 0x9b, 0x79, 0x3c, 0xc8, 0xbf,
	0xaf, 0x57, 0xed, 0x3a, 0x0b, 0xf2, 0x76, 0x41, 0x66, 0xed, 0x2c, 0x7d, 0xdd, 0xbd, 0x2d, 0x2b,
	0xb3, 0xe0, 0xe3, 0xe2, 0xd7, 0xd1, 0xd6, 0x0e, 0x47, 0x5b, 0xfb, 0x73, 0xb4, 0xb5, 0x1f, 0x27,
	0xbb, 0x73, 0x38, 0xd9, 0x9d, 0xdf, 0x27, 0xbb, 0xf3, 0xf5, 0x71, 0xcd, 0x55, 0xb8, 0x5d, 0x51,
	0x5f, 0x44, 0x8e, 0x4a, 0x44, 0xfa, 0x0d, 0x10, 0x37, 0x90, 0x38, 0x55, 0x50, 0x63, 0x5f, 0x60,
	0x24, 0xd0, 0xd9, 0x17, 0x77, 0x97, 0x1f, 0xf2, 0xea, 0x26, 0xbf, 0xb2, 0xe9, 0xdf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x7e, 0x54, 0x83, 0x4f, 0x1e, 0x03, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.AppSequence != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.AppSequence))
		i--
		dAtA[i] = 0x38
	}
	if m.PostDispatchSequence != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.PostDispatchSequence))
		i--
		dAtA[i] = 0x30
	}
	if m.IsmSequence != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.IsmSequence))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Messages) > 0 {
		for iNdEx := len(m.Messages) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Messages[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Mailboxes) > 0 {
		for iNdEx := len(m.Mailboxes) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Mailboxes[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.PostDispatchGenesis != nil {
		{
			size, err := m.PostDispatchGenesis.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.IsmGenesis != nil {
		{
			size, err := m.IsmGenesis.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MailboxMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MailboxMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MailboxMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.MessageId) > 0 {
		i -= len(m.MessageId)
		copy(dAtA[i:], m.MessageId)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.MessageId)))
		i--
		dAtA[i] = 0x12
	}
	if m.MailboxId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.MailboxId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.IsmGenesis != nil {
		l = m.IsmGenesis.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.PostDispatchGenesis != nil {
		l = m.PostDispatchGenesis.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Mailboxes) > 0 {
		for _, e := range m.Mailboxes {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Messages) > 0 {
		for _, e := range m.Messages {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.IsmSequence != 0 {
		n += 1 + sovGenesis(uint64(m.IsmSequence))
	}
	if m.PostDispatchSequence != 0 {
		n += 1 + sovGenesis(uint64(m.PostDispatchSequence))
	}
	if m.AppSequence != 0 {
		n += 1 + sovGenesis(uint64(m.AppSequence))
	}
	return n
}

func (m *MailboxMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MailboxId != 0 {
		n += 1 + sovGenesis(uint64(m.MailboxId))
	}
	l = len(m.MessageId)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsmGenesis", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.IsmGenesis == nil {
				m.IsmGenesis = &types.GenesisState{}
			}
			if err := m.IsmGenesis.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PostDispatchGenesis", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PostDispatchGenesis == nil {
				m.PostDispatchGenesis = &types1.GenesisState{}
			}
			if err := m.PostDispatchGenesis.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Mailboxes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Mailboxes = append(m.Mailboxes, &Mailbox{})
			if err := m.Mailboxes[len(m.Mailboxes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Messages", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Messages = append(m.Messages, &MailboxMessage{})
			if err := m.Messages[len(m.Messages)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsmSequence", wireType)
			}
			m.IsmSequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IsmSequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PostDispatchSequence", wireType)
			}
			m.PostDispatchSequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PostDispatchSequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AppSequence", wireType)
			}
			m.AppSequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AppSequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *MailboxMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: MailboxMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MailboxMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MailboxId", wireType)
			}
			m.MailboxId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MailboxId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageId", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessageId = append(m.MessageId[:0], dAtA[iNdEx:postIndex]...)
			if m.MessageId == nil {
				m.MessageId = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
