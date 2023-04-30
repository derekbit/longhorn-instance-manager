#!/bin/bash

set -e

# check and download dependency for gRPC code generate
if [ ! -e ./proto/vendor/protobuf/src/google/protobuf ]; then
    rm -rf ./proto/vendor/protobuf/src/google/protobuf
    DIR="./proto/vendor/protobuf/src/google/protobuf"
    mkdir -p $DIR
    wget https://raw.githubusercontent.com/protocolbuffers/protobuf/v3.9.0/src/google/protobuf/empty.proto -P $DIR
fi
DIR="./proto/vendor/github.com/longhorn/longhorn-engine/proto/ptypes"
if [ ! -e ${DIR} ]; then
    rm -rf ${DIR}
    mkdir -p $DIR
    wget https://raw.githubusercontent.com/longhorn/longhorn-engine/master/proto/ptypes/common.proto -P $DIR
    wget https://raw.githubusercontent.com/longhorn/longhorn-engine/master/proto/ptypes/controller.proto -P $DIR
    wget https://raw.githubusercontent.com/longhorn/longhorn-engine/master/proto/ptypes/syncagent.proto -P $DIR
fi

# common
python3 -m grpc_tools.protoc -I pkg/imrpc -I proto/vendor/protobuf/src/ --python_out=integration/rpc/common --grpc_python_out=integration/rpc/common pkg/imrpc/imcommon.proto
protoc -I pkg/imrpc/ -I proto/vendor/protobuf/src/ pkg/imrpc/imcommon.proto --go_out=plugins=grpc:pkg/imrpc

# instance manager (process manager)
python3 -m grpc_tools.protoc -I pkg/imrpc -I proto/vendor/protobuf/src/ --python_out=integration/rpc/instance_manager --grpc_python_out=integration/rpc/instance_manager pkg/imrpc/imrpc.proto
protoc -I pkg/imrpc/ -I proto/vendor/protobuf/src/ pkg/imrpc/imrpc.proto --go_out=plugins=grpc:pkg/imrpc

# proxy
python3 -m grpc_tools.protoc -I pkg/imrpc -I proto/vendor/ -I proto/vendor/protobuf/src/ --python_out=integration/rpc/proxy --grpc_python_out=integration/rpc/proxy pkg/imrpc/proxy.proto
protoc -I pkg/imrpc/ -I proto/vendor/ -I proto/vendor/protobuf/src/ pkg/imrpc/proxy.proto --go_out=plugins=grpc:pkg/imrpc/

# disk
python3 -m grpc_tools.protoc -I pkg/imrpc -I proto/vendor/ -I proto/vendor/protobuf/src/ --python_out=integration/rpc/disk --grpc_python_out=integration/rpc/disk pkg/imrpc/disk.proto
protoc -I pkg/imrpc/ -I proto/vendor/ -I proto/vendor/protobuf/src/ pkg/imrpc/disk.proto --go_out=plugins=grpc:pkg/imrpc/

# instance
python3 -m grpc_tools.protoc -I pkg/imrpc -I proto/vendor/ -I proto/vendor/protobuf/src/ --python_out=integration/rpc/instance --grpc_python_out=integration/rpc/instance pkg/imrpc/instance.proto
protoc -I pkg/imrpc/ -I proto/vendor/ -I proto/vendor/protobuf/src/ pkg/imrpc/instance.proto --go_out=plugins=grpc:pkg/imrpc/

# spdk
python3 -m grpc_tools.protoc -I pkg/imrpc -I proto/vendor/ -I proto/vendor/protobuf/src/ --python_out=integration/rpc/spdk --grpc_python_out=integration/rpc/spdk pkg/imrpc/spdk.proto
protoc -I pkg/imrpc/ -I proto/vendor/ -I proto/vendor/protobuf/src/ pkg/imrpc/spdk.proto --go_out=plugins=grpc:pkg/imrpc/
