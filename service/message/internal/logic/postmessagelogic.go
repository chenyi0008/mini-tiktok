package logic

import (
	"context"

	"mini-tiktok/service/message/internal/svc"
	"mini-tiktok/service/message/pb/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostMessageLogic {
	return &PostMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostMessageLogic) PostMessage(in *message.PostMessageRequest) (*message.Response, error) {
	resp := new(message.Response)
	_, err := l.svcCtx.MessageModel.Create(uint(in.ToUserId), uint(in.FormUserId), in.Content)
	if err != nil {
		resp.Message = err.Error()
		resp.Code = 1
	}
	return resp, nil
}
