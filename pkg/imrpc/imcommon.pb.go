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
	Version                                 string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	GitCommit                               string   `protobuf:"bytes,2,opt,name=gitCommit,proto3" json:"gitCommit,omitempty"`
	BuildDate                               string   `protobuf:"bytes,3,opt,name=buildDate,proto3" json:"buildDate,omitempty"`
	InstanceManagerAPIVersion               int64    `protobuf:"varint,4,opt,name=instanceManagerAPIVersion,proto3" json:"instanceManagerAPIVersion,omitempty"`
	InstanceManagerAPIMinVersion            int64    `protobuf:"varint,5,opt,name=instanceManagerAPIMinVersion,proto3" json:"instanceManagerAPIMinVersion,omitempty"`
	InstanceManagerProxyAPIVersion          int64    `protobuf:"varint,6,opt,name=instanceManagerProxyAPIVersion,proto3" json:"instanceManagerProxyAPIVersion,omitempty"`
	InstanceManagerProxyAPIMinVersion       int64    `protobuf:"varint,7,opt,name=instanceManagerProxyAPIMinVersion,proto3" json:"instanceManagerProxyAPIMinVersion,omitempty"`
	InstanceManagerDiskServiceAPIVersion    int64    `protobuf:"varint,8,opt,name=instanceManagerDiskServiceAPIVersion,proto3" json:"instanceManagerDiskServiceAPIVersion,omitempty"`
	InstanceManagerDiskServiceAPIMinVersion int64    `protobuf:"varint,9,opt,name=instanceManagerDiskServiceAPIMinVersion,proto3" json:"instanceManagerDiskServiceAPIMinVersion,omitempty"`
	XXX_NoUnkeyedLiteral                    struct{} `json:"-"`
	XXX_unrecognized                        []byte   `json:"-"`
	XXX_sizecache                           int32    `json:"-"`
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
	// 298 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x86, 0x8d, 0x4d, 0xbf, 0x46, 0xd1, 0x32, 0xa7, 0x08, 0x45, 0xda, 0x22, 0x58, 0x44, 0xbc,
	0x78, 0xf5, 0x62, 0x2d, 0x82, 0xd0, 0x42, 0x49, 0x41, 0xbc, 0xa6, 0xe9, 0x10, 0x87, 0x36, 0x3b,
	0x61, 0x77, 0x0d, 0xfa, 0xa3, 0xfc, 0x8f, 0xe2, 0xf6, 0xc3, 0x50, 0xe9, 0xc7, 0x2d, 0x99, 0xf7,
	0xe1, 0xd9, 0xf7, 0xf0, 0xc2, 0x19, 0xa7, 0xb1, 0xa4, 0xa9, 0xa8, 0xbb, 0x4c, 0x8b, 0x15, 0x2c,
	0x73, 0xaa, 0xb3, 0xb8, 0xf3, 0xed, 0xc3, 0xf9, 0x2b, 0x69, 0xc3, 0xa2, 0x42, 0x32, 0x99, 0x28,
	0x43, 0x18, 0x40, 0x35, 0x5f, 0x9c, 0x02, 0xaf, 0xe5, 0x75, 0xeb, 0xe1, 0xea, 0x17, 0x9b, 0x50,
	0x4f, 0xd8, 0x3e, 0x49, 0x9a, 0xb2, 0x0d, 0x8e, 0x5d, 0xf6, 0x77, 0xf8, 0x4d, 0x27, 0x1f, 0x3c,
	0x9f, 0xf6, 0x23, 0x4b, 0x41, 0x69, 0x91, 0xae, 0x0f, 0xf8, 0x00, 0x17, 0xac, 0x8c, 0x8d, 0x54,
	0x4c, 0xc3, 0x48, 0x45, 0x09, 0xe9, 0xc7, 0xd1, 0xcb, 0xf2, 0xe9, 0xc0, 0x6f, 0x79, 0xdd, 0x52,
	0xb8, 0x1d, 0xc0, 0x1e, 0x34, 0xff, 0x87, 0x43, 0x56, 0x2b, 0x41, 0xd9, 0x09, 0x76, 0x32, 0xf8,
	0x0c, 0x97, 0x1b, 0xf9, 0x48, 0xcb, 0xe7, 0x57, 0xa1, 0x46, 0xc5, 0x59, 0xf6, 0x50, 0x38, 0x80,
	0xf6, 0x16, 0xa2, 0x50, 0xa8, 0xea, 0x54, 0xfb, 0x41, 0x0c, 0xe1, 0x6a, 0x03, 0xea, 0xb3, 0x99,
	0x8d, 0x49, 0xe7, 0x1c, 0x53, 0xa1, 0x5b, 0xcd, 0x09, 0x0f, 0x62, 0xf1, 0x0d, 0xae, 0x77, 0x72,
	0x85, 0x9e, 0x75, 0xa7, 0x3d, 0x14, 0xef, 0xb4, 0xe1, 0x64, 0x20, 0xc9, 0x7a, 0x2a, 0x08, 0xfe,
	0x9c, 0x15, 0x2d, 0xb7, 0xe0, 0xbe, 0x6f, 0x6e, 0x01, 0x7b, 0x51, 0x3c, 0x23, 0x35, 0x1d, 0x5b,
	0xd1, 0xd4, 0xd7, 0x9c, 0x93, 0xc6, 0x53, 0xa8, 0xcd, 0x45, 0x25, 0xef, 0xa2, 0x55, 0xe3, 0x08,
	0x6b, 0xe0, 0x9b, 0x6c, 0x3a, 0x6b, 0x78, 0x93, 0x8a, 0x9b, 0xe3, 0xfd, 0x4f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xb8, 0x41, 0x5a, 0x8e, 0xa0, 0x02, 0x00, 0x00,
}