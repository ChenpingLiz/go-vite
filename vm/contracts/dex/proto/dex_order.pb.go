// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dex_order.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Order struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Address              []byte   `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
	TradeToken           []byte   `protobuf:"bytes,3,opt,name=TradeToken,proto3" json:"TradeToken,omitempty"`
	QuoteToken           []byte   `protobuf:"bytes,4,opt,name=QuoteToken,proto3" json:"QuoteToken,omitempty"`
	Side                 bool     `protobuf:"varint,5,opt,name=Side,proto3" json:"Side,omitempty"`
	Type                 int32    `protobuf:"varint,6,opt,name=Type,proto3" json:"Type,omitempty"`
	Price                string   `protobuf:"bytes,7,opt,name=Price,proto3" json:"Price,omitempty"`
	Quantity             []byte   `protobuf:"bytes,8,opt,name=Quantity,proto3" json:"Quantity,omitempty"`
	Amount               []byte   `protobuf:"bytes,9,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Status               int32    `protobuf:"varint,10,opt,name=Status,proto3" json:"Status,omitempty"`
	CancelReason         int32    `protobuf:"varint,11,opt,name=cancelReason,proto3" json:"cancelReason,omitempty"`
	Timestamp            int64    `protobuf:"varint,12,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	ExecutedQuantity     []byte   `protobuf:"bytes,13,opt,name=ExecutedQuantity,proto3" json:"ExecutedQuantity,omitempty"`
	ExecutedAmount       []byte   `protobuf:"bytes,14,opt,name=ExecutedAmount,proto3" json:"ExecutedAmount,omitempty"`
	RefundToken          []byte   `protobuf:"bytes,15,opt,name=refundToken,proto3" json:"refundToken,omitempty"`
	RefundQuantity       []byte   `protobuf:"bytes,16,opt,name=refundQuantity,proto3" json:"refundQuantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()    {}
func (*Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_5df90bf1b4327f51, []int{0}
}

func (m *Order) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Order.Unmarshal(m, b)
}
func (m *Order) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Order.Marshal(b, m, deterministic)
}
func (m *Order) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Order.Merge(m, src)
}
func (m *Order) XXX_Size() int {
	return xxx_messageInfo_Order.Size(m)
}
func (m *Order) XXX_DiscardUnknown() {
	xxx_messageInfo_Order.DiscardUnknown(m)
}

var xxx_messageInfo_Order proto.InternalMessageInfo

func (m *Order) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Order) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *Order) GetTradeToken() []byte {
	if m != nil {
		return m.TradeToken
	}
	return nil
}

func (m *Order) GetQuoteToken() []byte {
	if m != nil {
		return m.QuoteToken
	}
	return nil
}

func (m *Order) GetSide() bool {
	if m != nil {
		return m.Side
	}
	return false
}

func (m *Order) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Order) GetPrice() string {
	if m != nil {
		return m.Price
	}
	return ""
}

func (m *Order) GetQuantity() []byte {
	if m != nil {
		return m.Quantity
	}
	return nil
}

func (m *Order) GetAmount() []byte {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *Order) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *Order) GetCancelReason() int32 {
	if m != nil {
		return m.CancelReason
	}
	return 0
}

func (m *Order) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Order) GetExecutedQuantity() []byte {
	if m != nil {
		return m.ExecutedQuantity
	}
	return nil
}

func (m *Order) GetExecutedAmount() []byte {
	if m != nil {
		return m.ExecutedAmount
	}
	return nil
}

func (m *Order) GetRefundToken() []byte {
	if m != nil {
		return m.RefundToken
	}
	return nil
}

func (m *Order) GetRefundQuantity() []byte {
	if m != nil {
		return m.RefundQuantity
	}
	return nil
}

type OrderNode struct {
	Order                *Order   `protobuf:"bytes,1,opt,name=Order,proto3" json:"Order,omitempty"`
	ForwardOnLevel       [][]byte `protobuf:"bytes,2,rep,name=ForwardOnLevel,proto3" json:"ForwardOnLevel,omitempty"`
	BackwardOnLevel      [][]byte `protobuf:"bytes,3,rep,name=backwardOnLevel,proto3" json:"backwardOnLevel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderNode) Reset()         { *m = OrderNode{} }
func (m *OrderNode) String() string { return proto.CompactTextString(m) }
func (*OrderNode) ProtoMessage()    {}
func (*OrderNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_5df90bf1b4327f51, []int{1}
}

func (m *OrderNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderNode.Unmarshal(m, b)
}
func (m *OrderNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderNode.Marshal(b, m, deterministic)
}
func (m *OrderNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderNode.Merge(m, src)
}
func (m *OrderNode) XXX_Size() int {
	return xxx_messageInfo_OrderNode.Size(m)
}
func (m *OrderNode) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderNode.DiscardUnknown(m)
}

var xxx_messageInfo_OrderNode proto.InternalMessageInfo

func (m *OrderNode) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *OrderNode) GetForwardOnLevel() [][]byte {
	if m != nil {
		return m.ForwardOnLevel
	}
	return nil
}

func (m *OrderNode) GetBackwardOnLevel() [][]byte {
	if m != nil {
		return m.BackwardOnLevel
	}
	return nil
}

type OrderListMeta struct {
	Header               []byte   `protobuf:"bytes,1,opt,name=Header,proto3" json:"Header,omitempty"`
	Tail                 []byte   `protobuf:"bytes,2,opt,name=Tail,proto3" json:"Tail,omitempty"`
	Length               int32    `protobuf:"varint,3,opt,name=Length,proto3" json:"Length,omitempty"`
	Level                int32    `protobuf:"varint,4,opt,name=Level,proto3" json:"Level,omitempty"`
	ForwardOnLevel       [][]byte `protobuf:"bytes,5,rep,name=ForwardOnLevel,proto3" json:"ForwardOnLevel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderListMeta) Reset()         { *m = OrderListMeta{} }
func (m *OrderListMeta) String() string { return proto.CompactTextString(m) }
func (*OrderListMeta) ProtoMessage()    {}
func (*OrderListMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_5df90bf1b4327f51, []int{2}
}

func (m *OrderListMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderListMeta.Unmarshal(m, b)
}
func (m *OrderListMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderListMeta.Marshal(b, m, deterministic)
}
func (m *OrderListMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderListMeta.Merge(m, src)
}
func (m *OrderListMeta) XXX_Size() int {
	return xxx_messageInfo_OrderListMeta.Size(m)
}
func (m *OrderListMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderListMeta.DiscardUnknown(m)
}

var xxx_messageInfo_OrderListMeta proto.InternalMessageInfo

func (m *OrderListMeta) GetHeader() []byte {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *OrderListMeta) GetTail() []byte {
	if m != nil {
		return m.Tail
	}
	return nil
}

func (m *OrderListMeta) GetLength() int32 {
	if m != nil {
		return m.Length
	}
	return 0
}

func (m *OrderListMeta) GetLevel() int32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *OrderListMeta) GetForwardOnLevel() [][]byte {
	if m != nil {
		return m.ForwardOnLevel
	}
	return nil
}

type OrderBook struct {
	MarketSide           []string `protobuf:"bytes,1,rep,name=MarketSide,proto3" json:"MarketSide,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderBook) Reset()         { *m = OrderBook{} }
func (m *OrderBook) String() string { return proto.CompactTextString(m) }
func (*OrderBook) ProtoMessage()    {}
func (*OrderBook) Descriptor() ([]byte, []int) {
	return fileDescriptor_5df90bf1b4327f51, []int{3}
}

func (m *OrderBook) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderBook.Unmarshal(m, b)
}
func (m *OrderBook) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderBook.Marshal(b, m, deterministic)
}
func (m *OrderBook) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderBook.Merge(m, src)
}
func (m *OrderBook) XXX_Size() int {
	return xxx_messageInfo_OrderBook.Size(m)
}
func (m *OrderBook) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderBook.DiscardUnknown(m)
}

var xxx_messageInfo_OrderBook proto.InternalMessageInfo

func (m *OrderBook) GetMarketSide() []string {
	if m != nil {
		return m.MarketSide
	}
	return nil
}

type Transaction struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	TakerSide            bool     `protobuf:"varint,2,opt,name=TakerSide,proto3" json:"TakerSide,omitempty"`
	TakerId              []byte   `protobuf:"bytes,3,opt,name=TakerId,proto3" json:"TakerId,omitempty"`
	MakerId              []byte   `protobuf:"bytes,4,opt,name=MakerId,proto3" json:"MakerId,omitempty"`
	Price                string   `protobuf:"bytes,5,opt,name=Price,proto3" json:"Price,omitempty"`
	Quantity             []byte   `protobuf:"bytes,6,opt,name=Quantity,proto3" json:"Quantity,omitempty"`
	Amount               []byte   `protobuf:"bytes,7,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Timestamp            int64    `protobuf:"varint,8,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_5df90bf1b4327f51, []int{4}
}

func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (m *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(m, src)
}
func (m *Transaction) XXX_Size() int {
	return xxx_messageInfo_Transaction.Size(m)
}
func (m *Transaction) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

func (m *Transaction) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Transaction) GetTakerSide() bool {
	if m != nil {
		return m.TakerSide
	}
	return false
}

func (m *Transaction) GetTakerId() []byte {
	if m != nil {
		return m.TakerId
	}
	return nil
}

func (m *Transaction) GetMakerId() []byte {
	if m != nil {
		return m.MakerId
	}
	return nil
}

func (m *Transaction) GetPrice() string {
	if m != nil {
		return m.Price
	}
	return ""
}

func (m *Transaction) GetQuantity() []byte {
	if m != nil {
		return m.Quantity
	}
	return nil
}

func (m *Transaction) GetAmount() []byte {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *Transaction) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func init() {
	proto.RegisterType((*Order)(nil), "proto.Order")
	proto.RegisterType((*OrderNode)(nil), "proto.OrderNode")
	proto.RegisterType((*OrderListMeta)(nil), "proto.OrderListMeta")
	proto.RegisterType((*OrderBook)(nil), "proto.OrderBook")
	proto.RegisterType((*Transaction)(nil), "proto.Transaction")
}

func init() { proto.RegisterFile("dex_order.proto", fileDescriptor_5df90bf1b4327f51) }

var fileDescriptor_5df90bf1b4327f51 = []byte{
	// 504 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x93, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0xc7, 0x95, 0x66, 0xe9, 0xc7, 0x69, 0xd7, 0x4e, 0x16, 0x42, 0x16, 0x42, 0x28, 0xca, 0xc5,
	0x14, 0x81, 0xb4, 0x0b, 0x78, 0x82, 0x21, 0x81, 0xa8, 0xd4, 0x32, 0xe6, 0xf5, 0x1e, 0x79, 0xf1,
	0x01, 0xa2, 0xb6, 0x76, 0xe5, 0x38, 0xb0, 0x5d, 0xf2, 0x08, 0x3c, 0x1d, 0x17, 0xbc, 0x0c, 0xf2,
	0x89, 0xfb, 0x91, 0x6e, 0xbb, 0xaa, 0xff, 0xbf, 0xe3, 0x1e, 0x1f, 0xc5, 0x3f, 0xc3, 0x44, 0xe1,
	0xdd, 0x57, 0x63, 0x15, 0xda, 0x8b, 0x8d, 0x35, 0xce, 0xb0, 0x84, 0x7e, 0xb2, 0x7f, 0x31, 0x24,
	0x57, 0x1e, 0xb3, 0x31, 0x74, 0xa6, 0x8a, 0x47, 0x69, 0x94, 0x8f, 0x44, 0x67, 0xaa, 0x18, 0x87,
	0xde, 0xa5, 0x52, 0x16, 0xab, 0x8a, 0x77, 0x08, 0x6e, 0x23, 0x7b, 0x05, 0xb0, 0xb0, 0x52, 0xe1,
	0xc2, 0x2c, 0x51, 0xf3, 0x98, 0x8a, 0x07, 0xc4, 0xd7, 0xaf, 0x6b, 0xe3, 0x42, 0xfd, 0xa4, 0xa9,
	0xef, 0x09, 0x63, 0x70, 0x72, 0x53, 0x2a, 0xe4, 0x49, 0x1a, 0xe5, 0x7d, 0x41, 0x6b, 0xcf, 0x16,
	0xf7, 0x1b, 0xe4, 0xdd, 0x34, 0xca, 0x13, 0x41, 0x6b, 0xf6, 0x0c, 0x92, 0x2f, 0xb6, 0x2c, 0x90,
	0xf7, 0xd2, 0x28, 0x1f, 0x88, 0x26, 0xb0, 0x17, 0xd0, 0xbf, 0xae, 0xa5, 0x76, 0xa5, 0xbb, 0xe7,
	0x7d, 0xea, 0xbd, 0xcb, 0xec, 0x39, 0x74, 0x2f, 0xd7, 0xa6, 0xd6, 0x8e, 0x0f, 0xa8, 0x12, 0x92,
	0xe7, 0x37, 0x4e, 0xba, 0xba, 0xe2, 0x40, 0xfd, 0x43, 0x62, 0x19, 0x8c, 0x0a, 0xa9, 0x0b, 0x5c,
	0x09, 0x94, 0x95, 0xd1, 0x7c, 0x48, 0xd5, 0x16, 0x63, 0x2f, 0x61, 0xb0, 0x28, 0xd7, 0x58, 0x39,
	0xb9, 0xde, 0xf0, 0x51, 0x1a, 0xe5, 0xb1, 0xd8, 0x03, 0xf6, 0x1a, 0xce, 0x3e, 0xdc, 0x61, 0x51,
	0x3b, 0x54, 0xbb, 0xa9, 0x4e, 0xe9, 0xec, 0x07, 0x9c, 0x9d, 0xc3, 0x78, 0xcb, 0xc2, 0x94, 0x63,
	0xda, 0x79, 0x44, 0x59, 0x0a, 0x43, 0x8b, 0xdf, 0x6a, 0xad, 0x9a, 0x0f, 0x38, 0xa1, 0x4d, 0x87,
	0xc8, 0x77, 0x6a, 0xe2, 0xee, 0xcc, 0xb3, 0xa6, 0x53, 0x9b, 0x66, 0xbf, 0x23, 0x18, 0xd0, 0xed,
	0x7e, 0x36, 0x0a, 0x59, 0x16, 0xae, 0x9a, 0x2e, 0x79, 0xf8, 0x76, 0xd4, 0x98, 0x70, 0x41, 0x4c,
	0x04, 0x0b, 0xce, 0x61, 0xfc, 0xd1, 0xd8, 0x5f, 0xd2, 0xaa, 0x2b, 0x3d, 0xc3, 0x9f, 0xb8, 0xe2,
	0x9d, 0x34, 0xf6, 0x9d, 0xdb, 0x94, 0xe5, 0x30, 0xb9, 0x95, 0xc5, 0xf2, 0x70, 0x63, 0x4c, 0x1b,
	0x8f, 0x71, 0xf6, 0x27, 0x82, 0x53, 0xea, 0x3d, 0x2b, 0x2b, 0x37, 0x47, 0x27, 0xfd, 0x6d, 0x7c,
	0x42, 0xb9, 0x1d, 0x64, 0x24, 0x42, 0x22, 0x07, 0x64, 0xb9, 0x0a, 0xba, 0xd1, 0xda, 0xef, 0x9d,
	0xa1, 0xfe, 0xee, 0x7e, 0x90, 0x67, 0x89, 0x08, 0xc9, 0xbb, 0xd1, 0x9c, 0x7a, 0x42, 0xb8, 0x09,
	0x8f, 0x4c, 0x9f, 0x3c, 0x36, 0x7d, 0xf6, 0x26, 0x7c, 0x96, 0xf7, 0xc6, 0x2c, 0xbd, 0xae, 0x73,
	0x69, 0x97, 0xe8, 0x48, 0xca, 0x28, 0x8d, 0xf3, 0x81, 0x38, 0x20, 0xd9, 0xdf, 0x08, 0x86, 0x0b,
	0x2b, 0x75, 0x25, 0x0b, 0x57, 0x1a, 0xfd, 0xe0, 0xa1, 0x78, 0x41, 0xe4, 0x12, 0x2d, 0xfd, 0xbd,
	0x43, 0x4e, 0xef, 0x81, 0x7f, 0x46, 0x14, 0xa6, 0x2a, 0xbc, 0x94, 0x6d, 0xf4, 0x95, 0x79, 0xa8,
	0x34, 0x6f, 0x64, 0x1b, 0xf7, 0xe2, 0x27, 0x4f, 0x89, 0xdf, 0x7d, 0x52, 0xfc, 0x5e, 0x4b, 0xfc,
	0x96, 0xbc, 0xfd, 0x23, 0x79, 0x6f, 0xbb, 0x24, 0xc0, 0xbb, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x5b, 0x60, 0xa4, 0xee, 0x1d, 0x04, 0x00, 0x00,
}