# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
import imrpc_pb2 as imrpc__pb2
import spdk_pb2 as spdk__pb2


class SPDKServiceStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.ReplicaCreate = channel.unary_unary(
        '/imrpc.SPDKService/ReplicaCreate',
        request_serializer=spdk__pb2.ReplicaCreateRequest.SerializeToString,
        response_deserializer=spdk__pb2.Replica.FromString,
        )
    self.ReplicaDelete = channel.unary_unary(
        '/imrpc.SPDKService/ReplicaDelete',
        request_serializer=spdk__pb2.ReplicaDeleteRequest.SerializeToString,
        response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
        )
    self.ReplicaGet = channel.unary_unary(
        '/imrpc.SPDKService/ReplicaGet',
        request_serializer=spdk__pb2.ReplicaGetRequest.SerializeToString,
        response_deserializer=spdk__pb2.Replica.FromString,
        )
    self.ReplicaList = channel.unary_unary(
        '/imrpc.SPDKService/ReplicaList',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=spdk__pb2.ReplicaListResponse.FromString,
        )
    self.ReplicaWatch = channel.unary_stream(
        '/imrpc.SPDKService/ReplicaWatch',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=spdk__pb2.Replica.FromString,
        )
    self.EngineCreate = channel.unary_unary(
        '/imrpc.SPDKService/EngineCreate',
        request_serializer=spdk__pb2.EngineCreateRequest.SerializeToString,
        response_deserializer=spdk__pb2.Engine.FromString,
        )
    self.EngineDelete = channel.unary_unary(
        '/imrpc.SPDKService/EngineDelete',
        request_serializer=spdk__pb2.EngineDeleteRequest.SerializeToString,
        response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
        )
    self.EngineGet = channel.unary_unary(
        '/imrpc.SPDKService/EngineGet',
        request_serializer=spdk__pb2.EngineGetRequest.SerializeToString,
        response_deserializer=spdk__pb2.Engine.FromString,
        )
    self.EngineList = channel.unary_unary(
        '/imrpc.SPDKService/EngineList',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=spdk__pb2.EngineListResponse.FromString,
        )
    self.EngineWatch = channel.unary_stream(
        '/imrpc.SPDKService/EngineWatch',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=spdk__pb2.Engine.FromString,
        )
    self.VersionGet = channel.unary_unary(
        '/imrpc.SPDKService/VersionGet',
        request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
        response_deserializer=imrpc__pb2.VersionResponse.FromString,
        )


class SPDKServiceServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def ReplicaCreate(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def ReplicaDelete(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def ReplicaGet(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def ReplicaList(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def ReplicaWatch(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def EngineCreate(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def EngineDelete(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def EngineGet(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def EngineList(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def EngineWatch(self, request, context):
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


def add_SPDKServiceServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'ReplicaCreate': grpc.unary_unary_rpc_method_handler(
          servicer.ReplicaCreate,
          request_deserializer=spdk__pb2.ReplicaCreateRequest.FromString,
          response_serializer=spdk__pb2.Replica.SerializeToString,
      ),
      'ReplicaDelete': grpc.unary_unary_rpc_method_handler(
          servicer.ReplicaDelete,
          request_deserializer=spdk__pb2.ReplicaDeleteRequest.FromString,
          response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
      ),
      'ReplicaGet': grpc.unary_unary_rpc_method_handler(
          servicer.ReplicaGet,
          request_deserializer=spdk__pb2.ReplicaGetRequest.FromString,
          response_serializer=spdk__pb2.Replica.SerializeToString,
      ),
      'ReplicaList': grpc.unary_unary_rpc_method_handler(
          servicer.ReplicaList,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=spdk__pb2.ReplicaListResponse.SerializeToString,
      ),
      'ReplicaWatch': grpc.unary_stream_rpc_method_handler(
          servicer.ReplicaWatch,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=spdk__pb2.Replica.SerializeToString,
      ),
      'EngineCreate': grpc.unary_unary_rpc_method_handler(
          servicer.EngineCreate,
          request_deserializer=spdk__pb2.EngineCreateRequest.FromString,
          response_serializer=spdk__pb2.Engine.SerializeToString,
      ),
      'EngineDelete': grpc.unary_unary_rpc_method_handler(
          servicer.EngineDelete,
          request_deserializer=spdk__pb2.EngineDeleteRequest.FromString,
          response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
      ),
      'EngineGet': grpc.unary_unary_rpc_method_handler(
          servicer.EngineGet,
          request_deserializer=spdk__pb2.EngineGetRequest.FromString,
          response_serializer=spdk__pb2.Engine.SerializeToString,
      ),
      'EngineList': grpc.unary_unary_rpc_method_handler(
          servicer.EngineList,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=spdk__pb2.EngineListResponse.SerializeToString,
      ),
      'EngineWatch': grpc.unary_stream_rpc_method_handler(
          servicer.EngineWatch,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=spdk__pb2.Engine.SerializeToString,
      ),
      'VersionGet': grpc.unary_unary_rpc_method_handler(
          servicer.VersionGet,
          request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
          response_serializer=imrpc__pb2.VersionResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'imrpc.SPDKService', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
