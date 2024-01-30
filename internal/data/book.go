package data

import (
	"context"

	"student/internal/biz"
	"student/internal/model"

	"github.com/go-kratos/kratos/v2/log"
)

type BookRepo struct {
	data *Data
	log  *log.Helper
}

func NewBookRepo(data *Data, logger log.Logger) biz.BookRepo {
	return &BookRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (uc BookRepo) GetBook(ctx context.Context, id int32) (*model.TBook, error) {
	var book model.TBook
	uc.data.gormDB.Where("id = ?", id).First(&book)
	uc.log.WithContext(ctx).Info("gormDB: GetStudent,id: ",id)
	return &book,nil;
}
