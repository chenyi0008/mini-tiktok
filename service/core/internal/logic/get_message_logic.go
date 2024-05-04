package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"mini-tiktok/service/message/pb/message"

	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageLogic {
	return &GetMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessageLogic) GetMessage(req *types.GetMessageRequest) (resp *types.GetMessageResponse, err error) {
	resp = new(types.GetMessageResponse)

	result, err := l.svcCtx.MessageRpc.GetMessage(l.ctx, &message.GetMessageRequest{
		ToUserId:   uint64(req.ToUserId),
		FromUserId: uint64(req.UserId),
	})

	if err != nil {
		resp.StatusCode = 1
		return
	}

	err = copier.Copy(&resp.MessageList, &result.MessageList)
	if err != nil {
		panic(err)
	}

	return
}
