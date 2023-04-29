# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: instance.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
import imrpc_pb2 as imrpc__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='instance.proto',
  package='imrpc',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\x0einstance.proto\x12\x05imrpc\x1a\x1bgoogle/protobuf/empty.proto\x1a\x0bimrpc.proto\"/\n\x0fProcessSpecific\x12\x0e\n\x06\x62inary\x18\x01 \x01(\t\x12\x0c\n\x04\x61rgs\x18\x02 \x03(\t\"\xc7\x01\n\x0cSpdkSpecific\x12\x0f\n\x07\x61\x64\x64ress\x18\x01 \x01(\t\x12G\n\x13replica_address_map\x18\x02 \x03(\x0b\x32*.imrpc.SpdkSpecific.ReplicaAddressMapEntry\x12\x10\n\x08\x66rontend\x18\x03 \x01(\t\x12\x11\n\tdisk_uuid\x18\x04 \x01(\t\x1a\x38\n\x16ReplicaAddressMapEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"\xf0\x01\n\x0cInstanceSpec\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04type\x18\x02 \x01(\t\x12\x1c\n\x14\x62\x61\x63kend_store_driver\x18\x03 \x01(\t\x12\x13\n\x0bvolume_name\x18\x04 \x01(\t\x12\x0c\n\x04size\x18\x05 \x01(\x04\x12\x12\n\nport_count\x18\x06 \x01(\x05\x12\x11\n\tport_args\x18\x07 \x03(\t\x12\x30\n\x10process_specific\x18\x08 \x01(\x0b\x32\x16.imrpc.ProcessSpecific\x12*\n\rspdk_specific\x18\t \x01(\x0b\x32\x13.imrpc.SpdkSpecific\"X\n\x0eInstanceStatus\x12\r\n\x05state\x18\x01 \x01(\t\x12\x11\n\terror_msg\x18\x02 \x01(\t\x12\x12\n\nport_start\x18\x03 \x01(\x05\x12\x10\n\x08port_end\x18\x04 \x01(\x05\":\n\x15InstanceCreateRequest\x12!\n\x04spec\x18\x01 \x01(\x0b\x32\x13.imrpc.InstanceSpec\"~\n\x15InstanceDeleteRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04type\x18\x02 \x01(\t\x12\x1c\n\x14\x62\x61\x63kend_store_driver\x18\x03 \x01(\t\x12\x11\n\tdisk_uuid\x18\x04 \x01(\t\x12\x18\n\x10\x63leanup_required\x18\x05 \x01(\x08\"a\n\x12InstanceGetRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04type\x18\x02 \x01(\t\x12\x1c\n\x14\x62\x61\x63kend_store_driver\x18\x03 \x01(\t\x12\x11\n\tdisk_uuid\x18\x04 \x01(\t\"m\n\x10InstanceResponse\x12!\n\x04spec\x18\x01 \x01(\x0b\x32\x13.imrpc.InstanceSpec\x12%\n\x06status\x18\x02 \x01(\x0b\x32\x15.imrpc.InstanceStatus\x12\x0f\n\x07\x64\x65leted\x18\x03 \x01(\x08\"\xa0\x01\n\x14InstanceListResponse\x12=\n\tinstances\x18\x01 \x03(\x0b\x32*.imrpc.InstanceListResponse.InstancesEntry\x1aI\n\x0eInstancesEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12&\n\x05value\x18\x02 \x01(\x0b\x32\x17.imrpc.InstanceResponse:\x02\x38\x01\"N\n\x12InstanceLogRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04type\x18\x02 \x01(\t\x12\x1c\n\x14\x62\x61\x63kend_store_driver\x18\x03 \x01(\t\"U\n\x16InstanceReplaceRequest\x12!\n\x04spec\x18\x01 \x01(\x0b\x32\x13.imrpc.InstanceSpec\x12\x18\n\x10terminate_signal\x18\x02 \x01(\t2\xba\x04\n\x0fInstanceService\x12I\n\x0eInstanceCreate\x12\x1c.imrpc.InstanceCreateRequest\x1a\x17.imrpc.InstanceResponse\"\x00\x12I\n\x0eInstanceDelete\x12\x1c.imrpc.InstanceDeleteRequest\x1a\x17.imrpc.InstanceResponse\"\x00\x12\x43\n\x0bInstanceGet\x12\x19.imrpc.InstanceGetRequest\x1a\x17.imrpc.InstanceResponse\"\x00\x12\x45\n\x0cInstanceList\x12\x16.google.protobuf.Empty\x1a\x1b.imrpc.InstanceListResponse\"\x00\x12:\n\x0bInstanceLog\x12\x19.imrpc.InstanceLogRequest\x1a\x0c.LogResponse\"\x00\x30\x01\x12\x44\n\rInstanceWatch\x12\x16.google.protobuf.Empty\x1a\x17.imrpc.InstanceResponse\"\x00\x30\x01\x12K\n\x0fInstanceReplace\x12\x1d.imrpc.InstanceReplaceRequest\x1a\x17.imrpc.InstanceResponse\"\x00\x12\x36\n\nVersionGet\x12\x16.google.protobuf.Empty\x1a\x10.VersionResponseb\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_empty__pb2.DESCRIPTOR,imrpc__pb2.DESCRIPTOR,])




_PROCESSSPECIFIC = _descriptor.Descriptor(
  name='ProcessSpecific',
  full_name='imrpc.ProcessSpecific',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='binary', full_name='imrpc.ProcessSpecific.binary', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='args', full_name='imrpc.ProcessSpecific.args', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=67,
  serialized_end=114,
)


_SPDKSPECIFIC_REPLICAADDRESSMAPENTRY = _descriptor.Descriptor(
  name='ReplicaAddressMapEntry',
  full_name='imrpc.SpdkSpecific.ReplicaAddressMapEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='imrpc.SpdkSpecific.ReplicaAddressMapEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='imrpc.SpdkSpecific.ReplicaAddressMapEntry.value', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=_b('8\001'),
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=260,
  serialized_end=316,
)

_SPDKSPECIFIC = _descriptor.Descriptor(
  name='SpdkSpecific',
  full_name='imrpc.SpdkSpecific',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='address', full_name='imrpc.SpdkSpecific.address', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='replica_address_map', full_name='imrpc.SpdkSpecific.replica_address_map', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='frontend', full_name='imrpc.SpdkSpecific.frontend', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='disk_uuid', full_name='imrpc.SpdkSpecific.disk_uuid', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_SPDKSPECIFIC_REPLICAADDRESSMAPENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=117,
  serialized_end=316,
)


_INSTANCESPEC = _descriptor.Descriptor(
  name='InstanceSpec',
  full_name='imrpc.InstanceSpec',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='imrpc.InstanceSpec.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='imrpc.InstanceSpec.type', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='backend_store_driver', full_name='imrpc.InstanceSpec.backend_store_driver', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='volume_name', full_name='imrpc.InstanceSpec.volume_name', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='size', full_name='imrpc.InstanceSpec.size', index=4,
      number=5, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='port_count', full_name='imrpc.InstanceSpec.port_count', index=5,
      number=6, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='port_args', full_name='imrpc.InstanceSpec.port_args', index=6,
      number=7, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='process_specific', full_name='imrpc.InstanceSpec.process_specific', index=7,
      number=8, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='spdk_specific', full_name='imrpc.InstanceSpec.spdk_specific', index=8,
      number=9, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=319,
  serialized_end=559,
)


_INSTANCESTATUS = _descriptor.Descriptor(
  name='InstanceStatus',
  full_name='imrpc.InstanceStatus',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='state', full_name='imrpc.InstanceStatus.state', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='error_msg', full_name='imrpc.InstanceStatus.error_msg', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='port_start', full_name='imrpc.InstanceStatus.port_start', index=2,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='port_end', full_name='imrpc.InstanceStatus.port_end', index=3,
      number=4, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=561,
  serialized_end=649,
)


_INSTANCECREATEREQUEST = _descriptor.Descriptor(
  name='InstanceCreateRequest',
  full_name='imrpc.InstanceCreateRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='spec', full_name='imrpc.InstanceCreateRequest.spec', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=651,
  serialized_end=709,
)


_INSTANCEDELETEREQUEST = _descriptor.Descriptor(
  name='InstanceDeleteRequest',
  full_name='imrpc.InstanceDeleteRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='imrpc.InstanceDeleteRequest.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='imrpc.InstanceDeleteRequest.type', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='backend_store_driver', full_name='imrpc.InstanceDeleteRequest.backend_store_driver', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='disk_uuid', full_name='imrpc.InstanceDeleteRequest.disk_uuid', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='cleanup_required', full_name='imrpc.InstanceDeleteRequest.cleanup_required', index=4,
      number=5, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=711,
  serialized_end=837,
)


_INSTANCEGETREQUEST = _descriptor.Descriptor(
  name='InstanceGetRequest',
  full_name='imrpc.InstanceGetRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='imrpc.InstanceGetRequest.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='imrpc.InstanceGetRequest.type', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='backend_store_driver', full_name='imrpc.InstanceGetRequest.backend_store_driver', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='disk_uuid', full_name='imrpc.InstanceGetRequest.disk_uuid', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=839,
  serialized_end=936,
)


_INSTANCERESPONSE = _descriptor.Descriptor(
  name='InstanceResponse',
  full_name='imrpc.InstanceResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='spec', full_name='imrpc.InstanceResponse.spec', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='status', full_name='imrpc.InstanceResponse.status', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='deleted', full_name='imrpc.InstanceResponse.deleted', index=2,
      number=3, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=938,
  serialized_end=1047,
)


_INSTANCELISTRESPONSE_INSTANCESENTRY = _descriptor.Descriptor(
  name='InstancesEntry',
  full_name='imrpc.InstanceListResponse.InstancesEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='imrpc.InstanceListResponse.InstancesEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='imrpc.InstanceListResponse.InstancesEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=_b('8\001'),
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1137,
  serialized_end=1210,
)

_INSTANCELISTRESPONSE = _descriptor.Descriptor(
  name='InstanceListResponse',
  full_name='imrpc.InstanceListResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='instances', full_name='imrpc.InstanceListResponse.instances', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_INSTANCELISTRESPONSE_INSTANCESENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1050,
  serialized_end=1210,
)


_INSTANCELOGREQUEST = _descriptor.Descriptor(
  name='InstanceLogRequest',
  full_name='imrpc.InstanceLogRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='imrpc.InstanceLogRequest.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='imrpc.InstanceLogRequest.type', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='backend_store_driver', full_name='imrpc.InstanceLogRequest.backend_store_driver', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1212,
  serialized_end=1290,
)


_INSTANCEREPLACEREQUEST = _descriptor.Descriptor(
  name='InstanceReplaceRequest',
  full_name='imrpc.InstanceReplaceRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='spec', full_name='imrpc.InstanceReplaceRequest.spec', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='terminate_signal', full_name='imrpc.InstanceReplaceRequest.terminate_signal', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1292,
  serialized_end=1377,
)

_SPDKSPECIFIC_REPLICAADDRESSMAPENTRY.containing_type = _SPDKSPECIFIC
_SPDKSPECIFIC.fields_by_name['replica_address_map'].message_type = _SPDKSPECIFIC_REPLICAADDRESSMAPENTRY
_INSTANCESPEC.fields_by_name['process_specific'].message_type = _PROCESSSPECIFIC
_INSTANCESPEC.fields_by_name['spdk_specific'].message_type = _SPDKSPECIFIC
_INSTANCECREATEREQUEST.fields_by_name['spec'].message_type = _INSTANCESPEC
_INSTANCERESPONSE.fields_by_name['spec'].message_type = _INSTANCESPEC
_INSTANCERESPONSE.fields_by_name['status'].message_type = _INSTANCESTATUS
_INSTANCELISTRESPONSE_INSTANCESENTRY.fields_by_name['value'].message_type = _INSTANCERESPONSE
_INSTANCELISTRESPONSE_INSTANCESENTRY.containing_type = _INSTANCELISTRESPONSE
_INSTANCELISTRESPONSE.fields_by_name['instances'].message_type = _INSTANCELISTRESPONSE_INSTANCESENTRY
_INSTANCEREPLACEREQUEST.fields_by_name['spec'].message_type = _INSTANCESPEC
DESCRIPTOR.message_types_by_name['ProcessSpecific'] = _PROCESSSPECIFIC
DESCRIPTOR.message_types_by_name['SpdkSpecific'] = _SPDKSPECIFIC
DESCRIPTOR.message_types_by_name['InstanceSpec'] = _INSTANCESPEC
DESCRIPTOR.message_types_by_name['InstanceStatus'] = _INSTANCESTATUS
DESCRIPTOR.message_types_by_name['InstanceCreateRequest'] = _INSTANCECREATEREQUEST
DESCRIPTOR.message_types_by_name['InstanceDeleteRequest'] = _INSTANCEDELETEREQUEST
DESCRIPTOR.message_types_by_name['InstanceGetRequest'] = _INSTANCEGETREQUEST
DESCRIPTOR.message_types_by_name['InstanceResponse'] = _INSTANCERESPONSE
DESCRIPTOR.message_types_by_name['InstanceListResponse'] = _INSTANCELISTRESPONSE
DESCRIPTOR.message_types_by_name['InstanceLogRequest'] = _INSTANCELOGREQUEST
DESCRIPTOR.message_types_by_name['InstanceReplaceRequest'] = _INSTANCEREPLACEREQUEST
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

ProcessSpecific = _reflection.GeneratedProtocolMessageType('ProcessSpecific', (_message.Message,), {
  'DESCRIPTOR' : _PROCESSSPECIFIC,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.ProcessSpecific)
  })
_sym_db.RegisterMessage(ProcessSpecific)

SpdkSpecific = _reflection.GeneratedProtocolMessageType('SpdkSpecific', (_message.Message,), {

  'ReplicaAddressMapEntry' : _reflection.GeneratedProtocolMessageType('ReplicaAddressMapEntry', (_message.Message,), {
    'DESCRIPTOR' : _SPDKSPECIFIC_REPLICAADDRESSMAPENTRY,
    '__module__' : 'instance_pb2'
    # @@protoc_insertion_point(class_scope:imrpc.SpdkSpecific.ReplicaAddressMapEntry)
    })
  ,
  'DESCRIPTOR' : _SPDKSPECIFIC,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.SpdkSpecific)
  })
_sym_db.RegisterMessage(SpdkSpecific)
_sym_db.RegisterMessage(SpdkSpecific.ReplicaAddressMapEntry)

InstanceSpec = _reflection.GeneratedProtocolMessageType('InstanceSpec', (_message.Message,), {
  'DESCRIPTOR' : _INSTANCESPEC,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.InstanceSpec)
  })
_sym_db.RegisterMessage(InstanceSpec)

InstanceStatus = _reflection.GeneratedProtocolMessageType('InstanceStatus', (_message.Message,), {
  'DESCRIPTOR' : _INSTANCESTATUS,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.InstanceStatus)
  })
_sym_db.RegisterMessage(InstanceStatus)

InstanceCreateRequest = _reflection.GeneratedProtocolMessageType('InstanceCreateRequest', (_message.Message,), {
  'DESCRIPTOR' : _INSTANCECREATEREQUEST,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.InstanceCreateRequest)
  })
_sym_db.RegisterMessage(InstanceCreateRequest)

InstanceDeleteRequest = _reflection.GeneratedProtocolMessageType('InstanceDeleteRequest', (_message.Message,), {
  'DESCRIPTOR' : _INSTANCEDELETEREQUEST,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.InstanceDeleteRequest)
  })
_sym_db.RegisterMessage(InstanceDeleteRequest)

InstanceGetRequest = _reflection.GeneratedProtocolMessageType('InstanceGetRequest', (_message.Message,), {
  'DESCRIPTOR' : _INSTANCEGETREQUEST,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.InstanceGetRequest)
  })
_sym_db.RegisterMessage(InstanceGetRequest)

InstanceResponse = _reflection.GeneratedProtocolMessageType('InstanceResponse', (_message.Message,), {
  'DESCRIPTOR' : _INSTANCERESPONSE,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.InstanceResponse)
  })
_sym_db.RegisterMessage(InstanceResponse)

InstanceListResponse = _reflection.GeneratedProtocolMessageType('InstanceListResponse', (_message.Message,), {

  'InstancesEntry' : _reflection.GeneratedProtocolMessageType('InstancesEntry', (_message.Message,), {
    'DESCRIPTOR' : _INSTANCELISTRESPONSE_INSTANCESENTRY,
    '__module__' : 'instance_pb2'
    # @@protoc_insertion_point(class_scope:imrpc.InstanceListResponse.InstancesEntry)
    })
  ,
  'DESCRIPTOR' : _INSTANCELISTRESPONSE,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.InstanceListResponse)
  })
_sym_db.RegisterMessage(InstanceListResponse)
_sym_db.RegisterMessage(InstanceListResponse.InstancesEntry)

InstanceLogRequest = _reflection.GeneratedProtocolMessageType('InstanceLogRequest', (_message.Message,), {
  'DESCRIPTOR' : _INSTANCELOGREQUEST,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.InstanceLogRequest)
  })
_sym_db.RegisterMessage(InstanceLogRequest)

InstanceReplaceRequest = _reflection.GeneratedProtocolMessageType('InstanceReplaceRequest', (_message.Message,), {
  'DESCRIPTOR' : _INSTANCEREPLACEREQUEST,
  '__module__' : 'instance_pb2'
  # @@protoc_insertion_point(class_scope:imrpc.InstanceReplaceRequest)
  })
_sym_db.RegisterMessage(InstanceReplaceRequest)


_SPDKSPECIFIC_REPLICAADDRESSMAPENTRY._options = None
_INSTANCELISTRESPONSE_INSTANCESENTRY._options = None

_INSTANCESERVICE = _descriptor.ServiceDescriptor(
  name='InstanceService',
  full_name='imrpc.InstanceService',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=1380,
  serialized_end=1950,
  methods=[
  _descriptor.MethodDescriptor(
    name='InstanceCreate',
    full_name='imrpc.InstanceService.InstanceCreate',
    index=0,
    containing_service=None,
    input_type=_INSTANCECREATEREQUEST,
    output_type=_INSTANCERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='InstanceDelete',
    full_name='imrpc.InstanceService.InstanceDelete',
    index=1,
    containing_service=None,
    input_type=_INSTANCEDELETEREQUEST,
    output_type=_INSTANCERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='InstanceGet',
    full_name='imrpc.InstanceService.InstanceGet',
    index=2,
    containing_service=None,
    input_type=_INSTANCEGETREQUEST,
    output_type=_INSTANCERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='InstanceList',
    full_name='imrpc.InstanceService.InstanceList',
    index=3,
    containing_service=None,
    input_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    output_type=_INSTANCELISTRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='InstanceLog',
    full_name='imrpc.InstanceService.InstanceLog',
    index=4,
    containing_service=None,
    input_type=_INSTANCELOGREQUEST,
    output_type=imrpc__pb2._LOGRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='InstanceWatch',
    full_name='imrpc.InstanceService.InstanceWatch',
    index=5,
    containing_service=None,
    input_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    output_type=_INSTANCERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='InstanceReplace',
    full_name='imrpc.InstanceService.InstanceReplace',
    index=6,
    containing_service=None,
    input_type=_INSTANCEREPLACEREQUEST,
    output_type=_INSTANCERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='VersionGet',
    full_name='imrpc.InstanceService.VersionGet',
    index=7,
    containing_service=None,
    input_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    output_type=imrpc__pb2._VERSIONRESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_INSTANCESERVICE)

DESCRIPTOR.services_by_name['InstanceService'] = _INSTANCESERVICE

# @@protoc_insertion_point(module_scope)
