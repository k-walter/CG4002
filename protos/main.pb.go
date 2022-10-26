// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: protos/main.proto

package __

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Action int32

const (
	// for eval
	Action_none    Action = 0
	Action_grenade Action = 1
	Action_reload  Action = 2
	Action_shoot   Action = 3
	Action_logout  Action = 4
	Action_shield  Action = 5
	// internal use
	Action_shot            Action = 6
	Action_grenaded        Action = 7
	Action_shieldAvailable Action = 8
	Action_checkFov        Action = 9 // for viz to check only, grenade to throw (after cfm with eval), grenaded to be hit
)

// Enum value maps for Action.
var (
	Action_name = map[int32]string{
		0: "none",
		1: "grenade",
		2: "reload",
		3: "shoot",
		4: "logout",
		5: "shield",
		6: "shot",
		7: "grenaded",
		8: "shieldAvailable",
		9: "checkFov",
	}
	Action_value = map[string]int32{
		"none":            0,
		"grenade":         1,
		"reload":          2,
		"shoot":           3,
		"logout":          4,
		"shield":          5,
		"shot":            6,
		"grenaded":        7,
		"shieldAvailable": 8,
		"checkFov":        9,
	}
)

func (x Action) Enum() *Action {
	p := new(Action)
	*p = x
	return p
}

func (x Action) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Action) Descriptor() protoreflect.EnumDescriptor {
	return file_protos_main_proto_enumTypes[0].Descriptor()
}

func (Action) Type() protoreflect.EnumType {
	return &file_protos_main_proto_enumTypes[0]
}

func (x Action) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Action.Descriptor instead.
func (Action) EnumDescriptor() ([]byte, []int) {
	return file_protos_main_proto_rawDescGZIP(), []int{0}
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Player uint32 `protobuf:"varint,1,opt,name=player,proto3" json:"player,omitempty"`
	Time   uint64 `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"` // No need to fill, synchronized on engine
	Rnd    uint32 `protobuf:"varint,3,opt,name=rnd,proto3" json:"rnd,omitempty"`   // Fill with logical clock on relay
	// Data
	Roll  int32  `protobuf:"varint,4,opt,name=roll,proto3" json:"roll,omitempty"`
	Pitch int32  `protobuf:"varint,5,opt,name=pitch,proto3" json:"pitch,omitempty"`
	Yaw   int32  `protobuf:"varint,6,opt,name=yaw,proto3" json:"yaw,omitempty"`
	X     int32  `protobuf:"varint,7,opt,name=x,proto3" json:"x,omitempty"`
	Y     int32  `protobuf:"varint,8,opt,name=y,proto3" json:"y,omitempty"`
	Z     int32  `protobuf:"varint,9,opt,name=z,proto3" json:"z,omitempty"`
	Index uint32 `protobuf:"varint,10,opt,name=index,proto3" json:"index,omitempty"` // 0 == reset
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_main_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_protos_main_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_protos_main_proto_rawDescGZIP(), []int{0}
}

func (x *Data) GetPlayer() uint32 {
	if x != nil {
		return x.Player
	}
	return 0
}

func (x *Data) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Data) GetRnd() uint32 {
	if x != nil {
		return x.Rnd
	}
	return 0
}

func (x *Data) GetRoll() int32 {
	if x != nil {
		return x.Roll
	}
	return 0
}

func (x *Data) GetPitch() int32 {
	if x != nil {
		return x.Pitch
	}
	return 0
}

func (x *Data) GetYaw() int32 {
	if x != nil {
		return x.Yaw
	}
	return 0
}

func (x *Data) GetX() int32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Data) GetY() int32 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *Data) GetZ() int32 {
	if x != nil {
		return x.Z
	}
	return 0
}

func (x *Data) GetIndex() uint32 {
	if x != nil {
		return x.Index
	}
	return 0
}

type SensorData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Data `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *SensorData) Reset() {
	*x = SensorData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_main_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SensorData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SensorData) ProtoMessage() {}

func (x *SensorData) ProtoReflect() protoreflect.Message {
	mi := &file_protos_main_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SensorData.ProtoReflect.Descriptor instead.
func (*SensorData) Descriptor() ([]byte, []int) {
	return file_protos_main_proto_rawDescGZIP(), []int{1}
}

func (x *SensorData) GetData() []*Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type RndResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rnd uint32 `protobuf:"varint,1,opt,name=rnd,proto3" json:"rnd,omitempty"`
}

func (x *RndResp) Reset() {
	*x = RndResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_main_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RndResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RndResp) ProtoMessage() {}

func (x *RndResp) ProtoReflect() protoreflect.Message {
	mi := &file_protos_main_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RndResp.ProtoReflect.Descriptor instead.
func (*RndResp) Descriptor() ([]byte, []int) {
	return file_protos_main_proto_rawDescGZIP(), []int{2}
}

func (x *RndResp) GetRnd() uint32 {
	if x != nil {
		return x.Rnd
	}
	return 0
}

type InFovResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Player uint32 `protobuf:"varint,1,opt,name=player,proto3" json:"player,omitempty"`
	Time   uint64 `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	InFov  bool   `protobuf:"varint,3,opt,name=inFov,proto3" json:"inFov,omitempty"`
}

func (x *InFovResp) Reset() {
	*x = InFovResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_main_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InFovResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InFovResp) ProtoMessage() {}

func (x *InFovResp) ProtoReflect() protoreflect.Message {
	mi := &file_protos_main_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InFovResp.ProtoReflect.Descriptor instead.
func (*InFovResp) Descriptor() ([]byte, []int) {
	return file_protos_main_proto_rawDescGZIP(), []int{3}
}

func (x *InFovResp) GetPlayer() uint32 {
	if x != nil {
		return x.Player
	}
	return 0
}

func (x *InFovResp) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *InFovResp) GetInFov() bool {
	if x != nil {
		return x.InFov
	}
	return false
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Player uint32 `protobuf:"varint,1,opt,name=player,proto3" json:"player,omitempty"`
	Time   uint64 `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	Rnd    uint32 `protobuf:"varint,3,opt,name=rnd,proto3" json:"rnd,omitempty"`
	Action Action `protobuf:"varint,4,opt,name=action,proto3,enum=Action" json:"action,omitempty"`
	// Action specific
	// OPTIMIZE oneof
	ShootID uint32 `protobuf:"varint,5,opt,name=shootID,proto3" json:"shootID,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_main_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_protos_main_proto_msgTypes[4]
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
	return file_protos_main_proto_rawDescGZIP(), []int{4}
}

func (x *Event) GetPlayer() uint32 {
	if x != nil {
		return x.Player
	}
	return 0
}

func (x *Event) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Event) GetRnd() uint32 {
	if x != nil {
		return x.Rnd
	}
	return 0
}

func (x *Event) GetAction() Action {
	if x != nil {
		return x.Action
	}
	return Action_none
}

func (x *Event) GetShootID() uint32 {
	if x != nil {
		return x.ShootID
	}
	return 0
}

type State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	P1 *PlayerState `protobuf:"bytes,1,opt,name=p1,proto3" json:"p1,omitempty"`
	P2 *PlayerState `protobuf:"bytes,2,opt,name=p2,proto3" json:"p2,omitempty"`
}

func (x *State) Reset() {
	*x = State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_main_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_protos_main_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_protos_main_proto_rawDescGZIP(), []int{5}
}

func (x *State) GetP1() *PlayerState {
	if x != nil {
		return x.P1
	}
	return nil
}

func (x *State) GetP2() *PlayerState {
	if x != nil {
		return x.P2
	}
	return nil
}

type PlayerState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hp           uint32  `protobuf:"varint,1,opt,name=hp,proto3" json:"hp,omitempty"`
	Action       Action  `protobuf:"varint,2,opt,name=action,proto3,enum=Action" json:"action,omitempty"`
	Bullets      uint32  `protobuf:"varint,3,opt,name=bullets,proto3" json:"bullets,omitempty"`
	Grenades     uint32  `protobuf:"varint,4,opt,name=grenades,proto3" json:"grenades,omitempty"`
	ShieldTime   float64 `protobuf:"fixed64,5,opt,name=shield_time,proto3" json:"shield_time,omitempty"`
	ShieldHealth uint32  `protobuf:"varint,6,opt,name=shield_health,proto3" json:"shield_health,omitempty"`
	NumDeaths    uint32  `protobuf:"varint,7,opt,name=num_deaths,proto3" json:"num_deaths,omitempty"`
	NumShield    uint32  `protobuf:"varint,8,opt,name=num_shield,proto3" json:"num_shield,omitempty"`
}

func (x *PlayerState) Reset() {
	*x = PlayerState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_main_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerState) ProtoMessage() {}

func (x *PlayerState) ProtoReflect() protoreflect.Message {
	mi := &file_protos_main_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerState.ProtoReflect.Descriptor instead.
func (*PlayerState) Descriptor() ([]byte, []int) {
	return file_protos_main_proto_rawDescGZIP(), []int{6}
}

func (x *PlayerState) GetHp() uint32 {
	if x != nil {
		return x.Hp
	}
	return 0
}

func (x *PlayerState) GetAction() Action {
	if x != nil {
		return x.Action
	}
	return Action_none
}

func (x *PlayerState) GetBullets() uint32 {
	if x != nil {
		return x.Bullets
	}
	return 0
}

func (x *PlayerState) GetGrenades() uint32 {
	if x != nil {
		return x.Grenades
	}
	return 0
}

func (x *PlayerState) GetShieldTime() float64 {
	if x != nil {
		return x.ShieldTime
	}
	return 0
}

func (x *PlayerState) GetShieldHealth() uint32 {
	if x != nil {
		return x.ShieldHealth
	}
	return 0
}

func (x *PlayerState) GetNumDeaths() uint32 {
	if x != nil {
		return x.NumDeaths
	}
	return 0
}

func (x *PlayerState) GetNumShield() uint32 {
	if x != nil {
		return x.NumShield
	}
	return 0
}

var File_protos_main_proto protoreflect.FileDescriptor

var file_protos_main_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xc0, 0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x03, 0x72, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x6c, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x70,
	0x69, 0x74, 0x63, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x70, 0x69, 0x74, 0x63,
	0x68, 0x12, 0x10, 0x0a, 0x03, 0x79, 0x61, 0x77, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x79, 0x61, 0x77, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01,
	0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x79, 0x12,
	0x0c, 0x0a, 0x01, 0x7a, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x7a, 0x12, 0x14, 0x0a,
	0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x22, 0x27, 0x0a, 0x0a, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x19, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x05, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x1b, 0x0a, 0x07,
	0x52, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x6e, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x72, 0x6e, 0x64, 0x22, 0x4d, 0x0a, 0x09, 0x49, 0x6e, 0x46,
	0x6f, 0x76, 0x52, 0x65, 0x73, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x69,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x46, 0x6f, 0x76, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x05, 0x69, 0x6e, 0x46, 0x6f, 0x76, 0x22, 0x80, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x72, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x72, 0x6e, 0x64,
	0x12, 0x1f, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x07, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x68, 0x6f, 0x6f, 0x74, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x73, 0x68, 0x6f, 0x6f, 0x74, 0x49, 0x44, 0x22, 0x43, 0x0a, 0x05, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x02, 0x70, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x02,
	0x70, 0x31, 0x12, 0x1c, 0x0a, 0x02, 0x70, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x02, 0x70, 0x32,
	0x22, 0xfc, 0x01, 0x0a, 0x0b, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x68, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x68, 0x70,
	0x12, 0x1f, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x07, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x62, 0x75, 0x6c, 0x6c, 0x65, 0x74, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x67,
	0x72, 0x65, 0x6e, 0x61, 0x64, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x67,
	0x72, 0x65, 0x6e, 0x61, 0x64, 0x65, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x68, 0x69, 0x65, 0x6c,
	0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x73, 0x68,
	0x69, 0x65, 0x6c, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x68, 0x69,
	0x65, 0x6c, 0x64, 0x5f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0d, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12,
	0x1e, 0x0a, 0x0a, 0x6e, 0x75, 0x6d, 0x5f, 0x64, 0x65, 0x61, 0x74, 0x68, 0x73, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0a, 0x6e, 0x75, 0x6d, 0x5f, 0x64, 0x65, 0x61, 0x74, 0x68, 0x73, 0x12,
	0x1e, 0x0a, 0x0a, 0x6e, 0x75, 0x6d, 0x5f, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0a, 0x6e, 0x75, 0x6d, 0x5f, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x2a,
	0x89, 0x01, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x08, 0x0a, 0x04, 0x6e, 0x6f,
	0x6e, 0x65, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x67, 0x72, 0x65, 0x6e, 0x61, 0x64, 0x65, 0x10,
	0x01, 0x12, 0x0a, 0x0a, 0x06, 0x72, 0x65, 0x6c, 0x6f, 0x61, 0x64, 0x10, 0x02, 0x12, 0x09, 0x0a,
	0x05, 0x73, 0x68, 0x6f, 0x6f, 0x74, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x6c, 0x6f, 0x67, 0x6f,
	0x75, 0x74, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x10, 0x05,
	0x12, 0x08, 0x0a, 0x04, 0x73, 0x68, 0x6f, 0x74, 0x10, 0x06, 0x12, 0x0c, 0x0a, 0x08, 0x67, 0x72,
	0x65, 0x6e, 0x61, 0x64, 0x65, 0x64, 0x10, 0x07, 0x12, 0x13, 0x0a, 0x0f, 0x73, 0x68, 0x69, 0x65,
	0x6c, 0x64, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x10, 0x08, 0x12, 0x0c, 0x0a,
	0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x46, 0x6f, 0x76, 0x10, 0x09, 0x32, 0xbc, 0x01, 0x0a, 0x05,
	0x52, 0x65, 0x6c, 0x61, 0x79, 0x12, 0x30, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x75, 0x6e,
	0x64, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x08, 0x2e, 0x52, 0x6e, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x22, 0x00, 0x30, 0x01, 0x12, 0x2c, 0x0a, 0x07, 0x47, 0x65, 0x73, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x05, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x28, 0x01, 0x12, 0x29, 0x0a, 0x05, 0x53, 0x68, 0x6f, 0x6f, 0x74, 0x12, 0x06,
	0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x12, 0x28, 0x0a, 0x04, 0x53, 0x68, 0x6f, 0x74, 0x12, 0x06, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x32, 0x50, 0x0a, 0x03, 0x56, 0x69,
	0x7a, 0x12, 0x2a, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x06, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x1d, 0x0a,
	0x05, 0x49, 0x6e, 0x46, 0x6f, 0x76, 0x12, 0x06, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x0a,
	0x2e, 0x49, 0x6e, 0x46, 0x6f, 0x76, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x32, 0x49, 0x0a, 0x04,
	0x50, 0x79, 0x6e, 0x71, 0x12, 0x17, 0x0a, 0x04, 0x45, 0x6d, 0x69, 0x74, 0x12, 0x05, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x1a, 0x06, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x28, 0x0a,
	0x04, 0x50, 0x6f, 0x6c, 0x6c, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x06, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_main_proto_rawDescOnce sync.Once
	file_protos_main_proto_rawDescData = file_protos_main_proto_rawDesc
)

func file_protos_main_proto_rawDescGZIP() []byte {
	file_protos_main_proto_rawDescOnce.Do(func() {
		file_protos_main_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_main_proto_rawDescData)
	})
	return file_protos_main_proto_rawDescData
}

var file_protos_main_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_protos_main_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_protos_main_proto_goTypes = []interface{}{
	(Action)(0),           // 0: Action
	(*Data)(nil),          // 1: Data
	(*SensorData)(nil),    // 2: SensorData
	(*RndResp)(nil),       // 3: RndResp
	(*InFovResp)(nil),     // 4: InFovResp
	(*Event)(nil),         // 5: Event
	(*State)(nil),         // 6: State
	(*PlayerState)(nil),   // 7: PlayerState
	(*emptypb.Empty)(nil), // 8: google.protobuf.Empty
}
var file_protos_main_proto_depIdxs = []int32{
	1,  // 0: SensorData.data:type_name -> Data
	0,  // 1: Event.action:type_name -> Action
	7,  // 2: State.p1:type_name -> PlayerState
	7,  // 3: State.p2:type_name -> PlayerState
	0,  // 4: PlayerState.action:type_name -> Action
	8,  // 5: Relay.GetRound:input_type -> google.protobuf.Empty
	1,  // 6: Relay.Gesture:input_type -> Data
	5,  // 7: Relay.Shoot:input_type -> Event
	5,  // 8: Relay.Shot:input_type -> Event
	6,  // 9: Viz.Update:input_type -> State
	5,  // 10: Viz.InFov:input_type -> Event
	1,  // 11: Pynq.Emit:input_type -> Data
	8,  // 12: Pynq.Poll:input_type -> google.protobuf.Empty
	3,  // 13: Relay.GetRound:output_type -> RndResp
	8,  // 14: Relay.Gesture:output_type -> google.protobuf.Empty
	8,  // 15: Relay.Shoot:output_type -> google.protobuf.Empty
	8,  // 16: Relay.Shot:output_type -> google.protobuf.Empty
	8,  // 17: Viz.Update:output_type -> google.protobuf.Empty
	4,  // 18: Viz.InFov:output_type -> InFovResp
	5,  // 19: Pynq.Emit:output_type -> Event
	5,  // 20: Pynq.Poll:output_type -> Event
	13, // [13:21] is the sub-list for method output_type
	5,  // [5:13] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_protos_main_proto_init() }
func file_protos_main_proto_init() {
	if File_protos_main_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_main_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
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
		file_protos_main_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SensorData); i {
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
		file_protos_main_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RndResp); i {
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
		file_protos_main_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InFovResp); i {
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
		file_protos_main_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_protos_main_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*State); i {
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
		file_protos_main_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerState); i {
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
			RawDescriptor: file_protos_main_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_protos_main_proto_goTypes,
		DependencyIndexes: file_protos_main_proto_depIdxs,
		EnumInfos:         file_protos_main_proto_enumTypes,
		MessageInfos:      file_protos_main_proto_msgTypes,
	}.Build()
	File_protos_main_proto = out.File
	file_protos_main_proto_rawDesc = nil
	file_protos_main_proto_goTypes = nil
	file_protos_main_proto_depIdxs = nil
}
