package logic

import (
	"context"
	"fmt"
	"mini-tiktok/service/favorite/internal/svc"
	"mini-tiktok/service/favorite/pb/favorite"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeListLogic {
	return &LikeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeListLogic) LikeList(in *favorite.LikeListRequest) (*favorite.LikeListResponse, error) {
	list, err := l.svcCtx.FavoriteModel.GetByUserId(in.UserId)

	resp := new(favorite.LikeListResponse)
	if err != nil {
		resp.Code = 1
		resp.Message = "查询失败"
		return resp, nil
	}
	var resList []uint64
	for _, item := range *list {
		fmt.Printf("%+v\n", item)
		resList = append(resList, uint64(item.VideoId))
	}

	return &favorite.LikeListResponse{
		Code:    0,
		VideoId: resList,
		Message: "查询成功",
	}, nil
	return resp, nil
}
