package service

import (
	"context"

	pb "student/api/helloworld/v1"
	"student/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"encoding/json"
)

type BookService struct {
	pb.UnimplementedBookServer
	book *biz.BookUsecase
	log *log.Helper
}

func NewBookService(book *biz.BookUsecase,logger log.Logger) *BookService {
	return &BookService{
		book : book,
		log: log.NewHelper(logger),
	}
}

func (uc *BookService) GetBook(ctx context.Context, req *pb.GetBooksRequest) (*pb.GetBooksReply, error) {
	val,_:=json.Marshal(req)
	uc.log.Infof("GetBook req: ",string(val))
	book,err := uc.book.GetOne(ctx,req.Id)
	if err!=nil{
		return nil,err
	}
	return &pb.GetBooksReply{
		Id:      book.Id,
		Name:    book.Name,
		Author:  book.Author,
		Price:   book.Price,
		Stock:   book.Stock,
		Sales:   book.Sales,
		ImaPath: book.ImgPath,
	},nil
}
