// Code generated by protoc-gen-go. DO NOT EDIT.
// source: vitepb/unconfirmed_meta_db.proto

package vitepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type UnconfirmedMeta struct {
	AccountId            []byte                `protobuf:"bytes,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
	TotalNumber          []byte                `protobuf:"bytes,2,opt,name=totalNumber,proto3" json:"totalNumber,omitempty"`
	UnconfirmedList      []*UnconfirmedByToken `protobuf:"bytes,3,rep,name=unconfirmedList,proto3" json:"unconfirmedList,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *UnconfirmedMeta) Reset()         { *m = UnconfirmedMeta{} }
func (m *UnconfirmedMeta) String() string { return proto.CompactTextString(m) }
func (*UnconfirmedMeta) ProtoMessage()    {}
func (*UnconfirmedMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_unconfirmed_meta_db_de6d3d4504e80e71, []int{0}
}
func (m *UnconfirmedMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnconfirmedMeta.Unmarshal(m, b)
}
func (m *UnconfirmedMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnconfirmedMeta.Marshal(b, m, deterministic)
}
func (dst *UnconfirmedMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnconfirmedMeta.Merge(dst, src)
}
func (m *UnconfirmedMeta) XXX_Size() int {
	return xxx_messageInfo_UnconfirmedMeta.Size(m)
}
func (m *UnconfirmedMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_UnconfirmedMeta.DiscardUnknown(m)
}

var xxx_messageInfo_UnconfirmedMeta proto.InternalMessageInfo

func (m *UnconfirmedMeta) GetAccountId() []byte {
	if m != nil {
		return m.AccountId
	}
	return nil
}

func (m *UnconfirmedMeta) GetTotalNumber() []byte {
	if m != nil {
		return m.TotalNumber
	}
	return nil
}

func (m *UnconfirmedMeta) GetUnconfirmedList() []*UnconfirmedByToken {
	if m != nil {
		return m.UnconfirmedList
	}
	return nil
}

type UnconfirmedByToken struct {
	TokenId              []byte   `protobuf:"bytes,1,opt,name=tokenId,proto3" json:"tokenId,omitempty"`
	TotalBalance         []byte   `protobuf:"bytes,2,opt,name=totalBalance,proto3" json:"totalBalance,omitempty"`
	HashList             [][]byte `protobuf:"bytes,3,rep,name=hashList,proto3" json:"hashList,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnconfirmedByToken) Reset()         { *m = UnconfirmedByToken{} }
func (m *UnconfirmedByToken) String() string { return proto.CompactTextString(m) }
func (*UnconfirmedByToken) ProtoMessage()    {}
func (*UnconfirmedByToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_unconfirmed_meta_db_de6d3d4504e80e71, []int{1}
}
func (m *UnconfirmedByToken) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnconfirmedByToken.Unmarshal(m, b)
}
func (m *UnconfirmedByToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnconfirmedByToken.Marshal(b, m, deterministic)
}
func (dst *UnconfirmedByToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnconfirmedByToken.Merge(dst, src)
}
func (m *UnconfirmedByToken) XXX_Size() int {
	return xxx_messageInfo_UnconfirmedByToken.Size(m)
}
func (m *UnconfirmedByToken) XXX_DiscardUnknown() {
	xxx_messageInfo_UnconfirmedByToken.DiscardUnknown(m)
}

var xxx_messageInfo_UnconfirmedByToken proto.InternalMessageInfo

func (m *UnconfirmedByToken) GetTokenId() []byte {
	if m != nil {
		return m.TokenId
	}
	return nil
}

func (m *UnconfirmedByToken) GetTotalBalance() []byte {
	if m != nil {
		return m.TotalBalance
	}
	return nil
}

func (m *UnconfirmedByToken) GetHashList() [][]byte {
	if m != nil {
		return m.HashList
	}
	return nil
}

func init() {
	proto.RegisterType((*UnconfirmedMeta)(nil), "vitepb.UnconfirmedMeta")
	proto.RegisterType((*UnconfirmedByToken)(nil), "vitepb.UnconfirmedByToken")
}

func init() {
	proto.RegisterFile("vitepb/unconfirmed_meta_db.proto", fileDescriptor_unconfirmed_meta_db_de6d3d4504e80e71)
}

var fileDescriptor_unconfirmed_meta_db_de6d3d4504e80e71 = []byte{
	// 205 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x28, 0xcb, 0x2c, 0x49,
	0x2d, 0x48, 0xd2, 0x2f, 0xcd, 0x4b, 0xce, 0xcf, 0x4b, 0xcb, 0x2c, 0xca, 0x4d, 0x4d, 0x89, 0xcf,
	0x4d, 0x2d, 0x49, 0x8c, 0x4f, 0x49, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0xa8,
	0x50, 0x9a, 0xce, 0xc8, 0xc5, 0x1f, 0x8a, 0x50, 0xe5, 0x9b, 0x5a, 0x92, 0x28, 0x24, 0xc3, 0xc5,
	0x99, 0x98, 0x9c, 0x9c, 0x5f, 0x9a, 0x57, 0xe2, 0x99, 0x22, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x13,
	0x84, 0x10, 0x10, 0x52, 0xe0, 0xe2, 0x2e, 0xc9, 0x2f, 0x49, 0xcc, 0xf1, 0x2b, 0xcd, 0x4d, 0x4a,
	0x2d, 0x92, 0x60, 0x02, 0xcb, 0x23, 0x0b, 0x09, 0xb9, 0x70, 0xf1, 0x23, 0x59, 0xec, 0x93, 0x59,
	0x5c, 0x22, 0xc1, 0xac, 0xc0, 0xac, 0xc1, 0x6d, 0x24, 0xa5, 0x07, 0xb1, 0x55, 0x0f, 0xc9, 0x46,
	0xa7, 0xca, 0x90, 0xfc, 0xec, 0xd4, 0xbc, 0x20, 0x74, 0x2d, 0x4a, 0x79, 0x5c, 0x42, 0x98, 0xca,
	0x84, 0x24, 0xb8, 0xd8, 0x4b, 0x40, 0x0c, 0xb8, 0xcb, 0x60, 0x5c, 0x21, 0x25, 0x2e, 0x1e, 0xb0,
	0x23, 0x9c, 0x12, 0x73, 0x12, 0xf3, 0x92, 0x53, 0xa1, 0x0e, 0x43, 0x11, 0x13, 0x92, 0xe2, 0xe2,
	0xc8, 0x48, 0x2c, 0xce, 0x80, 0x3b, 0x89, 0x27, 0x08, 0xce, 0x4f, 0x62, 0x03, 0x07, 0x8c, 0x31,
	0x20, 0x00, 0x00, 0xff, 0xff, 0xe9, 0x1f, 0xf5, 0x05, 0x3c, 0x01, 0x00, 0x00,
}
