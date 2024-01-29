package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Student struct {
	ID        int32
	Name      string
	Info      string
	Status    int32
	UpdatedAt time.Time
	CreatedAt time.Time
}

// 定义Student得操作接口
type StudentRepo interface {
	GetStudent(context.Context, int32) (*Student, error)
}

type StudentUsecase struct {
	repo StudentRepo
	log  *log.Helper
}

func NewStudentUsecase(repo StudentRepo, logger log.Logger) *StudentUsecase {
	return &StudentUsecase{repo: repo, log: log.NewHelper(logger)}
}

// 通过id获取student信息
func (uc *StudentUsecase) Get(ctx context.Context, id int32) (*Student, error) {
	uc.log.WithContext(ctx).Infof("biz.get: %d", id)
	return uc.repo.GetStudent(ctx, id)
}
