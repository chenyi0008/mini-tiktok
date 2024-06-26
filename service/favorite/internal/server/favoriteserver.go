// Code generated by goctl. DO NOT EDIT.
// Source: favorite.proto

package server

import (
	"context"

	"mini-tiktok/service/favorite/internal/logic"
	"mini-tiktok/service/favorite/internal/svc"
	"mini-tiktok/service/favorite/pb/favorite"
)

type FavoriteServer struct {
	svcCtx *svc.ServiceContext
	favorite.UnimplementedFavoriteServer
}

func NewFavoriteServer(svcCtx *svc.ServiceContext) *FavoriteServer {
	return &FavoriteServer{
		svcCtx: svcCtx,
	}
}

func (s *FavoriteServer) GiveLike(ctx context.Context, in *favorite.GiveLikeRequest) (*favorite.Response, error) {
	l := logic.NewGiveLikeLogic(ctx, s.svcCtx)
	return l.GiveLike(in)
}

func (s *FavoriteServer) CancelLike(ctx context.Context, in *favorite.CancelLikeRequest) (*favorite.Response, error) {
	l := logic.NewCancelLikeLogic(ctx, s.svcCtx)
	return l.CancelLike(in)
}

func (s *FavoriteServer) LikeList(ctx context.Context, in *favorite.LikeListRequest) (*favorite.LikeListResponse, error) {
	l := logic.NewLikeListLogic(ctx, s.svcCtx)
	return l.LikeList(in)
}

func (s *FavoriteServer) GetCommentList(ctx context.Context, in *favorite.GetCommentRequest) (*favorite.GetCommentResponse, error) {
	l := logic.NewGetCommentListLogic(ctx, s.svcCtx)
	return l.GetCommentList(in)
}

func (s *FavoriteServer) PostComment(ctx context.Context, in *favorite.PostCommentRequest) (*favorite.PostCommentResponse, error) {
	l := logic.NewPostCommentLogic(ctx, s.svcCtx)
	return l.PostComment(in)
}

func (s *FavoriteServer) GetCommentCount(ctx context.Context, in *favorite.GetCommentCountRequest) (*favorite.GetCommentCountResponse, error) {
	l := logic.NewGetCommentCountLogic(ctx, s.svcCtx)
	return l.GetCommentCount(in)
}

func (s *FavoriteServer) GetCommentCountBatch(ctx context.Context, in *favorite.GetCommentCountBatchRequest) (*favorite.GetCommentCountBatchResponse, error) {
	l := logic.NewGetCommentCountBatchLogic(ctx, s.svcCtx)
	return l.GetCommentCountBatch(in)
}

func (s *FavoriteServer) GetFavoriteCount(ctx context.Context, in *favorite.GetFavoriteCountRequest) (*favorite.GetFavoriteCountResponse, error) {
	l := logic.NewGetFavoriteCountLogic(ctx, s.svcCtx)
	return l.GetFavoriteCount(in)
}

func (s *FavoriteServer) IsFavorite(ctx context.Context, in *favorite.IsFavoriteRequest) (*favorite.IsFavoriteResponse, error) {
	l := logic.NewIsFavoriteLogic(ctx, s.svcCtx)
	return l.IsFavorite(in)
}

func (s *FavoriteServer) IsFavoriteBatch(ctx context.Context, in *favorite.IsFavoriteBatchRequest) (*favorite.IsFavoriteBatchResponse, error) {
	l := logic.NewIsFavoriteBatchLogic(ctx, s.svcCtx)
	return l.IsFavoriteBatch(in)
}

func (s *FavoriteServer) GetFavoriteCountBatch(ctx context.Context, in *favorite.GetFavoriteCountBatchRequest) (*favorite.GetFavoriteCountBatchResponse, error) {
	l := logic.NewGetFavoriteCountBatchLogic(ctx, s.svcCtx)
	return l.GetFavoriteCountBatch(in)
}
