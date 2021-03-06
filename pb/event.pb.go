// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.7.1
// source: event.proto

package pb

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type SubscribeType int32

const (
	// 区块事件，payload为BlockFilter
	SubscribeType_BLOCK SubscribeType = 0
)

// Enum value maps for SubscribeType.
var (
	SubscribeType_name = map[int32]string{
		0: "BLOCK",
	}
	SubscribeType_value = map[string]int32{
		"BLOCK": 0,
	}
)

func (x SubscribeType) Enum() *SubscribeType {
	p := new(SubscribeType)
	*p = x
	return p
}

func (x SubscribeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SubscribeType) Descriptor() protoreflect.EnumDescriptor {
	return file_event_proto_enumTypes[0].Descriptor()
}

func (SubscribeType) Type() protoreflect.EnumType {
	return &file_event_proto_enumTypes[0]
}

func (x SubscribeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SubscribeType.Descriptor instead.
func (SubscribeType) EnumDescriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{0}
}

type SubscribeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   SubscribeType `protobuf:"varint,1,opt,name=type,proto3,enum=pb.SubscribeType" json:"type,omitempty"`
	Filter []byte        `protobuf:"bytes,2,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *SubscribeRequest) Reset() {
	*x = SubscribeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeRequest) ProtoMessage() {}

func (x *SubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeRequest.ProtoReflect.Descriptor instead.
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{0}
}

func (x *SubscribeRequest) GetType() SubscribeType {
	if x != nil {
		return x.Type
	}
	return SubscribeType_BLOCK
}

func (x *SubscribeRequest) GetFilter() []byte {
	if x != nil {
		return x.Filter
	}
	return nil
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{1}
}

func (x *Event) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type BlockRange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start string `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End   string `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
}

func (x *BlockRange) Reset() {
	*x = BlockRange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockRange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockRange) ProtoMessage() {}

func (x *BlockRange) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockRange.ProtoReflect.Descriptor instead.
func (*BlockRange) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{2}
}

func (x *BlockRange) GetStart() string {
	if x != nil {
		return x.Start
	}
	return ""
}

func (x *BlockRange) GetEnd() string {
	if x != nil {
		return x.End
	}
	return ""
}

type BlockFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bcname         string      `protobuf:"bytes,1,opt,name=bcname,proto3" json:"bcname,omitempty"`
	Range          *BlockRange `protobuf:"bytes,2,opt,name=range,proto3" json:"range,omitempty"`
	ExcludeTx      bool        `protobuf:"varint,3,opt,name=exclude_tx,json=excludeTx,proto3" json:"exclude_tx,omitempty"`
	ExcludeTxEvent bool        `protobuf:"varint,4,opt,name=exclude_tx_event,json=excludeTxEvent,proto3" json:"exclude_tx_event,omitempty"`
	Contract       string      `protobuf:"bytes,10,opt,name=contract,proto3" json:"contract,omitempty"`
	EventName      string      `protobuf:"bytes,11,opt,name=event_name,json=eventName,proto3" json:"event_name,omitempty"`
	Initiator      string      `protobuf:"bytes,12,opt,name=initiator,proto3" json:"initiator,omitempty"`
	AuthRequire    string      `protobuf:"bytes,13,opt,name=auth_require,json=authRequire,proto3" json:"auth_require,omitempty"`
	FromAddr       string      `protobuf:"bytes,14,opt,name=from_addr,json=fromAddr,proto3" json:"from_addr,omitempty"`
	ToAddr         string      `protobuf:"bytes,15,opt,name=to_addr,json=toAddr,proto3" json:"to_addr,omitempty"`
}

func (x *BlockFilter) Reset() {
	*x = BlockFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockFilter) ProtoMessage() {}

func (x *BlockFilter) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockFilter.ProtoReflect.Descriptor instead.
func (*BlockFilter) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{3}
}

func (x *BlockFilter) GetBcname() string {
	if x != nil {
		return x.Bcname
	}
	return ""
}

func (x *BlockFilter) GetRange() *BlockRange {
	if x != nil {
		return x.Range
	}
	return nil
}

func (x *BlockFilter) GetExcludeTx() bool {
	if x != nil {
		return x.ExcludeTx
	}
	return false
}

func (x *BlockFilter) GetExcludeTxEvent() bool {
	if x != nil {
		return x.ExcludeTxEvent
	}
	return false
}

func (x *BlockFilter) GetContract() string {
	if x != nil {
		return x.Contract
	}
	return ""
}

func (x *BlockFilter) GetEventName() string {
	if x != nil {
		return x.EventName
	}
	return ""
}

func (x *BlockFilter) GetInitiator() string {
	if x != nil {
		return x.Initiator
	}
	return ""
}

func (x *BlockFilter) GetAuthRequire() string {
	if x != nil {
		return x.AuthRequire
	}
	return ""
}

func (x *BlockFilter) GetFromAddr() string {
	if x != nil {
		return x.FromAddr
	}
	return ""
}

func (x *BlockFilter) GetToAddr() string {
	if x != nil {
		return x.ToAddr
	}
	return ""
}

type FilteredTransaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txid   string           `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Events []*ContractEvent `protobuf:"bytes,2,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *FilteredTransaction) Reset() {
	*x = FilteredTransaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilteredTransaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilteredTransaction) ProtoMessage() {}

func (x *FilteredTransaction) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilteredTransaction.ProtoReflect.Descriptor instead.
func (*FilteredTransaction) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{4}
}

func (x *FilteredTransaction) GetTxid() string {
	if x != nil {
		return x.Txid
	}
	return ""
}

func (x *FilteredTransaction) GetEvents() []*ContractEvent {
	if x != nil {
		return x.Events
	}
	return nil
}

type FilteredBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bcname      string                 `protobuf:"bytes,1,opt,name=bcname,proto3" json:"bcname,omitempty"`
	Blockid     string                 `protobuf:"bytes,2,opt,name=blockid,proto3" json:"blockid,omitempty"`
	BlockHeight int64                  `protobuf:"varint,3,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	Txs         []*FilteredTransaction `protobuf:"bytes,4,rep,name=txs,proto3" json:"txs,omitempty"`
}

func (x *FilteredBlock) Reset() {
	*x = FilteredBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilteredBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilteredBlock) ProtoMessage() {}

func (x *FilteredBlock) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilteredBlock.ProtoReflect.Descriptor instead.
func (*FilteredBlock) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{5}
}

func (x *FilteredBlock) GetBcname() string {
	if x != nil {
		return x.Bcname
	}
	return ""
}

func (x *FilteredBlock) GetBlockid() string {
	if x != nil {
		return x.Blockid
	}
	return ""
}

func (x *FilteredBlock) GetBlockHeight() int64 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

func (x *FilteredBlock) GetTxs() []*FilteredTransaction {
	if x != nil {
		return x.Txs
	}
	return nil
}

type Logs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events []*Event `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *Logs) Reset() {
	*x = Logs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Logs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Logs) ProtoMessage() {}

func (x *Logs) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Logs.ProtoReflect.Descriptor instead.
func (*Logs) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{6}
}

func (x *Logs) GetEvents() []*Event {
	if x != nil {
		return x.Events
	}
	return nil
}

var File_event_proto protoreflect.FileDescriptor

var file_event_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70,
	0x62, 0x1a, 0x0c, 0x78, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x51, 0x0a, 0x10, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x22, 0x21, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x34, 0x0a, 0x0a, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x61,
	0x6e, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x22, 0xc6, 0x02, 0x0a, 0x0b,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x62,
	0x63, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x63, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x05, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x61, 0x6e,
	0x67, 0x65, 0x52, 0x05, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x63,
	0x6c, 0x75, 0x64, 0x65, 0x5f, 0x74, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x65,
	0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x54, 0x78, 0x12, 0x28, 0x0a, 0x10, 0x65, 0x78, 0x63, 0x6c,
	0x75, 0x64, 0x65, 0x5f, 0x74, 0x78, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0e, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x54, 0x78, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x61,
	0x75, 0x74, 0x68, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x61, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x12, 0x17, 0x0a, 0x07, 0x74,
	0x6f, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x6f,
	0x41, 0x64, 0x64, 0x72, 0x22, 0x54, 0x0a, 0x13, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x65, 0x64,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x78, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x78, 0x69, 0x64, 0x12,
	0x29, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x11, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x8f, 0x01, 0x0a, 0x0d, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06,
	0x62, 0x63, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x63,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x69, 0x64, 0x12, 0x21,
	0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x12, 0x29, 0x0a, 0x03, 0x74, 0x78, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x74, 0x78, 0x73, 0x22, 0x29, 0x0a, 0x04,
	0x4c, 0x6f, 0x67, 0x73, 0x12, 0x21, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52,
	0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2a, 0x1a, 0x0a, 0x0d, 0x53, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x42, 0x4c, 0x4f, 0x43,
	0x4b, 0x10, 0x00, 0x32, 0x3e, 0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x2e, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x12, 0x14, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x30, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_event_proto_rawDescOnce sync.Once
	file_event_proto_rawDescData = file_event_proto_rawDesc
)

func file_event_proto_rawDescGZIP() []byte {
	file_event_proto_rawDescOnce.Do(func() {
		file_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_event_proto_rawDescData)
	})
	return file_event_proto_rawDescData
}

var file_event_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_event_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_event_proto_goTypes = []interface{}{
	(SubscribeType)(0),          // 0: pb.SubscribeType
	(*SubscribeRequest)(nil),    // 1: pb.SubscribeRequest
	(*Event)(nil),               // 2: pb.Event
	(*BlockRange)(nil),          // 3: pb.BlockRange
	(*BlockFilter)(nil),         // 4: pb.BlockFilter
	(*FilteredTransaction)(nil), // 5: pb.FilteredTransaction
	(*FilteredBlock)(nil),       // 6: pb.FilteredBlock
	(*Logs)(nil),                // 7: pb.Logs
	(*ContractEvent)(nil),       // 8: pb.ContractEvent
}
var file_event_proto_depIdxs = []int32{
	0, // 0: pb.SubscribeRequest.type:type_name -> pb.SubscribeType
	3, // 1: pb.BlockFilter.range:type_name -> pb.BlockRange
	8, // 2: pb.FilteredTransaction.events:type_name -> pb.ContractEvent
	5, // 3: pb.FilteredBlock.txs:type_name -> pb.FilteredTransaction
	2, // 4: pb.Logs.events:type_name -> pb.Event
	1, // 5: pb.EventService.Subscribe:input_type -> pb.SubscribeRequest
	2, // 6: pb.EventService.Subscribe:output_type -> pb.Event
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_event_proto_init() }
func file_event_proto_init() {
	if File_event_proto != nil {
		return
	}
	file_xchain_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockRange); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockFilter); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilteredTransaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilteredBlock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Logs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_event_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_event_proto_goTypes,
		DependencyIndexes: file_event_proto_depIdxs,
		EnumInfos:         file_event_proto_enumTypes,
		MessageInfos:      file_event_proto_msgTypes,
	}.Build()
	File_event_proto = out.File
	file_event_proto_rawDesc = nil
	file_event_proto_goTypes = nil
	file_event_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EventServiceClient interface {
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (EventService_SubscribeClient, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
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

func (*UnimplementedEventServiceServer) Subscribe(*SubscribeRequest, EventService_SubscribeServer) error {
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
