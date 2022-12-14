# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: main.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\nmain.proto\x1a\x1bgoogle/protobuf/empty.proto\"\x8b\x01\n\x04\x44\x61ta\x12\x0e\n\x06player\x18\x01 \x01(\r\x12\x0c\n\x04time\x18\x02 \x01(\x04\x12\x0b\n\x03rnd\x18\x03 \x01(\r\x12\x0c\n\x04roll\x18\x04 \x01(\x05\x12\r\n\x05pitch\x18\x05 \x01(\x05\x12\x0b\n\x03yaw\x18\x06 \x01(\x05\x12\t\n\x01x\x18\x07 \x01(\x05\x12\t\n\x01y\x18\x08 \x01(\x05\x12\t\n\x01z\x18\t \x01(\x05\x12\r\n\x05index\x18\n \x01(\r\"!\n\nSensorData\x12\x13\n\x04\x64\x61ta\x18\x01 \x03(\x0b\x32\x05.Data\"\x16\n\x07RndResp\x12\x0b\n\x03rnd\x18\x01 \x01(\r\"E\n\tInFovResp\x12\x0e\n\x06player\x18\x01 \x01(\r\x12\x0c\n\x04time\x18\x02 \x01(\x04\x12\x0b\n\x03rnd\x18\x03 \x01(\r\x12\r\n\x05inFov\x18\x04 \x01(\x08\"\\\n\x05\x45vent\x12\x0e\n\x06player\x18\x01 \x01(\r\x12\x0c\n\x04time\x18\x02 \x01(\x04\x12\x0b\n\x03rnd\x18\x03 \x01(\r\x12\x17\n\x06\x61\x63tion\x18\x04 \x01(\x0e\x32\x07.Action\x12\x0f\n\x07shootID\x18\x05 \x01(\r\"C\n\x05State\x12\x1c\n\x02p1\x18\x01 \x01(\x0b\x32\x0c.PlayerStateR\x02p1\x12\x1c\n\x02p2\x18\x02 \x01(\x0b\x32\x0c.PlayerStateR\x02p2\"\xfc\x01\n\x0bPlayerState\x12\x0e\n\x02hp\x18\x01 \x01(\rR\x02hp\x12\x1f\n\x06\x61\x63tion\x18\x02 \x01(\x0e\x32\x07.ActionR\x06\x61\x63tion\x12\x18\n\x07\x62ullets\x18\x03 \x01(\rR\x07\x62ullets\x12\x1a\n\x08grenades\x18\x04 \x01(\rR\x08grenades\x12 \n\x0bshield_time\x18\x05 \x01(\x01R\x0bshield_time\x12$\n\rshield_health\x18\x06 \x01(\rR\rshield_health\x12\x1e\n\nnum_deaths\x18\x07 \x01(\rR\nnum_deaths\x12\x1e\n\nnum_shield\x18\x08 \x01(\rR\nnum_shield*\x93\x01\n\x06\x41\x63tion\x12\x08\n\x04none\x10\x00\x12\x0b\n\x07grenade\x10\x01\x12\n\n\x06reload\x10\x02\x12\t\n\x05shoot\x10\x03\x12\n\n\x06logout\x10\x04\x12\n\n\x06shield\x10\x05\x12\x08\n\x04shot\x10\x06\x12\x0c\n\x08grenaded\x10\x07\x12\x13\n\x0fshieldAvailable\x10\x08\x12\x0c\n\x08\x63heckFov\x10\t\x12\x08\n\x04\x64one\x10\n2\xc0\x01\n\x05Relay\x12\x30\n\x08GetRound\x12\x16.google.protobuf.Empty\x1a\x08.RndResp\"\x00\x30\x01\x12,\n\x07Gesture\x12\x05.Data\x1a\x16.google.protobuf.Empty\"\x00(\x01\x12+\n\x05Shoot\x12\x06.Event\x1a\x16.google.protobuf.Empty\"\x00(\x01\x12*\n\x04Shot\x12\x06.Event\x1a\x16.google.protobuf.Empty\"\x00(\x01\x32P\n\x03Viz\x12*\n\x06Update\x12\x06.State\x1a\x16.google.protobuf.Empty\"\x00\x12\x1d\n\x05InFov\x12\x06.Event\x1a\n.InFovResp\"\x00\x32I\n\x04Pynq\x12\x17\n\x04\x45mit\x12\x05.Data\x1a\x06.Event\"\x00\x12(\n\x04Poll\x12\x16.google.protobuf.Empty\x1a\x06.Event\"\x00\x42\x04Z\x02./b\x06proto3')

_ACTION = DESCRIPTOR.enum_types_by_name['Action']
Action = enum_type_wrapper.EnumTypeWrapper(_ACTION)
none = 0
grenade = 1
reload = 2
shoot = 3
logout = 4
shield = 5
shot = 6
grenaded = 7
shieldAvailable = 8
checkFov = 9
done = 10


_DATA = DESCRIPTOR.message_types_by_name['Data']
_SENSORDATA = DESCRIPTOR.message_types_by_name['SensorData']
_RNDRESP = DESCRIPTOR.message_types_by_name['RndResp']
_INFOVRESP = DESCRIPTOR.message_types_by_name['InFovResp']
_EVENT = DESCRIPTOR.message_types_by_name['Event']
_STATE = DESCRIPTOR.message_types_by_name['State']
_PLAYERSTATE = DESCRIPTOR.message_types_by_name['PlayerState']
Data = _reflection.GeneratedProtocolMessageType('Data', (_message.Message,), {
  'DESCRIPTOR' : _DATA,
  '__module__' : 'main_pb2'
  # @@protoc_insertion_point(class_scope:Data)
  })
_sym_db.RegisterMessage(Data)

SensorData = _reflection.GeneratedProtocolMessageType('SensorData', (_message.Message,), {
  'DESCRIPTOR' : _SENSORDATA,
  '__module__' : 'main_pb2'
  # @@protoc_insertion_point(class_scope:SensorData)
  })
_sym_db.RegisterMessage(SensorData)

RndResp = _reflection.GeneratedProtocolMessageType('RndResp', (_message.Message,), {
  'DESCRIPTOR' : _RNDRESP,
  '__module__' : 'main_pb2'
  # @@protoc_insertion_point(class_scope:RndResp)
  })
_sym_db.RegisterMessage(RndResp)

InFovResp = _reflection.GeneratedProtocolMessageType('InFovResp', (_message.Message,), {
  'DESCRIPTOR' : _INFOVRESP,
  '__module__' : 'main_pb2'
  # @@protoc_insertion_point(class_scope:InFovResp)
  })
_sym_db.RegisterMessage(InFovResp)

Event = _reflection.GeneratedProtocolMessageType('Event', (_message.Message,), {
  'DESCRIPTOR' : _EVENT,
  '__module__' : 'main_pb2'
  # @@protoc_insertion_point(class_scope:Event)
  })
_sym_db.RegisterMessage(Event)

State = _reflection.GeneratedProtocolMessageType('State', (_message.Message,), {
  'DESCRIPTOR' : _STATE,
  '__module__' : 'main_pb2'
  # @@protoc_insertion_point(class_scope:State)
  })
_sym_db.RegisterMessage(State)

PlayerState = _reflection.GeneratedProtocolMessageType('PlayerState', (_message.Message,), {
  'DESCRIPTOR' : _PLAYERSTATE,
  '__module__' : 'main_pb2'
  # @@protoc_insertion_point(class_scope:PlayerState)
  })
_sym_db.RegisterMessage(PlayerState)

_RELAY = DESCRIPTOR.services_by_name['Relay']
_VIZ = DESCRIPTOR.services_by_name['Viz']
_PYNQ = DESCRIPTOR.services_by_name['Pynq']
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\002./'
  _ACTION._serialized_start=734
  _ACTION._serialized_end=881
  _DATA._serialized_start=44
  _DATA._serialized_end=183
  _SENSORDATA._serialized_start=185
  _SENSORDATA._serialized_end=218
  _RNDRESP._serialized_start=220
  _RNDRESP._serialized_end=242
  _INFOVRESP._serialized_start=244
  _INFOVRESP._serialized_end=313
  _EVENT._serialized_start=315
  _EVENT._serialized_end=407
  _STATE._serialized_start=409
  _STATE._serialized_end=476
  _PLAYERSTATE._serialized_start=479
  _PLAYERSTATE._serialized_end=731
  _RELAY._serialized_start=884
  _RELAY._serialized_end=1076
  _VIZ._serialized_start=1078
  _VIZ._serialized_end=1158
  _PYNQ._serialized_start=1160
  _PYNQ._serialized_end=1233
# @@protoc_insertion_point(module_scope)
