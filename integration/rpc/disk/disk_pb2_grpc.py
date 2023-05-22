# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import disk_pb2 as disk__pb2
from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
import imcommon_pb2 as imcommon__pb2


class DiskServiceStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.DiskCreate = channel.unary_unary(
        '/imrpc.DiskService/DiskCreate',
        request_serializer=disk__pb2.DiskCreateRequest.SerializeToString,
        response_deserializer=disk__pb2.Disk.FromString,
        )
    self.DiskDelete = channel.unary_unary(
        '/imrpc.DiskService/DiskDelete',
        request_serializer=disk__pb2.DiskDeleteRequest.SerializeToString,
        response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
        )
    self.DiskGet = channel.unary_unary(
        '/imrpc.DiskService/DiskGet',
        request_serializer=disk__pb2.DiskGetRequest.SerializeToString,
        response_deserializer=disk__pb2.Disk.FromString,
        )
    self.DiskReplicaInstanceList = channel.unary_unary(
        '/imrpc.DiskService/DiskReplicaInstanceList',
        request_serializer=disk__pb2.DiskReplicaInstanceListRequest.SerializeToString,
        response_deserializer=disk__pb2.DiskReplicaInstanceListResponse.FromString,
        )
    self.DiskReplicaInstanceDelete = channel.unary_unary(
        '/imrpc.DiskService/DiskReplicaInstanceDelete',
        request_serializer=disk__pb2.DiskReplicaInstanceDeleteRequest.SerializeToString,
        response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
        )
    self.VersionGet = channel.unary_unary(
        '/imrpc.DiskService/VersionGet',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=imcommon__pb2.VersionResponse.FromString,
        )


class DiskServiceServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def DiskCreate(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def DiskDelete(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def DiskGet(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def DiskReplicaInstanceList(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def DiskReplicaInstanceDelete(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def VersionGet(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_DiskServiceServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'DiskCreate': grpc.unary_unary_rpc_method_handler(
          servicer.DiskCreate,
          request_deserializer=disk__pb2.DiskCreateRequest.FromString,
          response_serializer=disk__pb2.Disk.SerializeToString,
      ),
      'DiskDelete': grpc.unary_unary_rpc_method_handler(
          servicer.DiskDelete,
          request_deserializer=disk__pb2.DiskDeleteRequest.FromString,
          response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
      ),
      'DiskGet': grpc.unary_unary_rpc_method_handler(
          servicer.DiskGet,
          request_deserializer=disk__pb2.DiskGetRequest.FromString,
          response_serializer=disk__pb2.Disk.SerializeToString,
      ),
      'DiskReplicaInstanceList': grpc.unary_unary_rpc_method_handler(
          servicer.DiskReplicaInstanceList,
          request_deserializer=disk__pb2.DiskReplicaInstanceListRequest.FromString,
          response_serializer=disk__pb2.DiskReplicaInstanceListResponse.SerializeToString,
      ),
      'DiskReplicaInstanceDelete': grpc.unary_unary_rpc_method_handler(
          servicer.DiskReplicaInstanceDelete,
          request_deserializer=disk__pb2.DiskReplicaInstanceDeleteRequest.FromString,
          response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
      ),
      'VersionGet': grpc.unary_unary_rpc_method_handler(
          servicer.VersionGet,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=imcommon__pb2.VersionResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'imrpc.DiskService', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
