package handler

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"mini-tiktok/core/internal/logic"
	"mini-tiktok/core/internal/svc"
	"mini-tiktok/core/internal/types"
)

func PostMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostMessageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		get := r.Header.Get("UserId")
		userId, _ := strconv.Atoi(get)
		req.UserId = uint(userId)

		l := logic.NewPostMessageLogic(r.Context(), svcCtx)
		resp, err := l.PostMessage(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}