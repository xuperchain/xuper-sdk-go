// Code generated by protoc-gen-go. DO NOT EDIT.
// source: event.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SubscribeType int32

const (
	// 区块事件，payload为BlockFilter
	SubscribeType_BLOCK SubscribeType = 0
)

var SubscribeType_name = map[int32]string{
	0: "BLOCK",
}

var SubscribeType_value = map[string]int32{
	"BLOCK": 0,
}

func (x SubscribeType) String() string {
	return proto.EnumName(SubscribeType_name, int32(x))
}

func (SubscribeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{0}
}

type SubscribeRequest struct {
	Type                 SubscribeType `protobuf:"varint,1,opt,name=type,proto3,enum=pb.SubscribeType" json:"type,omitempty"`
	Filter               []byte        `protobuf:"bytes,2,opt,name=filter,proto3" json:"filter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *SubscribeRequest) Reset()         { *m = SubscribeRequest{} }
func (m *SubscribeRequest) String() string { return proto.CompactTextString(m) }
func (*SubscribeRequest) ProtoMessage()    {}
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{0}
}

func (m *SubscribeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscribeRequest.Unmarshal(m, b)
}
func (m *SubscribeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscribeRequest.Marshal(b, m, deterministic)
}
func (m *SubscribeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscribeRequest.Merge(m, src)
}
func (m *SubscribeRequest) XXX_Size() int {
	return xxx_messageInfo_SubscribeRequest.Size(m)
}
func (m *SubscribeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscribeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubscribeRequest proto.InternalMessageInfo

func (m *SubscribeRequest) GetType() SubscribeType {
	if m != nil {
		return m.Type
	}
	return SubscribeType_BLOCK
}

func (m *SubscribeRequest) GetFilter() []byte {
	if m != nil {
		return m.Filter
	}
	return nil
}

type Event struct {
	Payload              []byte   `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{1}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type BlockRange struct {
	Start                string   `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End                  string   `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockRange) Reset()         { *m = BlockRange{} }
func (m *BlockRange) String() string { return proto.CompactTextString(m) }
func (*BlockRange) ProtoMessage()    {}
func (*BlockRange) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{2}
}

func (m *BlockRange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockRange.Unmarshal(m, b)
}
func (m *BlockRange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockRange.Marshal(b, m, deterministic)
}
func (m *BlockRange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockRange.Merge(m, src)
}
func (m *BlockRange) XXX_Size() int {
	return xxx_messageInfo_BlockRange.Size(m)
}
func (m *BlockRange) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockRange.DiscardUnknown(m)
}

var xxx_messageInfo_BlockRange proto.InternalMessageInfo

func (m *BlockRange) GetStart() string {
	if m != nil {
		return m.Start
	}
	return ""
}

func (m *BlockRange) GetEnd() string {
	if m != nil {
		return m.End
	}
	return ""
}

type BlockFilter struct {
	Bcname               string      `protobuf:"bytes,1,opt,name=bcname,proto3" json:"bcname,omitempty"`
	Range                *BlockRange `protobuf:"bytes,2,opt,name=range,proto3" json:"range,omitempty"`
	ExcludeTx            bool        `protobuf:"varint,3,opt,name=exclude_tx,json=excludeTx,proto3" json:"exclude_tx,omitempty"`
	ExcludeTxEvent       bool        `protobuf:"varint,4,opt,name=exclude_tx_event,json=excludeTxEvent,proto3" json:"exclude_tx_event,omitempty"`
	Contract             string      `protobuf:"bytes,10,opt,name=contract,proto3" json:"contract,omitempty"`
	EventName            string      `protobuf:"bytes,11,opt,name=event_name,json=eventName,proto3" json:"event_name,omitempty"`
	Initiator            string      `protobuf:"bytes,12,opt,name=initiator,proto3" json:"initiator,omitempty"`
	AuthRequire          string      `protobuf:"bytes,13,opt,name=auth_require,json=authRequire,proto3" json:"auth_require,omitempty"`
	FromAddr             string      `protobuf:"bytes,14,opt,name=from_addr,json=fromAddr,proto3" json:"from_addr,omitempty"`
	ToAddr               string      `protobuf:"bytes,15,opt,name=to_addr,json=toAddr,proto3" json:"to_addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *BlockFilter) Reset()         { *m = BlockFilter{} }
func (m *BlockFilter) String() string { return proto.CompactTextString(m) }
func (*BlockFilter) ProtoMessage()    {}
func (*BlockFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{3}
}

func (m *BlockFilter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockFilter.Unmarshal(m, b)
}
func (m *BlockFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockFilter.Marshal(b, m, deterministic)
}
func (m *BlockFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockFilter.Merge(m, src)
}
func (m *BlockFilter) XXX_Size() int {
	return xxx_messageInfo_BlockFilter.Size(m)
}
func (m *BlockFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockFilter.DiscardUnknown(m)
}

var xxx_messageInfo_BlockFilter proto.InternalMessageInfo

func (m *BlockFilter) GetBcname() string {
	if m != nil {
		return m.Bcname
	}
	return ""
}

func (m *BlockFilter) GetRange() *BlockRange {
	if m != nil {
		return m.Range
	}
	return nil
}

func (m *BlockFilter) GetExcludeTx() bool {
	if m != nil {
		return m.ExcludeTx
	}
	return false
}

func (m *BlockFilter) GetExcludeTxEvent() bool {
	if m != nil {
		return m.ExcludeTxEvent
	}
	return false
}

func (m *BlockFilter) GetContract() string {
	if m != nil {
		return m.Contract
	}
	return ""
}

func (m *BlockFilter) GetEventName() string {
	if m != nil {
		return m.EventName
	}
	return ""
}

func (m *BlockFilter) GetInitiator() string {
	if m != nil {
		return m.Initiator
	}
	return ""
}

func (m *BlockFilter) GetAuthRequire() string {
	if m != nil {
		return m.AuthRequire
	}
	return ""
}

func (m *BlockFilter) GetFromAddr() string {
	if m != nil {
		return m.FromAddr
	}
	return ""
}

func (m *BlockFilter) GetToAddr() string {
	if m != nil {
		return m.ToAddr
	}
	return ""
}

type ContractEvent struct {
	Contract string `json:"contract,omitempty"`
	Name     string `json:"name,omitempty"`
	Body     string `json:"body,omitempty"`
}

type FilteredTransaction struct {
	Txid                 string           `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Events               []*ContractEvent `protobuf:"bytes,2,rep,name=events,proto3" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *FilteredTransaction) Reset()         { *m = FilteredTransaction{} }
func (m *FilteredTransaction) String() string { return proto.CompactTextString(m) }
func (*FilteredTransaction) ProtoMessage()    {}
func (*FilteredTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{4}
}

func (m *FilteredTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilteredTransaction.Unmarshal(m, b)
}
func (m *FilteredTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilteredTransaction.Marshal(b, m, deterministic)
}
func (m *FilteredTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilteredTransaction.Merge(m, src)
}
func (m *FilteredTransaction) XXX_Size() int {
	return xxx_messageInfo_FilteredTransaction.Size(m)
}
func (m *FilteredTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_FilteredTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_FilteredTransaction proto.InternalMessageInfo

func (m *FilteredTransaction) GetTxid() string {
	if m != nil {
		return m.Txid
	}
	return ""
}

func (m *FilteredTransaction) GetEvents() []*ContractEvent {
	if m != nil {
		return m.Events
	}
	return nil
}

type FilteredBlock struct {
	Bcname               string                 `protobuf:"bytes,1,opt,name=bcname,proto3" json:"bcname,omitempty"`
	Blockid              string                 `protobuf:"bytes,2,opt,name=blockid,proto3" json:"blockid,omitempty"`
	BlockHeight          int64                  `protobuf:"varint,3,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	Txs                  []*FilteredTransaction `protobuf:"bytes,4,rep,name=txs,proto3" json:"txs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *FilteredBlock) Reset()         { *m = FilteredBlock{} }
func (m *FilteredBlock) String() string { return proto.CompactTextString(m) }
func (*FilteredBlock) ProtoMessage()    {}
func (*FilteredBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_2d17a9d3f0ddf27e, []int{5}
}

func (m *FilteredBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilteredBlock.Unmarshal(m, b)
}
func (m *FilteredBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilteredBlock.Marshal(b, m, deterministic)
}
func (m *FilteredBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilteredBlock.Merge(m, src)
}
func (m *FilteredBlock) XXX_Size() int {
	return xxx_messageInfo_FilteredBlock.Size(m)
}
func (m *FilteredBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_FilteredBlock.DiscardUnknown(m)
}

var xxx_messageInfo_FilteredBlock proto.InternalMessageInfo

func (m *FilteredBlock) GetBcname() string {
	if m != nil {
		return m.Bcname
	}
	return ""
}

func (m *FilteredBlock) GetBlockid() string {
	if m != nil {
		return m.Blockid
	}
	return ""
}

func (m *FilteredBlock) GetBlockHeight() int64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *FilteredBlock) GetTxs() []*FilteredTransaction {
	if m != nil {
		return m.Txs
	}
	return nil
}

func init() {
	proto.RegisterEnum("pb.SubscribeType", SubscribeType_name, SubscribeType_value)
	proto.RegisterType((*SubscribeRequest)(nil), "pb.SubscribeRequest")
	proto.RegisterType((*Event)(nil), "pb.Event")
	proto.RegisterType((*BlockRange)(nil), "pb.BlockRange")
	proto.RegisterType((*BlockFilter)(nil), "pb.BlockFilter")
	proto.RegisterType((*FilteredTransaction)(nil), "pb.FilteredTransaction")
	proto.RegisterType((*FilteredBlock)(nil), "pb.FilteredBlock")
}

func init() { proto.RegisterFile("event.proto", fileDescriptor_2d17a9d3f0ddf27e) }

var fileDescriptor_2d17a9d3f0ddf27e = []byte{
	// 502 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x53, 0xdf, 0x6f, 0xd3, 0x30,
	0x10, 0x26, 0x4d, 0x7f, 0x2c, 0x97, 0xb4, 0x14, 0x33, 0x31, 0xab, 0x80, 0xd4, 0x45, 0x20, 0x65,
	0x3c, 0x54, 0xa8, 0xf0, 0x8c, 0xc4, 0x26, 0x10, 0x12, 0x08, 0x84, 0xd7, 0xf7, 0xc8, 0x49, 0xbc,
	0xd5, 0xa2, 0x8b, 0x33, 0xc7, 0x99, 0xd2, 0xbf, 0x82, 0xff, 0x8a, 0xbf, 0x0b, 0xf9, 0x9c, 0x66,
	0x9a, 0x04, 0x6f, 0xbe, 0xef, 0xfb, 0xee, 0xee, 0xbb, 0xb3, 0x0d, 0xa1, 0xb8, 0x13, 0xa5, 0x59,
	0x55, 0x5a, 0x19, 0x45, 0x06, 0x55, 0xb6, 0x88, 0xda, 0x7c, 0xcb, 0x65, 0xe9, 0x90, 0xf8, 0x27,
	0xcc, 0x2f, 0x9b, 0xac, 0xce, 0xb5, 0xcc, 0x04, 0x13, 0xb7, 0x8d, 0xa8, 0x0d, 0x79, 0x0d, 0x43,
	0xb3, 0xaf, 0x04, 0xf5, 0x96, 0x5e, 0x32, 0x5b, 0x3f, 0x59, 0x55, 0xd9, 0xaa, 0xd7, 0x6c, 0xf6,
	0x95, 0x60, 0x48, 0x93, 0x67, 0x30, 0xbe, 0x92, 0x3b, 0x23, 0x34, 0x1d, 0x2c, 0xbd, 0x24, 0x62,
	0x5d, 0x14, 0x9f, 0xc2, 0xe8, 0x93, 0xed, 0x49, 0x28, 0x4c, 0x2a, 0xbe, 0xdf, 0x29, 0x5e, 0x60,
	0xa9, 0x88, 0x1d, 0xc2, 0xf8, 0x3d, 0xc0, 0xf9, 0x4e, 0xe5, 0xbf, 0x18, 0x2f, 0xaf, 0x05, 0x39,
	0x86, 0x51, 0x6d, 0xb8, 0x36, 0xa8, 0x0a, 0x98, 0x0b, 0xc8, 0x1c, 0x7c, 0x51, 0x16, 0x58, 0x3b,
	0x60, 0xf6, 0x18, 0xff, 0x19, 0x40, 0x88, 0x69, 0x9f, 0xb1, 0x91, 0x35, 0x90, 0xe5, 0x25, 0xbf,
	0x11, 0x5d, 0x62, 0x17, 0x91, 0x57, 0x30, 0xd2, 0xb6, 0x30, 0xe6, 0x86, 0xeb, 0x99, 0x1d, 0xe0,
	0xbe, 0x1d, 0x73, 0x24, 0x79, 0x09, 0x20, 0xda, 0x7c, 0xd7, 0x14, 0x22, 0x35, 0x2d, 0xf5, 0x97,
	0x5e, 0x72, 0xc4, 0x82, 0x0e, 0xd9, 0xb4, 0x24, 0x81, 0xf9, 0x3d, 0x9d, 0xe2, 0x12, 0xe9, 0x10,
	0x45, 0xb3, 0x5e, 0xe4, 0xc6, 0x5c, 0xc0, 0x51, 0xae, 0x4a, 0xa3, 0x79, 0x6e, 0x28, 0xa0, 0x91,
	0x3e, 0xc6, 0x26, 0x56, 0x94, 0xa2, 0xcd, 0x10, 0xd9, 0x00, 0x91, 0xef, 0xd6, 0xe9, 0x0b, 0x08,
	0x64, 0x29, 0x8d, 0xe4, 0x46, 0x69, 0x1a, 0x39, 0xb6, 0x07, 0xc8, 0x29, 0x44, 0xbc, 0x31, 0xdb,
	0x54, 0x8b, 0xdb, 0x46, 0x6a, 0x41, 0xa7, 0x28, 0x08, 0x2d, 0xc6, 0x1c, 0x44, 0x9e, 0x43, 0x70,
	0xa5, 0xd5, 0x4d, 0xca, 0x8b, 0x42, 0xd3, 0x99, 0x6b, 0x6e, 0x81, 0x8f, 0x45, 0xa1, 0xc9, 0x09,
	0x4c, 0x8c, 0x72, 0xd4, 0x63, 0xb7, 0x20, 0xa3, 0x2c, 0x11, 0x6f, 0xe0, 0xa9, 0x5b, 0xa1, 0x28,
	0x36, 0x9a, 0x97, 0x35, 0xcf, 0x8d, 0x54, 0x25, 0x21, 0x30, 0x34, 0xad, 0x2c, 0xba, 0x6d, 0xe2,
	0x99, 0x9c, 0xc1, 0x18, 0xed, 0xd6, 0x74, 0xb0, 0xf4, 0x93, 0xd0, 0xbd, 0x86, 0x8b, 0x6e, 0x3c,
	0x9c, 0x9f, 0x75, 0x82, 0xf8, 0xb7, 0x07, 0xd3, 0x43, 0x59, 0x5c, 0xf7, 0x7f, 0x2f, 0x88, 0xc2,
	0x24, 0xb3, 0x02, 0x79, 0xb8, 0xde, 0x43, 0x68, 0x47, 0xc6, 0x63, 0xba, 0x15, 0xf2, 0x7a, 0x6b,
	0xf0, 0x5a, 0x7c, 0x16, 0x22, 0xf6, 0x05, 0x21, 0x72, 0x06, 0xbe, 0x69, 0x6b, 0x3a, 0x44, 0x3b,
	0x27, 0xd6, 0xce, 0x3f, 0x66, 0x61, 0x56, 0xf3, 0x66, 0x01, 0xd3, 0x07, 0x0f, 0x97, 0x04, 0x30,
	0x3a, 0xff, 0xf6, 0xe3, 0xe2, 0xeb, 0xfc, 0xd1, 0xfa, 0x03, 0x44, 0x68, 0xff, 0x52, 0xe8, 0x3b,
	0x99, 0x0b, 0xb2, 0x82, 0xa0, 0xd7, 0x92, 0xe3, 0x07, 0x6f, 0xbe, 0xfb, 0x17, 0x8b, 0xc0, 0xa2,
	0x98, 0xf4, 0xd6, 0xcb, 0xc6, 0xf8, 0x7f, 0xde, 0xfd, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x74, 0x59,
	0xa4, 0xba, 0x60, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EventServiceClient interface {
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (EventService_SubscribeClient, error)
}

type eventServiceClient struct {
	cc *grpc.ClientConn
}

func NewEventServiceClient(cc *grpc.ClientConn) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (EventService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_EventService_serviceDesc.Streams[0], "/pb.EventService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EventService_SubscribeClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *eventServiceSubscribeClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventServiceServer is the server API for EventService service.
type EventServiceServer interface {
	Subscribe(*SubscribeRequest, EventService_SubscribeServer) error
}

// UnimplementedEventServiceServer can be embedded to have forward compatible implementations.
type UnimplementedEventServiceServer struct {
}

func (*UnimplementedEventServiceServer) Subscribe(req *SubscribeRequest, srv EventService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}

func RegisterEventServiceServer(s *grpc.Server, srv EventServiceServer) {
	s.RegisterService(&_EventService_serviceDesc, srv)
}

func _EventService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EventServiceServer).Subscribe(m, &eventServiceSubscribeServer{stream})
}

type EventService_SubscribeServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type eventServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *eventServiceSubscribeServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

var _EventService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _EventService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "event.proto",
}
