// Code generated by protoc-gen-go. DO NOT EDIT.
// source: imcommon.proto

package imrpc

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type BackendStoreDriver int32

const (
	BackendStoreDriver_longhorn BackendStoreDriver = 0
	BackendStoreDriver_spdk     BackendStoreDriver = 1
)

var BackendStoreDriver_name = map[int32]string{
	0: "longhorn",
	1: "spdk",
}

var BackendStoreDriver_value = map[string]int32{
	"longhorn": 0,
	"spdk":     1,
}

func (x BackendStoreDriver) String() string {
	return proto.EnumName(BackendStoreDriver_name, int32(x))
}

func (BackendStoreDriver) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d185168b042b46f6, []int{0}
}

type VersionResponse struct {
	Version                                     string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	GitCommit                                   string   `protobuf:"bytes,2,opt,name=gitCommit,proto3" json:"gitCommit,omitempty"`
	BuildDate                                   string   `protobuf:"bytes,3,opt,name=buildDate,proto3" json:"buildDate,omitempty"`
	InstanceManagerAPIVersion                   int64    `protobuf:"varint,4,opt,name=instanceManagerAPIVersion,proto3" json:"instanceManagerAPIVersion,omitempty"`
	InstanceManagerAPIMinVersion                int64    `protobuf:"varint,5,opt,name=instanceManagerAPIMinVersion,proto3" json:"instanceManagerAPIMinVersion,omitempty"`
	InstanceManagerProxyAPIVersion              int64    `protobuf:"varint,6,opt,name=instanceManagerProxyAPIVersion,proto3" json:"instanceManagerProxyAPIVersion,omitempty"`
	InstanceManagerProxyAPIMinVersion           int64    `protobuf:"varint,7,opt,name=instanceManagerProxyAPIMinVersion,proto3" json:"instanceManagerProxyAPIMinVersion,omitempty"`
	InstanceManagerDiskServiceAPIVersion        int64    `protobuf:"varint,8,opt,name=instanceManagerDiskServiceAPIVersion,proto3" json:"instanceManagerDiskServiceAPIVersion,omitempty"`
	InstanceManagerDiskServiceAPIMinVersion     int64    `protobuf:"varint,9,opt,name=instanceManagerDiskServiceAPIMinVersion,proto3" json:"instanceManagerDiskServiceAPIMinVersion,omitempty"`
	InstanceManagerInstanceServiceAPIVersion    int64    `protobuf:"varint,10,opt,name=instanceManagerInstanceServiceAPIVersion,proto3" json:"instanceManagerInstanceServiceAPIVersion,omitempty"`
	InstanceManagerInstanceServiceAPIMinVersion int64    `protobuf:"varint,11,opt,name=instanceManagerInstanceServiceAPIMinVersion,proto3" json:"instanceManagerInstanceServiceAPIMinVersion,omitempty"`
	XXX_NoUnkeyedLiteral                        struct{} `json:"-"`
	XXX_unrecognized                            []byte   `json:"-"`
	XXX_sizecache                               int32    `json:"-"`
}

func (m *VersionResponse) Reset()         { *m = VersionResponse{} }
func (m *VersionResponse) String() string { return proto.CompactTextString(m) }
func (*VersionResponse) ProtoMessage()    {}
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d185168b042b46f6, []int{0}
}

func (m *VersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionResponse.Unmarshal(m, b)
}
func (m *VersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionResponse.Marshal(b, m, deterministic)
}
func (m *VersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionResponse.Merge(m, src)
}
func (m *VersionResponse) XXX_Size() int {
	return xxx_messageInfo_VersionResponse.Size(m)
}
func (m *VersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VersionResponse proto.InternalMessageInfo

func (m *VersionResponse) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *VersionResponse) GetGitCommit() string {
	if m != nil {
		return m.GitCommit
	}
	return ""
}

func (m *VersionResponse) GetBuildDate() string {
	if m != nil {
		return m.BuildDate
	}
	return ""
}

func (m *VersionResponse) GetInstanceManagerAPIVersion() int64 {
	if m != nil {
		return m.InstanceManagerAPIVersion
	}
	return 0
}

func (m *VersionResponse) GetInstanceManagerAPIMinVersion() int64 {
	if m != nil {
		return m.InstanceManagerAPIMinVersion
	}
	return 0
}

func (m *VersionResponse) GetInstanceManagerProxyAPIVersion() int64 {
	if m != nil {
		return m.InstanceManagerProxyAPIVersion
	}
	return 0
}

func (m *VersionResponse) GetInstanceManagerProxyAPIMinVersion() int64 {
	if m != nil {
		return m.InstanceManagerProxyAPIMinVersion
	}
	return 0
}

func (m *VersionResponse) GetInstanceManagerDiskServiceAPIVersion() int64 {
	if m != nil {
		return m.InstanceManagerDiskServiceAPIVersion
	}
	return 0
}

func (m *VersionResponse) GetInstanceManagerDiskServiceAPIMinVersion() int64 {
	if m != nil {
		return m.InstanceManagerDiskServiceAPIMinVersion
	}
	return 0
}

func (m *VersionResponse) GetInstanceManagerInstanceServiceAPIVersion() int64 {
	if m != nil {
		return m.InstanceManagerInstanceServiceAPIVersion
	}
	return 0
}

func (m *VersionResponse) GetInstanceManagerInstanceServiceAPIMinVersion() int64 {
	if m != nil {
		return m.InstanceManagerInstanceServiceAPIMinVersion
	}
	return 0
}

type LogResponse struct {
	Line                 string   `protobuf:"bytes,2,opt,name=line,proto3" json:"line,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogResponse) Reset()         { *m = LogResponse{} }
func (m *LogResponse) String() string { return proto.CompactTextString(m) }
func (*LogResponse) ProtoMessage()    {}
func (*LogResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d185168b042b46f6, []int{1}
}

func (m *LogResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogResponse.Unmarshal(m, b)
}
func (m *LogResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogResponse.Marshal(b, m, deterministic)
}
func (m *LogResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogResponse.Merge(m, src)
}
func (m *LogResponse) XXX_Size() int {
	return xxx_messageInfo_LogResponse.Size(m)
}
func (m *LogResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LogResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LogResponse proto.InternalMessageInfo

func (m *LogResponse) GetLine() string {
	if m != nil {
		return m.Line
	}
	return ""
}

func init() {
	proto.RegisterEnum("imrpc.BackendStoreDriver", BackendStoreDriver_name, BackendStoreDriver_value)
	proto.RegisterType((*VersionResponse)(nil), "imrpc.VersionResponse")
	proto.RegisterType((*LogResponse)(nil), "imrpc.LogResponse")
}

func init() { proto.RegisterFile("imcommon.proto", fileDescriptor_d185168b042b46f6) }

var fileDescriptor_d185168b042b46f6 = []byte{
	// 326 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xcd, 0x4a, 0xf3, 0x50,
	0x10, 0x86, 0xbf, 0x7c, 0xfd, 0x9f, 0x8a, 0x96, 0x59, 0x45, 0x28, 0xd2, 0x16, 0xc1, 0xa2, 0xe2,
	0xc6, 0xad, 0x1b, 0x6b, 0x11, 0x0a, 0x2d, 0x94, 0x14, 0x44, 0x5c, 0x99, 0xa6, 0x43, 0x1c, 0xda,
	0x9c, 0x09, 0x27, 0xc7, 0xa2, 0xd7, 0xec, 0x4d, 0x88, 0xa7, 0x3f, 0x86, 0x96, 0xb6, 0x71, 0x97,
	0x99, 0xf7, 0xc9, 0x93, 0x77, 0x31, 0x81, 0x63, 0x8e, 0x02, 0x89, 0x22, 0x51, 0x37, 0xb1, 0x16,
	0x23, 0x58, 0xe0, 0x48, 0xc7, 0x41, 0xeb, 0xab, 0x00, 0x27, 0x4f, 0xa4, 0x13, 0x16, 0xe5, 0x51,
	0x12, 0x8b, 0x4a, 0x08, 0x5d, 0x28, 0xcd, 0x17, 0x2b, 0xd7, 0x69, 0x38, 0xed, 0x8a, 0xb7, 0x1a,
	0xb1, 0x0e, 0x95, 0x90, 0xcd, 0x83, 0x44, 0x11, 0x1b, 0xf7, 0xbf, 0xcd, 0x7e, 0x17, 0x3f, 0xe9,
	0xf8, 0x9d, 0x67, 0x93, 0xae, 0x6f, 0xc8, 0xcd, 0x2d, 0xd2, 0xf5, 0x02, 0xef, 0xe0, 0x94, 0x55,
	0x62, 0x7c, 0x15, 0xd0, 0xc0, 0x57, 0x7e, 0x48, 0xfa, 0x7e, 0xd8, 0x5b, 0x7e, 0xda, 0xcd, 0x37,
	0x9c, 0x76, 0xce, 0xdb, 0x0d, 0x60, 0x07, 0xea, 0xdb, 0xe1, 0x80, 0xd5, 0x4a, 0x50, 0xb0, 0x82,
	0xbd, 0x0c, 0x3e, 0xc2, 0xd9, 0x46, 0x3e, 0xd4, 0xf2, 0xf1, 0x99, 0xaa, 0x51, 0xb4, 0x96, 0x03,
	0x14, 0xf6, 0xa1, 0xb9, 0x83, 0x48, 0x15, 0x2a, 0x59, 0xd5, 0x61, 0x10, 0x3d, 0x38, 0xdf, 0x80,
	0xba, 0x9c, 0x4c, 0x47, 0xa4, 0xe7, 0x1c, 0x50, 0xaa, 0x5b, 0xd9, 0x0a, 0x33, 0xb1, 0xf8, 0x0c,
	0x17, 0x7b, 0xb9, 0x54, 0xcf, 0x8a, 0xd5, 0x66, 0xc5, 0xf1, 0x05, 0xda, 0x1b, 0x68, 0x6f, 0x39,
	0x6e, 0x37, 0x06, 0xab, 0xce, 0xcc, 0xe3, 0x2b, 0x5c, 0x1d, 0x64, 0x53, 0xcd, 0xab, 0x56, 0xff,
	0x97, 0x57, 0x5a, 0x4d, 0xa8, 0xf6, 0x25, 0x5c, 0x1f, 0x3a, 0x42, 0x7e, 0xc6, 0x8a, 0x96, 0x97,
	0x6c, 0x9f, 0x2f, 0xaf, 0x01, 0x3b, 0x7e, 0x30, 0x25, 0x35, 0x19, 0x19, 0xd1, 0xd4, 0xd5, 0x3c,
	0x27, 0x8d, 0x47, 0x50, 0x9e, 0x89, 0x0a, 0xdf, 0x44, 0xab, 0xda, 0x3f, 0x2c, 0x43, 0x3e, 0x89,
	0x27, 0xd3, 0x9a, 0x33, 0x2e, 0xda, 0x9f, 0xe9, 0xf6, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x3d, 0xa1,
	0x91, 0xf8, 0x5e, 0x03, 0x00, 0x00,
}