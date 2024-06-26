// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.9
// source: favorite.proto

package favorite

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Favorite_GiveLike_FullMethodName              = "/favorite.Favorite/GiveLike"
	Favorite_CancelLike_FullMethodName            = "/favorite.Favorite/CancelLike"
	Favorite_LikeList_FullMethodName              = "/favorite.Favorite/LikeList"
	Favorite_GetCommentList_FullMethodName        = "/favorite.Favorite/GetCommentList"
	Favorite_PostComment_FullMethodName           = "/favorite.Favorite/PostComment"
	Favorite_GetCommentCount_FullMethodName       = "/favorite.Favorite/GetCommentCount"
	Favorite_GetCommentCountBatch_FullMethodName  = "/favorite.Favorite/GetCommentCountBatch"
	Favorite_GetFavoriteCount_FullMethodName      = "/favorite.Favorite/GetFavoriteCount"
	Favorite_IsFavorite_FullMethodName            = "/favorite.Favorite/IsFavorite"
	Favorite_IsFavoriteBatch_FullMethodName       = "/favorite.Favorite/IsFavoriteBatch"
	Favorite_GetFavoriteCountBatch_FullMethodName = "/favorite.Favorite/GetFavoriteCountBatch"
)

// FavoriteClient is the client API for Favorite service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FavoriteClient interface {
	GiveLike(ctx context.Context, in *GiveLikeRequest, opts ...grpc.CallOption) (*Response, error)
	CancelLike(ctx context.Context, in *CancelLikeRequest, opts ...grpc.CallOption) (*Response, error)
	LikeList(ctx context.Context, in *LikeListRequest, opts ...grpc.CallOption) (*LikeListResponse, error)
	GetCommentList(ctx context.Context, in *GetCommentRequest, opts ...grpc.CallOption) (*GetCommentResponse, error)
	PostComment(ctx context.Context, in *PostCommentRequest, opts ...grpc.CallOption) (*PostCommentResponse, error)
	GetCommentCount(ctx context.Context, in *GetCommentCountRequest, opts ...grpc.CallOption) (*GetCommentCountResponse, error)
	GetCommentCountBatch(ctx context.Context, in *GetCommentCountBatchRequest, opts ...grpc.CallOption) (*GetCommentCountBatchResponse, error)
	GetFavoriteCount(ctx context.Context, in *GetFavoriteCountRequest, opts ...grpc.CallOption) (*GetFavoriteCountResponse, error)
	IsFavorite(ctx context.Context, in *IsFavoriteRequest, opts ...grpc.CallOption) (*IsFavoriteResponse, error)
	IsFavoriteBatch(ctx context.Context, in *IsFavoriteBatchRequest, opts ...grpc.CallOption) (*IsFavoriteBatchResponse, error)
	GetFavoriteCountBatch(ctx context.Context, in *GetFavoriteCountBatchRequest, opts ...grpc.CallOption) (*GetFavoriteCountBatchResponse, error)
}

type favoriteClient struct {
	cc grpc.ClientConnInterface
}

func NewFavoriteClient(cc grpc.ClientConnInterface) FavoriteClient {
	return &favoriteClient{cc}
}

func (c *favoriteClient) GiveLike(ctx context.Context, in *GiveLikeRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Favorite_GiveLike_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) CancelLike(ctx context.Context, in *CancelLikeRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Favorite_CancelLike_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) LikeList(ctx context.Context, in *LikeListRequest, opts ...grpc.CallOption) (*LikeListResponse, error) {
	out := new(LikeListResponse)
	err := c.cc.Invoke(ctx, Favorite_LikeList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) GetCommentList(ctx context.Context, in *GetCommentRequest, opts ...grpc.CallOption) (*GetCommentResponse, error) {
	out := new(GetCommentResponse)
	err := c.cc.Invoke(ctx, Favorite_GetCommentList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) PostComment(ctx context.Context, in *PostCommentRequest, opts ...grpc.CallOption) (*PostCommentResponse, error) {
	out := new(PostCommentResponse)
	err := c.cc.Invoke(ctx, Favorite_PostComment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) GetCommentCount(ctx context.Context, in *GetCommentCountRequest, opts ...grpc.CallOption) (*GetCommentCountResponse, error) {
	out := new(GetCommentCountResponse)
	err := c.cc.Invoke(ctx, Favorite_GetCommentCount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) GetCommentCountBatch(ctx context.Context, in *GetCommentCountBatchRequest, opts ...grpc.CallOption) (*GetCommentCountBatchResponse, error) {
	out := new(GetCommentCountBatchResponse)
	err := c.cc.Invoke(ctx, Favorite_GetCommentCountBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) GetFavoriteCount(ctx context.Context, in *GetFavoriteCountRequest, opts ...grpc.CallOption) (*GetFavoriteCountResponse, error) {
	out := new(GetFavoriteCountResponse)
	err := c.cc.Invoke(ctx, Favorite_GetFavoriteCount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) IsFavorite(ctx context.Context, in *IsFavoriteRequest, opts ...grpc.CallOption) (*IsFavoriteResponse, error) {
	out := new(IsFavoriteResponse)
	err := c.cc.Invoke(ctx, Favorite_IsFavorite_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) IsFavoriteBatch(ctx context.Context, in *IsFavoriteBatchRequest, opts ...grpc.CallOption) (*IsFavoriteBatchResponse, error) {
	out := new(IsFavoriteBatchResponse)
	err := c.cc.Invoke(ctx, Favorite_IsFavoriteBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) GetFavoriteCountBatch(ctx context.Context, in *GetFavoriteCountBatchRequest, opts ...grpc.CallOption) (*GetFavoriteCountBatchResponse, error) {
	out := new(GetFavoriteCountBatchResponse)
	err := c.cc.Invoke(ctx, Favorite_GetFavoriteCountBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FavoriteServer is the server API for Favorite service.
// All implementations must embed UnimplementedFavoriteServer
// for forward compatibility
type FavoriteServer interface {
	GiveLike(context.Context, *GiveLikeRequest) (*Response, error)
	CancelLike(context.Context, *CancelLikeRequest) (*Response, error)
	LikeList(context.Context, *LikeListRequest) (*LikeListResponse, error)
	GetCommentList(context.Context, *GetCommentRequest) (*GetCommentResponse, error)
	PostComment(context.Context, *PostCommentRequest) (*PostCommentResponse, error)
	GetCommentCount(context.Context, *GetCommentCountRequest) (*GetCommentCountResponse, error)
	GetCommentCountBatch(context.Context, *GetCommentCountBatchRequest) (*GetCommentCountBatchResponse, error)
	GetFavoriteCount(context.Context, *GetFavoriteCountRequest) (*GetFavoriteCountResponse, error)
	IsFavorite(context.Context, *IsFavoriteRequest) (*IsFavoriteResponse, error)
	IsFavoriteBatch(context.Context, *IsFavoriteBatchRequest) (*IsFavoriteBatchResponse, error)
	GetFavoriteCountBatch(context.Context, *GetFavoriteCountBatchRequest) (*GetFavoriteCountBatchResponse, error)
	mustEmbedUnimplementedFavoriteServer()
}

// UnimplementedFavoriteServer must be embedded to have forward compatible implementations.
type UnimplementedFavoriteServer struct {
}

func (UnimplementedFavoriteServer) GiveLike(context.Context, *GiveLikeRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GiveLike not implemented")
}
func (UnimplementedFavoriteServer) CancelLike(context.Context, *CancelLikeRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelLike not implemented")
}
func (UnimplementedFavoriteServer) LikeList(context.Context, *LikeListRequest) (*LikeListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikeList not implemented")
}
func (UnimplementedFavoriteServer) GetCommentList(context.Context, *GetCommentRequest) (*GetCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommentList not implemented")
}
func (UnimplementedFavoriteServer) PostComment(context.Context, *PostCommentRequest) (*PostCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostComment not implemented")
}
func (UnimplementedFavoriteServer) GetCommentCount(context.Context, *GetCommentCountRequest) (*GetCommentCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommentCount not implemented")
}
func (UnimplementedFavoriteServer) GetCommentCountBatch(context.Context, *GetCommentCountBatchRequest) (*GetCommentCountBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommentCountBatch not implemented")
}
func (UnimplementedFavoriteServer) GetFavoriteCount(context.Context, *GetFavoriteCountRequest) (*GetFavoriteCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFavoriteCount not implemented")
}
func (UnimplementedFavoriteServer) IsFavorite(context.Context, *IsFavoriteRequest) (*IsFavoriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsFavorite not implemented")
}
func (UnimplementedFavoriteServer) IsFavoriteBatch(context.Context, *IsFavoriteBatchRequest) (*IsFavoriteBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsFavoriteBatch not implemented")
}
func (UnimplementedFavoriteServer) GetFavoriteCountBatch(context.Context, *GetFavoriteCountBatchRequest) (*GetFavoriteCountBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFavoriteCountBatch not implemented")
}
func (UnimplementedFavoriteServer) mustEmbedUnimplementedFavoriteServer() {}

// UnsafeFavoriteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FavoriteServer will
// result in compilation errors.
type UnsafeFavoriteServer interface {
	mustEmbedUnimplementedFavoriteServer()
}

func RegisterFavoriteServer(s grpc.ServiceRegistrar, srv FavoriteServer) {
	s.RegisterService(&Favorite_ServiceDesc, srv)
}

func _Favorite_GiveLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GiveLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).GiveLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_GiveLike_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).GiveLike(ctx, req.(*GiveLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_CancelLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).CancelLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_CancelLike_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).CancelLike(ctx, req.(*CancelLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_LikeList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikeListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).LikeList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_LikeList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).LikeList(ctx, req.(*LikeListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_GetCommentList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).GetCommentList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_GetCommentList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).GetCommentList(ctx, req.(*GetCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_PostComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).PostComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_PostComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).PostComment(ctx, req.(*PostCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_GetCommentCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).GetCommentCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_GetCommentCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).GetCommentCount(ctx, req.(*GetCommentCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_GetCommentCountBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentCountBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).GetCommentCountBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_GetCommentCountBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).GetCommentCountBatch(ctx, req.(*GetCommentCountBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_GetFavoriteCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFavoriteCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).GetFavoriteCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_GetFavoriteCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).GetFavoriteCount(ctx, req.(*GetFavoriteCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_IsFavorite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsFavoriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).IsFavorite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_IsFavorite_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).IsFavorite(ctx, req.(*IsFavoriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_IsFavoriteBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsFavoriteBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).IsFavoriteBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_IsFavoriteBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).IsFavoriteBatch(ctx, req.(*IsFavoriteBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_GetFavoriteCountBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFavoriteCountBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).GetFavoriteCountBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Favorite_GetFavoriteCountBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).GetFavoriteCountBatch(ctx, req.(*GetFavoriteCountBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Favorite_ServiceDesc is the grpc.ServiceDesc for Favorite service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Favorite_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "favorite.Favorite",
	HandlerType: (*FavoriteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GiveLike",
			Handler:    _Favorite_GiveLike_Handler,
		},
		{
			MethodName: "CancelLike",
			Handler:    _Favorite_CancelLike_Handler,
		},
		{
			MethodName: "LikeList",
			Handler:    _Favorite_LikeList_Handler,
		},
		{
			MethodName: "GetCommentList",
			Handler:    _Favorite_GetCommentList_Handler,
		},
		{
			MethodName: "PostComment",
			Handler:    _Favorite_PostComment_Handler,
		},
		{
			MethodName: "GetCommentCount",
			Handler:    _Favorite_GetCommentCount_Handler,
		},
		{
			MethodName: "GetCommentCountBatch",
			Handler:    _Favorite_GetCommentCountBatch_Handler,
		},
		{
			MethodName: "GetFavoriteCount",
			Handler:    _Favorite_GetFavoriteCount_Handler,
		},
		{
			MethodName: "IsFavorite",
			Handler:    _Favorite_IsFavorite_Handler,
		},
		{
			MethodName: "IsFavoriteBatch",
			Handler:    _Favorite_IsFavoriteBatch_Handler,
		},
		{
			MethodName: "GetFavoriteCountBatch",
			Handler:    _Favorite_GetFavoriteCountBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "favorite.proto",
}
