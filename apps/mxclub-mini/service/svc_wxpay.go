package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/domain/user/repo"
)

func init() {
	jet.Provide(NewWxPayService)
}

type WxPayService struct {
	userRepo repo.IUserRepo
}

func NewWxPayService(repo repo.IUserRepo) *WxPayService {
	return &WxPayService{userRepo: repo}
}
