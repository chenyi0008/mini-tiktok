package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"mini-tiktok/common/consts"
	"mini-tiktok/service/favorite/pb/favorite"

	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLikeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLikeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLikeListLogic {
	return &GetLikeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLikeListLogic) GetLikeList(req *types.GetLikeListRequest) (resp *types.GetLikeListResponse, err error) {
	resp = new(types.GetLikeListResponse)
	result, err := l.svcCtx.FavoriteRpc.LikeList(l.ctx, &favorite.LikeListRequest{
		UserId: uint64(req.UserId),
	})
	if err != nil {
		return
	}

	arr := make([]int, 0)
	for _, id := range result.VideoId {
		arr = append(arr, int(id))
	}

	resp.StatusMsg = result.Message
	if result.Code == consts.SUCCEED {
		return
	} else {
		//batch, err := l.svcCtx.RedisCli.GetVideoInfoBatch(l.ctx, arr)
		list, err2 := l.svcCtx.VideoModel.ListInIds(result.VideoId)
		if err2 != nil {
			return nil, err2
		}
		copier.Copy(&resp.VideoList, &list)
	}

	return
}
