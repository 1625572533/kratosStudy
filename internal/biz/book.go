package biz

import (
	"context"
	"student/internal/model"

	"github.com/go-kratos/kratos/v2/log"
)

//定义student的操作接口
type BookRepo interface{
	GetBook(context.Context,int32)(*model.TBook,error)
}

type BookUsecase struct{
	repo BookRepo
	log *log.Helper
}

func NewBookUsecase(repo BookRepo,logger log.Logger) *BookUsecase{
	return &BookUsecase{repo:repo,log:log.NewHelper(logger)}
}

//通过id获取student信息
func (uc *BookUsecase) GetOne(ctx context.Context,id int32) (*model.TBook,error){
	uc.log.WithContext(ctx).Infof("biz.get: %d",id)
	return uc.repo.GetBook(ctx,id)
}