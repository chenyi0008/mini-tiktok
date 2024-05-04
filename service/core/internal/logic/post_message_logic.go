package logic

import (
	"context"
	"mini-tiktok/service/core/internal/svc"
	"mini-tiktok/service/core/internal/types"
	"mini-tiktok/service/message/pb/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostMessageLogic {
	return &PostMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostMessageLogic) PostMessage(req *types.PostMessageRequest) (resp *types.PostMessageResponse, err error) {
	resp = new(types.PostMessageResponse)
	result, err := l.svcCtx.MessageRpc.PostMessage(l.ctx, &message.PostMessageRequest{
		ToUserId:   uint64(req.ToUserId),
		FormUserId: uint64(req.UserId),
		Content:    req.Content,
		ActionType: uint32(req.ActionType),
	})

	if err != nil {
		resp.StatusMsg = result.Message
		resp.StatusCode = 1
	}

	return
}
