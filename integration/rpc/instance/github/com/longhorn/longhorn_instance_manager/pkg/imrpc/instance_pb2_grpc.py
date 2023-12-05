# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

from github.com.longhorn.longhorn_instance_manager.pkg.imrpc import imrpc_pb2 as github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_imrpc__pb2
from github.com.longhorn.longhorn_instance_manager.pkg.imrpc import instance_pb2 as github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2
from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


class InstanceServiceStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.InstanceCreate = channel.unary_unary(
        '/imrpc.InstanceService/InstanceCreate',
        request_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceCreateRequest.SerializeToString,
        response_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceResponse.FromString,
        )
    self.InstanceDelete = channel.unary_unary(
        '/imrpc.InstanceService/InstanceDelete',
        request_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceDeleteRequest.SerializeToString,
        response_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceResponse.FromString,
        )
    self.InstanceGet = channel.unary_unary(
        '/imrpc.InstanceService/InstanceGet',
        request_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceGetRequest.SerializeToString,
        response_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceResponse.FromString,
        )
    self.InstanceList = channel.unary_unary(
        '/imrpc.InstanceService/InstanceList',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceListResponse.FromString,
        )
    self.InstanceLog = channel.unary_stream(
        '/imrpc.InstanceService/InstanceLog',
        request_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceLogRequest.SerializeToString,
        response_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_imrpc__pb2.LogResponse.FromString,
        )
    self.InstanceWatch = channel.unary_stream(
        '/imrpc.InstanceService/InstanceWatch',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
        )
    self.InstanceReplace = channel.unary_unary(
        '/imrpc.InstanceService/InstanceReplace',
        request_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceReplaceRequest.SerializeToString,
        response_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceResponse.FromString,
        )
    self.InstanceSuspend = channel.unary_unary(
        '/imrpc.InstanceService/InstanceSuspend',
        request_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceSuspendRequest.SerializeToString,
        response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
        )
    self.VersionGet = channel.unary_unary(
        '/imrpc.InstanceService/VersionGet',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_imrpc__pb2.VersionResponse.FromString,
        )


class InstanceServiceServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def InstanceCreate(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def InstanceDelete(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def InstanceGet(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def InstanceList(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def InstanceLog(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def InstanceWatch(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def InstanceReplace(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def InstanceSuspend(self, request, context):
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


def add_InstanceServiceServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'InstanceCreate': grpc.unary_unary_rpc_method_handler(
          servicer.InstanceCreate,
          request_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceCreateRequest.FromString,
          response_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceResponse.SerializeToString,
      ),
      'InstanceDelete': grpc.unary_unary_rpc_method_handler(
          servicer.InstanceDelete,
          request_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceDeleteRequest.FromString,
          response_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceResponse.SerializeToString,
      ),
      'InstanceGet': grpc.unary_unary_rpc_method_handler(
          servicer.InstanceGet,
          request_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceGetRequest.FromString,
          response_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceResponse.SerializeToString,
      ),
      'InstanceList': grpc.unary_unary_rpc_method_handler(
          servicer.InstanceList,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceListResponse.SerializeToString,
      ),
      'InstanceLog': grpc.unary_stream_rpc_method_handler(
          servicer.InstanceLog,
          request_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceLogRequest.FromString,
          response_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_imrpc__pb2.LogResponse.SerializeToString,
      ),
      'InstanceWatch': grpc.unary_stream_rpc_method_handler(
          servicer.InstanceWatch,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
      ),
      'InstanceReplace': grpc.unary_unary_rpc_method_handler(
          servicer.InstanceReplace,
          request_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceReplaceRequest.FromString,
          response_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceResponse.SerializeToString,
      ),
      'InstanceSuspend': grpc.unary_unary_rpc_method_handler(
          servicer.InstanceSuspend,
          request_deserializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_instance__pb2.InstanceSuspendRequest.FromString,
          response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
      ),
      'VersionGet': grpc.unary_unary_rpc_method_handler(
          servicer.VersionGet,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=github_dot_com_dot_longhorn_dot_longhorn__instance__manager_dot_pkg_dot_imrpc_dot_imrpc__pb2.VersionResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'imrpc.InstanceService', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
