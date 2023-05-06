// Code generated by protoc-gen-go. DO NOT EDIT.
// source: disk.proto

package imrpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type DiskType int32

const (
	DiskType_filesystem DiskType = 0
	DiskType_block      DiskType = 1
)

var DiskType_name = map[int32]string{
	0: "filesystem",
	1: "block",
}

var DiskType_value = map[string]int32{
	"filesystem": 0,
	"block":      1,
}

func (x DiskType) String() string {
	return proto.EnumName(DiskType_name, int32(x))
}

func (DiskType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f96b80c5532b4167, []int{0}
}

type Disk struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Uuid                 string   `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Path                 string   `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
	Type                 string   `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	TotalSize            int64    `protobuf:"varint,5,opt,name=total_size,json=totalSize,proto3" json:"total_size,omitempty"`
	FreeSize             int64    `protobuf:"varint,6,opt,name=free_size,json=freeSize,proto3" json:"free_size,omitempty"`
	TotalBlocks          int64    `protobuf:"varint,7,opt,name=total_blocks,json=totalBlocks,proto3" json:"total_blocks,omitempty"`
	FreeBlocks           int64    `protobuf:"varint,8,opt,name=free_blocks,json=freeBlocks,proto3" json:"free_blocks,omitempty"`
	BlockSize            int64    `protobuf:"varint,9,opt,name=block_size,json=blockSize,proto3" json:"block_size,omitempty"`
	ClusterSize          int64    `protobuf:"varint,10,opt,name=cluster_size,json=clusterSize,proto3" json:"cluster_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Disk) Reset()         { *m = Disk{} }
func (m *Disk) String() string { return proto.CompactTextString(m) }
func (*Disk) ProtoMessage()    {}
func (*Disk) Descriptor() ([]byte, []int) {
	return fileDescriptor_f96b80c5532b4167, []int{0}
}

func (m *Disk) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Disk.Unmarshal(m, b)
}
func (m *Disk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Disk.Marshal(b, m, deterministic)
}
func (m *Disk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Disk.Merge(m, src)
}
func (m *Disk) XXX_Size() int {
	return xxx_messageInfo_Disk.Size(m)
}
func (m *Disk) XXX_DiscardUnknown() {
	xxx_messageInfo_Disk.DiscardUnknown(m)
}

var xxx_messageInfo_Disk proto.InternalMessageInfo

func (m *Disk) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Disk) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Disk) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Disk) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Disk) GetTotalSize() int64 {
	if m != nil {
		return m.TotalSize
	}
	return 0
}

func (m *Disk) GetFreeSize() int64 {
	if m != nil {
		return m.FreeSize
	}
	return 0
}

func (m *Disk) GetTotalBlocks() int64 {
	if m != nil {
		return m.TotalBlocks
	}
	return 0
}

func (m *Disk) GetFreeBlocks() int64 {
	if m != nil {
		return m.FreeBlocks
	}
	return 0
}

func (m *Disk) GetBlockSize() int64 {
	if m != nil {
		return m.BlockSize
	}
	return 0
}

func (m *Disk) GetClusterSize() int64 {
	if m != nil {
		return m.ClusterSize
	}
	return 0
}

type DiskCreateRequest struct {
	DiskType             DiskType `protobuf:"varint,1,opt,name=disk_type,json=diskType,proto3,enum=imrpc.DiskType" json:"disk_type,omitempty"`
	DiskName             string   `protobuf:"bytes,2,opt,name=disk_name,json=diskName,proto3" json:"disk_name,omitempty"`
	DiskPath             string   `protobuf:"bytes,3,opt,name=disk_path,json=diskPath,proto3" json:"disk_path,omitempty"`
	BlockSize            int64    `protobuf:"varint,4,opt,name=block_size,json=blockSize,proto3" json:"block_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiskCreateRequest) Reset()         { *m = DiskCreateRequest{} }
func (m *DiskCreateRequest) String() string { return proto.CompactTextString(m) }
func (*DiskCreateRequest) ProtoMessage()    {}
func (*DiskCreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f96b80c5532b4167, []int{1}
}

func (m *DiskCreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiskCreateRequest.Unmarshal(m, b)
}
func (m *DiskCreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiskCreateRequest.Marshal(b, m, deterministic)
}
func (m *DiskCreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiskCreateRequest.Merge(m, src)
}
func (m *DiskCreateRequest) XXX_Size() int {
	return xxx_messageInfo_DiskCreateRequest.Size(m)
}
func (m *DiskCreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DiskCreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DiskCreateRequest proto.InternalMessageInfo

func (m *DiskCreateRequest) GetDiskType() DiskType {
	if m != nil {
		return m.DiskType
	}
	return DiskType_filesystem
}

func (m *DiskCreateRequest) GetDiskName() string {
	if m != nil {
		return m.DiskName
	}
	return ""
}

func (m *DiskCreateRequest) GetDiskPath() string {
	if m != nil {
		return m.DiskPath
	}
	return ""
}

func (m *DiskCreateRequest) GetBlockSize() int64 {
	if m != nil {
		return m.BlockSize
	}
	return 0
}

type DiskGetRequest struct {
	DiskType             DiskType `protobuf:"varint,1,opt,name=disk_type,json=diskType,proto3,enum=imrpc.DiskType" json:"disk_type,omitempty"`
	DiskName             string   `protobuf:"bytes,2,opt,name=disk_name,json=diskName,proto3" json:"disk_name,omitempty"`
	DiskPath             string   `protobuf:"bytes,3,opt,name=disk_path,json=diskPath,proto3" json:"disk_path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiskGetRequest) Reset()         { *m = DiskGetRequest{} }
func (m *DiskGetRequest) String() string { return proto.CompactTextString(m) }
func (*DiskGetRequest) ProtoMessage()    {}
func (*DiskGetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f96b80c5532b4167, []int{2}
}

func (m *DiskGetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiskGetRequest.Unmarshal(m, b)
}
func (m *DiskGetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiskGetRequest.Marshal(b, m, deterministic)
}
func (m *DiskGetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiskGetRequest.Merge(m, src)
}
func (m *DiskGetRequest) XXX_Size() int {
	return xxx_messageInfo_DiskGetRequest.Size(m)
}
func (m *DiskGetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DiskGetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DiskGetRequest proto.InternalMessageInfo

func (m *DiskGetRequest) GetDiskType() DiskType {
	if m != nil {
		return m.DiskType
	}
	return DiskType_filesystem
}

func (m *DiskGetRequest) GetDiskName() string {
	if m != nil {
		return m.DiskName
	}
	return ""
}

func (m *DiskGetRequest) GetDiskPath() string {
	if m != nil {
		return m.DiskPath
	}
	return ""
}

type DiskDeleteRequest struct {
	DiskName             string   `protobuf:"bytes,1,opt,name=disk_name,json=diskName,proto3" json:"disk_name,omitempty"`
	DiskUuid             string   `protobuf:"bytes,2,opt,name=disk_uuid,json=diskUuid,proto3" json:"disk_uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiskDeleteRequest) Reset()         { *m = DiskDeleteRequest{} }
func (m *DiskDeleteRequest) String() string { return proto.CompactTextString(m) }
func (*DiskDeleteRequest) ProtoMessage()    {}
func (*DiskDeleteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f96b80c5532b4167, []int{3}
}

func (m *DiskDeleteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiskDeleteRequest.Unmarshal(m, b)
}
func (m *DiskDeleteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiskDeleteRequest.Marshal(b, m, deterministic)
}
func (m *DiskDeleteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiskDeleteRequest.Merge(m, src)
}
func (m *DiskDeleteRequest) XXX_Size() int {
	return xxx_messageInfo_DiskDeleteRequest.Size(m)
}
func (m *DiskDeleteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DiskDeleteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DiskDeleteRequest proto.InternalMessageInfo

func (m *DiskDeleteRequest) GetDiskName() string {
	if m != nil {
		return m.DiskName
	}
	return ""
}

func (m *DiskDeleteRequest) GetDiskUuid() string {
	if m != nil {
		return m.DiskUuid
	}
	return ""
}

func init() {
	proto.RegisterEnum("imrpc.DiskType", DiskType_name, DiskType_value)
	proto.RegisterType((*Disk)(nil), "imrpc.Disk")
	proto.RegisterType((*DiskCreateRequest)(nil), "imrpc.DiskCreateRequest")
	proto.RegisterType((*DiskGetRequest)(nil), "imrpc.DiskGetRequest")
	proto.RegisterType((*DiskDeleteRequest)(nil), "imrpc.DiskDeleteRequest")
}

func init() { proto.RegisterFile("disk.proto", fileDescriptor_f96b80c5532b4167) }

var fileDescriptor_f96b80c5532b4167 = []byte{
	// 460 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x53, 0x5f, 0x6b, 0xd4, 0x40,
	0x10, 0x37, 0xd7, 0xbb, 0xf6, 0x32, 0x27, 0xb1, 0x2e, 0x58, 0x42, 0x4a, 0xb1, 0x1e, 0x08, 0x45,
	0x34, 0x85, 0xf6, 0x55, 0x7c, 0xd0, 0x8a, 0x4f, 0x8a, 0xa4, 0xea, 0x6b, 0xc9, 0x25, 0x73, 0x75,
	0xb9, 0xec, 0x6d, 0xcc, 0x6e, 0xc4, 0xeb, 0xe7, 0xf0, 0xc9, 0x2f, 0xe9, 0x57, 0x90, 0x99, 0xdd,
	0x23, 0x89, 0x70, 0xaf, 0x7d, 0x9b, 0xfc, 0xfe, 0xec, 0xcc, 0xfc, 0x36, 0x0b, 0x50, 0x4a, 0xb3,
	0x4a, 0xeb, 0x46, 0x5b, 0x2d, 0x26, 0x52, 0x35, 0x75, 0x91, 0x1c, 0xdf, 0x6a, 0x7d, 0x5b, 0xe1,
	0x39, 0x83, 0x8b, 0x76, 0x79, 0x8e, 0xaa, 0xb6, 0x1b, 0xa7, 0x49, 0x22, 0xa9, 0x0a, 0xad, 0x94,
	0x5e, 0xbb, 0xef, 0xf9, 0xef, 0x11, 0x8c, 0xaf, 0xa4, 0x59, 0x89, 0x08, 0x46, 0xb2, 0x8c, 0x83,
	0xd3, 0xe0, 0x2c, 0xcc, 0x46, 0xb2, 0x14, 0x02, 0xc6, 0x6d, 0x2b, 0xcb, 0x78, 0xc4, 0x08, 0xd7,
	0x84, 0xd5, 0xb9, 0xfd, 0x1e, 0xef, 0x39, 0x8c, 0x6a, 0xc2, 0xec, 0xa6, 0xc6, 0x78, 0xec, 0x30,
	0xaa, 0xc5, 0x09, 0x80, 0xd5, 0x36, 0xaf, 0x6e, 0x8c, 0xbc, 0xc3, 0x78, 0x72, 0x1a, 0x9c, 0xed,
	0x65, 0x21, 0x23, 0xd7, 0xf2, 0x0e, 0xc5, 0x31, 0x84, 0xcb, 0x06, 0xd1, 0xb1, 0xfb, 0xcc, 0x4e,
	0x09, 0x60, 0xf2, 0x19, 0x3c, 0x74, 0xde, 0x45, 0xa5, 0x8b, 0x95, 0x89, 0x0f, 0x98, 0x9f, 0x31,
	0xf6, 0x96, 0x21, 0xf1, 0x14, 0x66, 0xec, 0xf7, 0x8a, 0x29, 0x2b, 0x80, 0x20, 0x2f, 0x38, 0x01,
	0x60, 0xce, 0x75, 0x08, 0x5d, 0x7f, 0x46, 0xb6, 0x2d, 0x8a, 0xaa, 0x35, 0x16, 0x1b, 0x27, 0x00,
	0xd7, 0xc2, 0x63, 0x24, 0x99, 0xff, 0x09, 0xe0, 0x31, 0xc5, 0xf2, 0xae, 0xc1, 0xdc, 0x62, 0x86,
	0x3f, 0x5a, 0x34, 0x56, 0xbc, 0x84, 0x90, 0xe2, 0xbe, 0xe1, 0x85, 0x29, 0xaa, 0xe8, 0xe2, 0x51,
	0xca, 0xa1, 0xa7, 0x24, 0xfe, 0xb2, 0xa9, 0x31, 0x9b, 0x96, 0xbe, 0xa2, 0x35, 0x59, 0xbd, 0xce,
	0x15, 0xfa, 0x18, 0x99, 0xfc, 0x94, 0xab, 0x8e, 0xec, 0xe5, 0xc9, 0xe4, 0x67, 0xca, 0x74, 0x38,
	0xff, 0xf8, 0xbf, 0xf9, 0xe7, 0xbf, 0x20, 0xa2, 0x76, 0x1f, 0xd0, 0xde, 0xf3, 0x60, 0xf3, 0x8f,
	0x2e, 0x95, 0x2b, 0xac, 0xb0, 0x4b, 0x65, 0x70, 0x5c, 0xb0, 0xe3, 0xb8, 0xde, 0xbf, 0xc4, 0xe4,
	0xd7, 0x56, 0x96, 0x2f, 0x9e, 0xc3, 0x74, 0x3b, 0x9e, 0x88, 0x00, 0x96, 0xb2, 0x42, 0xb3, 0x31,
	0x16, 0xd5, 0xe1, 0x03, 0x11, 0xc2, 0x84, 0x37, 0x3e, 0x0c, 0x2e, 0xfe, 0x06, 0x30, 0x23, 0xdd,
	0x35, 0x36, 0x3f, 0x65, 0x81, 0xe2, 0x12, 0xa0, 0xbb, 0x1b, 0x11, 0xf7, 0x16, 0x1d, 0x5c, 0x57,
	0x32, 0xeb, 0x31, 0xe2, 0x8d, 0x33, 0xb9, 0xd1, 0x07, 0xa6, 0xc1, 0x36, 0xc9, 0x51, 0xea, 0x9e,
	0x4f, 0xba, 0x7d, 0x3e, 0xe9, 0x7b, 0x7a, 0x3e, 0xe2, 0x15, 0x1c, 0xf8, 0xd0, 0xc5, 0x93, 0x9e,
	0xb9, 0xbb, 0x84, 0x61, 0xbb, 0xd7, 0x00, 0xdf, 0xb0, 0x31, 0x52, 0xaf, 0xc9, 0xb1, 0xe3, 0xd0,
	0xe4, 0xc8, 0x5b, 0xbc, 0x34, 0x43, 0x53, 0xeb, 0xb5, 0xc1, 0xc5, 0x3e, 0xeb, 0x2e, 0xff, 0x05,
	0x00, 0x00, 0xff, 0xff, 0xd8, 0xb2, 0xb6, 0x78, 0xde, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DiskServiceClient is the client API for DiskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DiskServiceClient interface {
	DiskCreate(ctx context.Context, in *DiskCreateRequest, opts ...grpc.CallOption) (*Disk, error)
	DiskDelete(ctx context.Context, in *DiskDeleteRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DiskGet(ctx context.Context, in *DiskGetRequest, opts ...grpc.CallOption) (*Disk, error)
	VersionGet(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*VersionResponse, error)
}

type diskServiceClient struct {
	cc *grpc.ClientConn
}

func NewDiskServiceClient(cc *grpc.ClientConn) DiskServiceClient {
	return &diskServiceClient{cc}
}

func (c *diskServiceClient) DiskCreate(ctx context.Context, in *DiskCreateRequest, opts ...grpc.CallOption) (*Disk, error) {
	out := new(Disk)
	err := c.cc.Invoke(ctx, "/imrpc.DiskService/DiskCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *diskServiceClient) DiskDelete(ctx context.Context, in *DiskDeleteRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/imrpc.DiskService/DiskDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *diskServiceClient) DiskGet(ctx context.Context, in *DiskGetRequest, opts ...grpc.CallOption) (*Disk, error) {
	out := new(Disk)
	err := c.cc.Invoke(ctx, "/imrpc.DiskService/DiskGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *diskServiceClient) VersionGet(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/imrpc.DiskService/VersionGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiskServiceServer is the server API for DiskService service.
type DiskServiceServer interface {
	DiskCreate(context.Context, *DiskCreateRequest) (*Disk, error)
	DiskDelete(context.Context, *DiskDeleteRequest) (*empty.Empty, error)
	DiskGet(context.Context, *DiskGetRequest) (*Disk, error)
	VersionGet(context.Context, *empty.Empty) (*VersionResponse, error)
}

// UnimplementedDiskServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDiskServiceServer struct {
}

func (*UnimplementedDiskServiceServer) DiskCreate(ctx context.Context, req *DiskCreateRequest) (*Disk, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiskCreate not implemented")
}
func (*UnimplementedDiskServiceServer) DiskDelete(ctx context.Context, req *DiskDeleteRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiskDelete not implemented")
}
func (*UnimplementedDiskServiceServer) DiskGet(ctx context.Context, req *DiskGetRequest) (*Disk, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiskGet not implemented")
}
func (*UnimplementedDiskServiceServer) VersionGet(ctx context.Context, req *empty.Empty) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VersionGet not implemented")
}

func RegisterDiskServiceServer(s *grpc.Server, srv DiskServiceServer) {
	s.RegisterService(&_DiskService_serviceDesc, srv)
}

func _DiskService_DiskCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiskCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiskServiceServer).DiskCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imrpc.DiskService/DiskCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiskServiceServer).DiskCreate(ctx, req.(*DiskCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiskService_DiskDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiskDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiskServiceServer).DiskDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imrpc.DiskService/DiskDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiskServiceServer).DiskDelete(ctx, req.(*DiskDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiskService_DiskGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiskGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiskServiceServer).DiskGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imrpc.DiskService/DiskGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiskServiceServer).DiskGet(ctx, req.(*DiskGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DiskService_VersionGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiskServiceServer).VersionGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/imrpc.DiskService/VersionGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiskServiceServer).VersionGet(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _DiskService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "imrpc.DiskService",
	HandlerType: (*DiskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DiskCreate",
			Handler:    _DiskService_DiskCreate_Handler,
		},
		{
			MethodName: "DiskDelete",
			Handler:    _DiskService_DiskDelete_Handler,
		},
		{
			MethodName: "DiskGet",
			Handler:    _DiskService_DiskGet_Handler,
		},
		{
			MethodName: "VersionGet",
			Handler:    _DiskService_VersionGet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "disk.proto",
}