# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: github.com/longhorn/longhorn-instance-manager/pkg/imrpc/common.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='github.com/longhorn/longhorn-instance-manager/pkg/imrpc/common.proto',
  package='imrpc',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\nDgithub.com/longhorn/longhorn-instance-manager/pkg/imrpc/common.proto\x12\x05imrpc*,\n\x12\x42\x61\x63kendStoreDriver\x12\x0c\n\x08longhorn\x10\x00\x12\x08\n\x04spdk\x10\x01\x62\x06proto3')
)

_BACKENDSTOREDRIVER = _descriptor.EnumDescriptor(
  name='BackendStoreDriver',
  full_name='imrpc.BackendStoreDriver',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='longhorn', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='spdk', index=1, number=1,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=79,
  serialized_end=123,
)
_sym_db.RegisterEnumDescriptor(_BACKENDSTOREDRIVER)

BackendStoreDriver = enum_type_wrapper.EnumTypeWrapper(_BACKENDSTOREDRIVER)
longhorn = 0
spdk = 1


DESCRIPTOR.enum_types_by_name['BackendStoreDriver'] = _BACKENDSTOREDRIVER
_sym_db.RegisterFileDescriptor(DESCRIPTOR)


# @@protoc_insertion_point(module_scope)
