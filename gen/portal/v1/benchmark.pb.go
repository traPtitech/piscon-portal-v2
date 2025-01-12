// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: portal/v1/benchmark.proto

package portalv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BenchmarkResult int32

const (
	BenchmarkResult_BENCHMARK_RESULT_UNSPECIFIED BenchmarkResult = 0
	BenchmarkResult_BENCHMARK_RESULT_PASSED      BenchmarkResult = 1
	BenchmarkResult_BENCHMARK_RESULT_FAILED      BenchmarkResult = 2
	BenchmarkResult_BENCHMARK_RESULT_ERROR       BenchmarkResult = 3
)

// Enum value maps for BenchmarkResult.
var (
	BenchmarkResult_name = map[int32]string{
		0: "BENCHMARK_RESULT_UNSPECIFIED",
		1: "BENCHMARK_RESULT_PASSED",
		2: "BENCHMARK_RESULT_FAILED",
		3: "BENCHMARK_RESULT_ERROR",
	}
	BenchmarkResult_value = map[string]int32{
		"BENCHMARK_RESULT_UNSPECIFIED": 0,
		"BENCHMARK_RESULT_PASSED":      1,
		"BENCHMARK_RESULT_FAILED":      2,
		"BENCHMARK_RESULT_ERROR":       3,
	}
)

func (x BenchmarkResult) Enum() *BenchmarkResult {
	p := new(BenchmarkResult)
	*p = x
	return p
}

func (x BenchmarkResult) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BenchmarkResult) Descriptor() protoreflect.EnumDescriptor {
	return file_portal_v1_benchmark_proto_enumTypes[0].Descriptor()
}

func (BenchmarkResult) Type() protoreflect.EnumType {
	return &file_portal_v1_benchmark_proto_enumTypes[0]
}

func (x BenchmarkResult) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BenchmarkResult.Descriptor instead.
func (BenchmarkResult) EnumDescriptor() ([]byte, []int) {
	return file_portal_v1_benchmark_proto_rawDescGZIP(), []int{0}
}

type GetBenchmarkJobRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetBenchmarkJobRequest) Reset() {
	*x = GetBenchmarkJobRequest{}
	mi := &file_portal_v1_benchmark_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBenchmarkJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBenchmarkJobRequest) ProtoMessage() {}

func (x *GetBenchmarkJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_portal_v1_benchmark_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBenchmarkJobRequest.ProtoReflect.Descriptor instead.
func (*GetBenchmarkJobRequest) Descriptor() ([]byte, []int) {
	return file_portal_v1_benchmark_proto_rawDescGZIP(), []int{0}
}

type GetBenchmarkJobResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BenchmarkId   *string                `protobuf:"bytes,1,opt,name=benchmark_id,json=benchmarkId,proto3,oneof" json:"benchmark_id,omitempty"`
	Exist         bool                   `protobuf:"varint,2,opt,name=exist,proto3" json:"exist,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetBenchmarkJobResponse) Reset() {
	*x = GetBenchmarkJobResponse{}
	mi := &file_portal_v1_benchmark_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBenchmarkJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBenchmarkJobResponse) ProtoMessage() {}

func (x *GetBenchmarkJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_portal_v1_benchmark_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBenchmarkJobResponse.ProtoReflect.Descriptor instead.
func (*GetBenchmarkJobResponse) Descriptor() ([]byte, []int) {
	return file_portal_v1_benchmark_proto_rawDescGZIP(), []int{1}
}

func (x *GetBenchmarkJobResponse) GetBenchmarkId() string {
	if x != nil && x.BenchmarkId != nil {
		return *x.BenchmarkId
	}
	return ""
}

func (x *GetBenchmarkJobResponse) GetExist() bool {
	if x != nil {
		return x.Exist
	}
	return false
}

type SendBenchmarkProgressRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BenchmarkId   string                 `protobuf:"bytes,1,opt,name=benchmark_id,json=benchmarkId,proto3" json:"benchmark_id,omitempty"`
	StartedAt     *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=started_at,json=startedAt,proto3" json:"started_at,omitempty"`
	Stdout        string                 `protobuf:"bytes,3,opt,name=stdout,proto3" json:"stdout,omitempty"` // stdout, stderrどちらも、実行中でも全てのログを送る
	Stderr        string                 `protobuf:"bytes,4,opt,name=stderr,proto3" json:"stderr,omitempty"`
	Score         int64                  `protobuf:"varint,5,opt,name=score,proto3" json:"score,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendBenchmarkProgressRequest) Reset() {
	*x = SendBenchmarkProgressRequest{}
	mi := &file_portal_v1_benchmark_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendBenchmarkProgressRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendBenchmarkProgressRequest) ProtoMessage() {}

func (x *SendBenchmarkProgressRequest) ProtoReflect() protoreflect.Message {
	mi := &file_portal_v1_benchmark_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendBenchmarkProgressRequest.ProtoReflect.Descriptor instead.
func (*SendBenchmarkProgressRequest) Descriptor() ([]byte, []int) {
	return file_portal_v1_benchmark_proto_rawDescGZIP(), []int{2}
}

func (x *SendBenchmarkProgressRequest) GetBenchmarkId() string {
	if x != nil {
		return x.BenchmarkId
	}
	return ""
}

func (x *SendBenchmarkProgressRequest) GetStartedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StartedAt
	}
	return nil
}

func (x *SendBenchmarkProgressRequest) GetStdout() string {
	if x != nil {
		return x.Stdout
	}
	return ""
}

func (x *SendBenchmarkProgressRequest) GetStderr() string {
	if x != nil {
		return x.Stderr
	}
	return ""
}

func (x *SendBenchmarkProgressRequest) GetScore() int64 {
	if x != nil {
		return x.Score
	}
	return 0
}

type SendBenchmarkProgressResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendBenchmarkProgressResponse) Reset() {
	*x = SendBenchmarkProgressResponse{}
	mi := &file_portal_v1_benchmark_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendBenchmarkProgressResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendBenchmarkProgressResponse) ProtoMessage() {}

func (x *SendBenchmarkProgressResponse) ProtoReflect() protoreflect.Message {
	mi := &file_portal_v1_benchmark_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendBenchmarkProgressResponse.ProtoReflect.Descriptor instead.
func (*SendBenchmarkProgressResponse) Descriptor() ([]byte, []int) {
	return file_portal_v1_benchmark_proto_rawDescGZIP(), []int{3}
}

type PostJobFinishedRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BenchmarkId   string                 `protobuf:"bytes,1,opt,name=benchmark_id,json=benchmarkId,proto3" json:"benchmark_id,omitempty"`
	Result        BenchmarkResult        `protobuf:"varint,2,opt,name=result,proto3,enum=portal.v1.BenchmarkResult" json:"result,omitempty"`
	FinishedAt    *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=finished_at,json=finishedAt,proto3" json:"finished_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PostJobFinishedRequest) Reset() {
	*x = PostJobFinishedRequest{}
	mi := &file_portal_v1_benchmark_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PostJobFinishedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostJobFinishedRequest) ProtoMessage() {}

func (x *PostJobFinishedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_portal_v1_benchmark_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostJobFinishedRequest.ProtoReflect.Descriptor instead.
func (*PostJobFinishedRequest) Descriptor() ([]byte, []int) {
	return file_portal_v1_benchmark_proto_rawDescGZIP(), []int{4}
}

func (x *PostJobFinishedRequest) GetBenchmarkId() string {
	if x != nil {
		return x.BenchmarkId
	}
	return ""
}

func (x *PostJobFinishedRequest) GetResult() BenchmarkResult {
	if x != nil {
		return x.Result
	}
	return BenchmarkResult_BENCHMARK_RESULT_UNSPECIFIED
}

func (x *PostJobFinishedRequest) GetFinishedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.FinishedAt
	}
	return nil
}

type PostJobFinishedResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PostJobFinishedResponse) Reset() {
	*x = PostJobFinishedResponse{}
	mi := &file_portal_v1_benchmark_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PostJobFinishedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostJobFinishedResponse) ProtoMessage() {}

func (x *PostJobFinishedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_portal_v1_benchmark_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostJobFinishedResponse.ProtoReflect.Descriptor instead.
func (*PostJobFinishedResponse) Descriptor() ([]byte, []int) {
	return file_portal_v1_benchmark_proto_rawDescGZIP(), []int{5}
}

var File_portal_v1_benchmark_proto protoreflect.FileDescriptor

var file_portal_v1_benchmark_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x65, 0x6e, 0x63,
	0x68, 0x6d, 0x61, 0x72, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x70, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x18, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x42, 0x65,
	0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x68, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72,
	0x6b, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x0c,
	0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x49,
	0x64, 0x88, 0x01, 0x01, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x78, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x05, 0x65, 0x78, 0x69, 0x73, 0x74, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x62,
	0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0x22, 0xc2, 0x01, 0x0a, 0x1c,
	0x53, 0x65, 0x6e, 0x64, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x50, 0x72, 0x6f,
	0x67, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c,
	0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x49, 0x64, 0x12,
	0x39, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x64, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x64, 0x6f,
	0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63,
	0x6f, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x22, 0x1f, 0x0a, 0x1d, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72,
	0x6b, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0xac, 0x01, 0x0a, 0x16, 0x50, 0x6f, 0x73, 0x74, 0x4a, 0x6f, 0x62, 0x46, 0x69, 0x6e,
	0x69, 0x73, 0x68, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c,
	0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x62, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x49, 0x64, 0x12,
	0x32, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1a, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x65, 0x6e, 0x63,
	0x68, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x3b, 0x0a, 0x0b, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x41, 0x74,
	0x22, 0x19, 0x0a, 0x17, 0x50, 0x6f, 0x73, 0x74, 0x4a, 0x6f, 0x62, 0x46, 0x69, 0x6e, 0x69, 0x73,
	0x68, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x89, 0x01, 0x0a, 0x0f,
	0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x20, 0x0a, 0x1c, 0x42, 0x45, 0x4e, 0x43, 0x48, 0x4d, 0x41, 0x52, 0x4b, 0x5f, 0x52, 0x45, 0x53,
	0x55, 0x4c, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x1b, 0x0a, 0x17, 0x42, 0x45, 0x4e, 0x43, 0x48, 0x4d, 0x41, 0x52, 0x4b, 0x5f, 0x52,
	0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f, 0x50, 0x41, 0x53, 0x53, 0x45, 0x44, 0x10, 0x01, 0x12, 0x1b,
	0x0a, 0x17, 0x42, 0x45, 0x4e, 0x43, 0x48, 0x4d, 0x41, 0x52, 0x4b, 0x5f, 0x52, 0x45, 0x53, 0x55,
	0x4c, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x42,
	0x45, 0x4e, 0x43, 0x48, 0x4d, 0x41, 0x52, 0x4b, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x5f,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x03, 0x32, 0xb4, 0x02, 0x0a, 0x10, 0x42, 0x65, 0x6e, 0x63,
	0x68, 0x6d, 0x61, 0x72, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x58, 0x0a, 0x0f,
	0x47, 0x65, 0x74, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x4a, 0x6f, 0x62, 0x12,
	0x21, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x42,
	0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x22, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x4a, 0x6f, 0x62, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6c, 0x0a, 0x15, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x65,
	0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x27, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64,
	0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x61,
	0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61,
	0x72, 0x6b, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x28, 0x01, 0x12, 0x58, 0x0a, 0x0f, 0x50, 0x6f, 0x73, 0x74, 0x4a, 0x6f, 0x62, 0x46,
	0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x12, 0x21, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x4a, 0x6f, 0x62, 0x46, 0x69, 0x6e, 0x69, 0x73,
	0x68, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x70, 0x6f, 0x72,
	0x74, 0x61, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x4a, 0x6f, 0x62, 0x46, 0x69,
	0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0xa9,
	0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x76, 0x31,
	0x42, 0x0e, 0x42, 0x65, 0x6e, 0x63, 0x68, 0x6d, 0x61, 0x72, 0x6b, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74,
	0x72, 0x61, 0x50, 0x74, 0x69, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x70, 0x69, 0x73, 0x63, 0x6f, 0x6e,
	0x2d, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2d, 0x76, 0x32, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2f, 0x76, 0x31, 0x3b, 0x70,
	0x6f, 0x72, 0x74, 0x61, 0x6c, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x09,
	0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x50, 0x6f, 0x72, 0x74,
	0x61, 0x6c, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x5c, 0x56,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a,
	0x50, 0x6f, 0x72, 0x74, 0x61, 0x6c, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_portal_v1_benchmark_proto_rawDescOnce sync.Once
	file_portal_v1_benchmark_proto_rawDescData = file_portal_v1_benchmark_proto_rawDesc
)

func file_portal_v1_benchmark_proto_rawDescGZIP() []byte {
	file_portal_v1_benchmark_proto_rawDescOnce.Do(func() {
		file_portal_v1_benchmark_proto_rawDescData = protoimpl.X.CompressGZIP(file_portal_v1_benchmark_proto_rawDescData)
	})
	return file_portal_v1_benchmark_proto_rawDescData
}

var file_portal_v1_benchmark_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_portal_v1_benchmark_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_portal_v1_benchmark_proto_goTypes = []any{
	(BenchmarkResult)(0),                  // 0: portal.v1.BenchmarkResult
	(*GetBenchmarkJobRequest)(nil),        // 1: portal.v1.GetBenchmarkJobRequest
	(*GetBenchmarkJobResponse)(nil),       // 2: portal.v1.GetBenchmarkJobResponse
	(*SendBenchmarkProgressRequest)(nil),  // 3: portal.v1.SendBenchmarkProgressRequest
	(*SendBenchmarkProgressResponse)(nil), // 4: portal.v1.SendBenchmarkProgressResponse
	(*PostJobFinishedRequest)(nil),        // 5: portal.v1.PostJobFinishedRequest
	(*PostJobFinishedResponse)(nil),       // 6: portal.v1.PostJobFinishedResponse
	(*timestamppb.Timestamp)(nil),         // 7: google.protobuf.Timestamp
}
var file_portal_v1_benchmark_proto_depIdxs = []int32{
	7, // 0: portal.v1.SendBenchmarkProgressRequest.started_at:type_name -> google.protobuf.Timestamp
	0, // 1: portal.v1.PostJobFinishedRequest.result:type_name -> portal.v1.BenchmarkResult
	7, // 2: portal.v1.PostJobFinishedRequest.finished_at:type_name -> google.protobuf.Timestamp
	1, // 3: portal.v1.BenchmarkService.GetBenchmarkJob:input_type -> portal.v1.GetBenchmarkJobRequest
	3, // 4: portal.v1.BenchmarkService.SendBenchmarkProgress:input_type -> portal.v1.SendBenchmarkProgressRequest
	5, // 5: portal.v1.BenchmarkService.PostJobFinished:input_type -> portal.v1.PostJobFinishedRequest
	2, // 6: portal.v1.BenchmarkService.GetBenchmarkJob:output_type -> portal.v1.GetBenchmarkJobResponse
	4, // 7: portal.v1.BenchmarkService.SendBenchmarkProgress:output_type -> portal.v1.SendBenchmarkProgressResponse
	6, // 8: portal.v1.BenchmarkService.PostJobFinished:output_type -> portal.v1.PostJobFinishedResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_portal_v1_benchmark_proto_init() }
func file_portal_v1_benchmark_proto_init() {
	if File_portal_v1_benchmark_proto != nil {
		return
	}
	file_portal_v1_benchmark_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_portal_v1_benchmark_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_portal_v1_benchmark_proto_goTypes,
		DependencyIndexes: file_portal_v1_benchmark_proto_depIdxs,
		EnumInfos:         file_portal_v1_benchmark_proto_enumTypes,
		MessageInfos:      file_portal_v1_benchmark_proto_msgTypes,
	}.Build()
	File_portal_v1_benchmark_proto = out.File
	file_portal_v1_benchmark_proto_rawDesc = nil
	file_portal_v1_benchmark_proto_goTypes = nil
	file_portal_v1_benchmark_proto_depIdxs = nil
}
