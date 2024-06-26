// Code generated by goctl. DO NOT EDIT.
// Source: follow.proto

package server

import (
	"context"

	"mini-tiktok/service/follow/internal/logic"
	"mini-tiktok/service/follow/internal/svc"
	"mini-tiktok/service/follow/pb/follow"
)

type FollowServer struct {
	svcCtx *svc.ServiceContext
	follow.UnimplementedFollowServer
}

func NewFollowServer(svcCtx *svc.ServiceContext) *FollowServer {
	return &FollowServer{
		svcCtx: svcCtx,
	}
}

func (s *FollowServer) PostFollow(ctx context.Context, in *follow.PostFollowRequest) (*follow.Response, error) {
	l := logic.NewPostFollowLogic(ctx, s.svcCtx)
	return l.PostFollow(in)
}

func (s *FollowServer) GetFans(ctx context.Context, in *follow.GetFansRequest) (*follow.GetFansResponse, error) {
	l := logic.NewGetFansLogic(ctx, s.svcCtx)
	return l.GetFans(in)
}

func (s *FollowServer) GetFollowingList(ctx context.Context, in *follow.GetFollowingListRequest) (*follow.GetFollowingListResponse, error) {
	l := logic.NewGetFollowingListLogic(ctx, s.svcCtx)
	return l.GetFollowingList(in)
}

func (s *FollowServer) GetFriendList(ctx context.Context, in *follow.GetFriendListRequest) (*follow.GetFriendListResponse, error) {
	l := logic.NewGetFriendListLogic(ctx, s.svcCtx)
	return l.GetFriendList(in)
}
