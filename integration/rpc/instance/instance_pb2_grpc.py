# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
import imrpc_pb2 as imrpc__pb2
import instance_pb2 as instance__pb2


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
        request_serializer=instance__pb2.InstanceCreateRequest.SerializeToString,
        response_deserializer=instance__pb2.InstanceResponse.FromString,
        )
    self.InstanceDelete = channel.unary_unary(
        '/imrpc.InstanceService/InstanceDelete',
        request_serializer=instance__pb2.InstanceDeleteRequest.SerializeToString,
        response_deserializer=instance__pb2.InstanceResponse.FromString,
        )
    self.InstanceGet = channel.unary_unary(
        '/imrpc.InstanceService/InstanceGet',
        request_serializer=instance__pb2.InstanceGetRequest.SerializeToString,
        response_deserializer=instance__pb2.InstanceResponse.FromString,
        )
    self.InstanceList = channel.unary_unary(
        '/imrpc.InstanceService/InstanceList',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=instance__pb2.InstanceListResponse.FromString,
        )
    self.InstanceLog = channel.unary_stream(
        '/imrpc.InstanceService/InstanceLog',
        request_serializer=instance__pb2.InstanceLogRequest.SerializeToString,
        response_deserializer=imrpc__pb2.LogResponse.FromString,
        )
    self.InstanceWatch = channel.unary_stream(
        '/imrpc.InstanceService/InstanceWatch',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=instance__pb2.InstanceResponse.FromString,
        )
    self.InstanceReplace = channel.unary_unary(
        '/imrpc.InstanceService/InstanceReplace',
        request_serializer=instance__pb2.InstanceReplaceRequest.SerializeToString,
        response_deserializer=instance__pb2.InstanceResponse.FromString,
        )
    self.VersionGet = channel.unary_unary(
        '/imrpc.InstanceService/VersionGet',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=imrpc__pb2.VersionResponse.FromString,
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
          request_deserializer=instance__pb2.InstanceCreateRequest.FromString,
          response_serializer=instance__pb2.InstanceResponse.SerializeToString,
      ),
      'InstanceDelete': grpc.unary_unary_rpc_method_handler(
          servicer.InstanceDelete,
          request_deserializer=instance__pb2.InstanceDeleteRequest.FromString,
          response_serializer=instance__pb2.InstanceResponse.SerializeToString,
      ),
      'InstanceGet': grpc.unary_unary_rpc_method_handler(
          servicer.InstanceGet,
          request_deserializer=instance__pb2.InstanceGetRequest.FromString,
          response_serializer=instance__pb2.InstanceResponse.SerializeToString,
      ),
      'InstanceList': grpc.unary_unary_rpc_method_handler(
          servicer.InstanceList,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=instance__pb2.InstanceListResponse.SerializeToString,
      ),
      'InstanceLog': grpc.unary_stream_rpc_method_handler(
          servicer.InstanceLog,
          request_deserializer=instance__pb2.InstanceLogRequest.FromString,
          response_serializer=imrpc__pb2.LogResponse.SerializeToString,
      ),
      'InstanceWatch': grpc.unary_stream_rpc_method_handler(
          servicer.InstanceWatch,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=instance__pb2.InstanceResponse.SerializeToString,
      ),
      'InstanceReplace': grpc.unary_unary_rpc_method_handler(
          servicer.InstanceReplace,
          request_deserializer=instance__pb2.InstanceReplaceRequest.FromString,
          response_serializer=instance__pb2.InstanceResponse.SerializeToString,
      ),
      'VersionGet': grpc.unary_unary_rpc_method_handler(
          servicer.VersionGet,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=imrpc__pb2.VersionResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'imrpc.InstanceService', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
