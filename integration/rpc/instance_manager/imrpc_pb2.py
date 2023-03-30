# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: imrpc.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='imrpc.proto',
  package='',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\x0bimrpc.proto\x1a\x1bgoogle/protobuf/empty.proto\"`\n\x0bProcessSpec\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0e\n\x06\x62inary\x18\x02 \x01(\t\x12\x0c\n\x04\x61rgs\x18\x03 \x03(\t\x12\x12\n\nport_count\x18\x04 \x01(\x05\x12\x11\n\tport_args\x18\x05 \x03(\t\"e\n\rProcessStatus\x12\r\n\x05state\x18\x01 \x01(\t\x12\x11\n\terror_msg\x18\x02 \x01(\t\x12\x12\n\nport_start\x18\x03 \x01(\x05\x12\x10\n\x08port_end\x18\x04 \x01(\x05\x12\x0c\n\x04uuid\x18\x05 \x01(\t\"2\n\x14ProcessCreateRequest\x12\x1a\n\x04spec\x18\x01 \x01(\x0b\x32\x0c.ProcessSpec\"$\n\x14ProcessDeleteRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\"!\n\x11ProcessGetRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\"^\n\x0fProcessResponse\x12\x1a\n\x04spec\x18\x01 \x01(\x0b\x32\x0c.ProcessSpec\x12\x1e\n\x06status\x18\x02 \x01(\x0b\x32\x0e.ProcessStatus\x12\x0f\n\x07\x64\x65leted\x18\x03 \x01(\x08\"\x14\n\x12ProcessListRequest\"\x91\x01\n\x13ProcessListResponse\x12\x36\n\tprocesses\x18\x01 \x03(\x0b\x32#.ProcessListResponse.ProcessesEntry\x1a\x42\n\x0eProcessesEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\x1f\n\x05value\x18\x02 \x01(\x0b\x32\x10.ProcessResponse:\x02\x38\x01\"\x1a\n\nLogRequest\x12\x0c\n\x04name\x18\x01 \x01(\t\"\x1b\n\x0bLogResponse\x12\x0c\n\x04line\x18\x02 \x01(\t\"M\n\x15ProcessReplaceRequest\x12\x1a\n\x04spec\x18\x01 \x01(\x0b\x32\x0c.ProcessSpec\x12\x18\n\x10terminate_signal\x18\x02 \x01(\t\"\xe4\x01\n\x0fVersionResponse\x12\x0f\n\x07version\x18\x01 \x01(\t\x12\x11\n\tgitCommit\x18\x02 \x01(\t\x12\x11\n\tbuildDate\x18\x03 \x01(\t\x12!\n\x19instanceManagerAPIVersion\x18\x04 \x01(\x03\x12$\n\x1cinstanceManagerAPIMinVersion\x18\x05 \x01(\x03\x12&\n\x1einstanceManagerProxyAPIVersion\x18\x06 \x01(\x03\x12)\n!instanceManagerProxyAPIMinVersion\x18\x07 \x01(\x03\x32\xe2\x03\n\x15ProcessManagerService\x12:\n\rProcessCreate\x12\x15.ProcessCreateRequest\x1a\x10.ProcessResponse\"\x00\x12:\n\rProcessDelete\x12\x15.ProcessDeleteRequest\x1a\x10.ProcessResponse\"\x00\x12\x34\n\nProcessGet\x12\x12.ProcessGetRequest\x1a\x10.ProcessResponse\"\x00\x12:\n\x0bProcessList\x12\x13.ProcessListRequest\x1a\x14.ProcessListResponse\"\x00\x12+\n\nProcessLog\x12\x0b.LogRequest\x1a\x0c.LogResponse\"\x00\x30\x01\x12<\n\x0cProcessWatch\x12\x16.google.protobuf.Empty\x1a\x10.ProcessResponse\"\x00\x30\x01\x12<\n\x0eProcessReplace\x12\x16.ProcessReplaceRequest\x1a\x10.ProcessResponse\"\x00\x12\x36\n\nVersionGet\x12\x16.google.protobuf.Empty\x1a\x10.VersionResponseb\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_empty__pb2.DESCRIPTOR,])




_PROCESSSPEC = _descriptor.Descriptor(
  name='ProcessSpec',
  full_name='ProcessSpec',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='ProcessSpec.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='binary', full_name='ProcessSpec.binary', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='args', full_name='ProcessSpec.args', index=2,
      number=3, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='port_count', full_name='ProcessSpec.port_count', index=3,
      number=4, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='port_args', full_name='ProcessSpec.port_args', index=4,
      number=5, type=9, cpp_type=9, label=3,
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
  serialized_start=44,
  serialized_end=140,
)


_PROCESSSTATUS = _descriptor.Descriptor(
  name='ProcessStatus',
  full_name='ProcessStatus',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='state', full_name='ProcessStatus.state', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='error_msg', full_name='ProcessStatus.error_msg', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='port_start', full_name='ProcessStatus.port_start', index=2,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='port_end', full_name='ProcessStatus.port_end', index=3,
      number=4, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='uuid', full_name='ProcessStatus.uuid', index=4,
      number=5, type=9, cpp_type=9, label=1,
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
  serialized_start=142,
  serialized_end=243,
)


_PROCESSCREATEREQUEST = _descriptor.Descriptor(
  name='ProcessCreateRequest',
  full_name='ProcessCreateRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='spec', full_name='ProcessCreateRequest.spec', index=0,
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
  serialized_start=245,
  serialized_end=295,
)


_PROCESSDELETEREQUEST = _descriptor.Descriptor(
  name='ProcessDeleteRequest',
  full_name='ProcessDeleteRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='ProcessDeleteRequest.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
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
  serialized_start=297,
  serialized_end=333,
)


_PROCESSGETREQUEST = _descriptor.Descriptor(
  name='ProcessGetRequest',
  full_name='ProcessGetRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='ProcessGetRequest.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
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
  serialized_start=335,
  serialized_end=368,
)


_PROCESSRESPONSE = _descriptor.Descriptor(
  name='ProcessResponse',
  full_name='ProcessResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='spec', full_name='ProcessResponse.spec', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='status', full_name='ProcessResponse.status', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='deleted', full_name='ProcessResponse.deleted', index=2,
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
  serialized_start=370,
  serialized_end=464,
)


_PROCESSLISTREQUEST = _descriptor.Descriptor(
  name='ProcessListRequest',
  full_name='ProcessListRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
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
  serialized_start=466,
  serialized_end=486,
)


_PROCESSLISTRESPONSE_PROCESSESENTRY = _descriptor.Descriptor(
  name='ProcessesEntry',
  full_name='ProcessListResponse.ProcessesEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='ProcessListResponse.ProcessesEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='ProcessListResponse.ProcessesEntry.value', index=1,
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
  serialized_start=568,
  serialized_end=634,
)

_PROCESSLISTRESPONSE = _descriptor.Descriptor(
  name='ProcessListResponse',
  full_name='ProcessListResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='processes', full_name='ProcessListResponse.processes', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_PROCESSLISTRESPONSE_PROCESSESENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=489,
  serialized_end=634,
)


_LOGREQUEST = _descriptor.Descriptor(
  name='LogRequest',
  full_name='LogRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='LogRequest.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
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
  serialized_start=636,
  serialized_end=662,
)


_LOGRESPONSE = _descriptor.Descriptor(
  name='LogResponse',
  full_name='LogResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='line', full_name='LogResponse.line', index=0,
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
  serialized_start=664,
  serialized_end=691,
)


_PROCESSREPLACEREQUEST = _descriptor.Descriptor(
  name='ProcessReplaceRequest',
  full_name='ProcessReplaceRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='spec', full_name='ProcessReplaceRequest.spec', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='terminate_signal', full_name='ProcessReplaceRequest.terminate_signal', index=1,
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
  serialized_start=693,
  serialized_end=770,
)


_VERSIONRESPONSE = _descriptor.Descriptor(
  name='VersionResponse',
  full_name='VersionResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='version', full_name='VersionResponse.version', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='gitCommit', full_name='VersionResponse.gitCommit', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='buildDate', full_name='VersionResponse.buildDate', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='instanceManagerAPIVersion', full_name='VersionResponse.instanceManagerAPIVersion', index=3,
      number=4, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='instanceManagerAPIMinVersion', full_name='VersionResponse.instanceManagerAPIMinVersion', index=4,
      number=5, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='instanceManagerProxyAPIVersion', full_name='VersionResponse.instanceManagerProxyAPIVersion', index=5,
      number=6, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='instanceManagerProxyAPIMinVersion', full_name='VersionResponse.instanceManagerProxyAPIMinVersion', index=6,
      number=7, type=3, cpp_type=2, label=1,
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
  serialized_start=773,
  serialized_end=1001,
)

_PROCESSCREATEREQUEST.fields_by_name['spec'].message_type = _PROCESSSPEC
_PROCESSRESPONSE.fields_by_name['spec'].message_type = _PROCESSSPEC
_PROCESSRESPONSE.fields_by_name['status'].message_type = _PROCESSSTATUS
_PROCESSLISTRESPONSE_PROCESSESENTRY.fields_by_name['value'].message_type = _PROCESSRESPONSE
_PROCESSLISTRESPONSE_PROCESSESENTRY.containing_type = _PROCESSLISTRESPONSE
_PROCESSLISTRESPONSE.fields_by_name['processes'].message_type = _PROCESSLISTRESPONSE_PROCESSESENTRY
_PROCESSREPLACEREQUEST.fields_by_name['spec'].message_type = _PROCESSSPEC
DESCRIPTOR.message_types_by_name['ProcessSpec'] = _PROCESSSPEC
DESCRIPTOR.message_types_by_name['ProcessStatus'] = _PROCESSSTATUS
DESCRIPTOR.message_types_by_name['ProcessCreateRequest'] = _PROCESSCREATEREQUEST
DESCRIPTOR.message_types_by_name['ProcessDeleteRequest'] = _PROCESSDELETEREQUEST
DESCRIPTOR.message_types_by_name['ProcessGetRequest'] = _PROCESSGETREQUEST
DESCRIPTOR.message_types_by_name['ProcessResponse'] = _PROCESSRESPONSE
DESCRIPTOR.message_types_by_name['ProcessListRequest'] = _PROCESSLISTREQUEST
DESCRIPTOR.message_types_by_name['ProcessListResponse'] = _PROCESSLISTRESPONSE
DESCRIPTOR.message_types_by_name['LogRequest'] = _LOGREQUEST
DESCRIPTOR.message_types_by_name['LogResponse'] = _LOGRESPONSE
DESCRIPTOR.message_types_by_name['ProcessReplaceRequest'] = _PROCESSREPLACEREQUEST
DESCRIPTOR.message_types_by_name['VersionResponse'] = _VERSIONRESPONSE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

ProcessSpec = _reflection.GeneratedProtocolMessageType('ProcessSpec', (_message.Message,), {
  'DESCRIPTOR' : _PROCESSSPEC,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:ProcessSpec)
  })
_sym_db.RegisterMessage(ProcessSpec)

ProcessStatus = _reflection.GeneratedProtocolMessageType('ProcessStatus', (_message.Message,), {
  'DESCRIPTOR' : _PROCESSSTATUS,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:ProcessStatus)
  })
_sym_db.RegisterMessage(ProcessStatus)

ProcessCreateRequest = _reflection.GeneratedProtocolMessageType('ProcessCreateRequest', (_message.Message,), {
  'DESCRIPTOR' : _PROCESSCREATEREQUEST,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:ProcessCreateRequest)
  })
_sym_db.RegisterMessage(ProcessCreateRequest)

ProcessDeleteRequest = _reflection.GeneratedProtocolMessageType('ProcessDeleteRequest', (_message.Message,), {
  'DESCRIPTOR' : _PROCESSDELETEREQUEST,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:ProcessDeleteRequest)
  })
_sym_db.RegisterMessage(ProcessDeleteRequest)

ProcessGetRequest = _reflection.GeneratedProtocolMessageType('ProcessGetRequest', (_message.Message,), {
  'DESCRIPTOR' : _PROCESSGETREQUEST,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:ProcessGetRequest)
  })
_sym_db.RegisterMessage(ProcessGetRequest)

ProcessResponse = _reflection.GeneratedProtocolMessageType('ProcessResponse', (_message.Message,), {
  'DESCRIPTOR' : _PROCESSRESPONSE,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:ProcessResponse)
  })
_sym_db.RegisterMessage(ProcessResponse)

ProcessListRequest = _reflection.GeneratedProtocolMessageType('ProcessListRequest', (_message.Message,), {
  'DESCRIPTOR' : _PROCESSLISTREQUEST,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:ProcessListRequest)
  })
_sym_db.RegisterMessage(ProcessListRequest)

ProcessListResponse = _reflection.GeneratedProtocolMessageType('ProcessListResponse', (_message.Message,), {

  'ProcessesEntry' : _reflection.GeneratedProtocolMessageType('ProcessesEntry', (_message.Message,), {
    'DESCRIPTOR' : _PROCESSLISTRESPONSE_PROCESSESENTRY,
    '__module__' : 'imrpc_pb2'
    # @@protoc_insertion_point(class_scope:ProcessListResponse.ProcessesEntry)
    })
  ,
  'DESCRIPTOR' : _PROCESSLISTRESPONSE,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:ProcessListResponse)
  })
_sym_db.RegisterMessage(ProcessListResponse)
_sym_db.RegisterMessage(ProcessListResponse.ProcessesEntry)

LogRequest = _reflection.GeneratedProtocolMessageType('LogRequest', (_message.Message,), {
  'DESCRIPTOR' : _LOGREQUEST,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:LogRequest)
  })
_sym_db.RegisterMessage(LogRequest)

LogResponse = _reflection.GeneratedProtocolMessageType('LogResponse', (_message.Message,), {
  'DESCRIPTOR' : _LOGRESPONSE,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:LogResponse)
  })
_sym_db.RegisterMessage(LogResponse)

ProcessReplaceRequest = _reflection.GeneratedProtocolMessageType('ProcessReplaceRequest', (_message.Message,), {
  'DESCRIPTOR' : _PROCESSREPLACEREQUEST,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:ProcessReplaceRequest)
  })
_sym_db.RegisterMessage(ProcessReplaceRequest)

VersionResponse = _reflection.GeneratedProtocolMessageType('VersionResponse', (_message.Message,), {
  'DESCRIPTOR' : _VERSIONRESPONSE,
  '__module__' : 'imrpc_pb2'
  # @@protoc_insertion_point(class_scope:VersionResponse)
  })
_sym_db.RegisterMessage(VersionResponse)


_PROCESSLISTRESPONSE_PROCESSESENTRY._options = None

_PROCESSMANAGERSERVICE = _descriptor.ServiceDescriptor(
  name='ProcessManagerService',
  full_name='ProcessManagerService',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=1004,
  serialized_end=1486,
  methods=[
  _descriptor.MethodDescriptor(
    name='ProcessCreate',
    full_name='ProcessManagerService.ProcessCreate',
    index=0,
    containing_service=None,
    input_type=_PROCESSCREATEREQUEST,
    output_type=_PROCESSRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ProcessDelete',
    full_name='ProcessManagerService.ProcessDelete',
    index=1,
    containing_service=None,
    input_type=_PROCESSDELETEREQUEST,
    output_type=_PROCESSRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ProcessGet',
    full_name='ProcessManagerService.ProcessGet',
    index=2,
    containing_service=None,
    input_type=_PROCESSGETREQUEST,
    output_type=_PROCESSRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ProcessList',
    full_name='ProcessManagerService.ProcessList',
    index=3,
    containing_service=None,
    input_type=_PROCESSLISTREQUEST,
    output_type=_PROCESSLISTRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ProcessLog',
    full_name='ProcessManagerService.ProcessLog',
    index=4,
    containing_service=None,
    input_type=_LOGREQUEST,
    output_type=_LOGRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ProcessWatch',
    full_name='ProcessManagerService.ProcessWatch',
    index=5,
    containing_service=None,
    input_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    output_type=_PROCESSRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='ProcessReplace',
    full_name='ProcessManagerService.ProcessReplace',
    index=6,
    containing_service=None,
    input_type=_PROCESSREPLACEREQUEST,
    output_type=_PROCESSRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='VersionGet',
    full_name='ProcessManagerService.VersionGet',
    index=7,
    containing_service=None,
    input_type=google_dot_protobuf_dot_empty__pb2._EMPTY,
    output_type=_VERSIONRESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_PROCESSMANAGERSERVICE)

DESCRIPTOR.services_by_name['ProcessManagerService'] = _PROCESSMANAGERSERVICE

# @@protoc_insertion_point(module_scope)
