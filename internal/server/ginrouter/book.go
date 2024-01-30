package ginrouter

import (
	"context"
	pb "student/api/helloworld/v1"
	ginx "student/internal/server/ginx"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type BookHttpServer interface {
	GetBook(context.Context, *pb.GetBooksRequest) (*pb.GetBooksReply, error)
}

type GinBookHttpServer struct {
	rg   *gin.RouterGroup
	book BookHttpServer
	log  *log.Helper
}

func RegisterBookHttpServer(rg *gin.RouterGroup, book BookHttpServer, logger log.Logger) {
	g := &GinBookHttpServer{
		rg:   rg,
		book: book,
		log:  log.NewHelper(logger),
	}
	rg.GET("/student", g.AsnPageListHandler(book))
}

func (g *GinBookHttpServer) AsnPageListHandler(book BookHttpServer) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var (
			ctx  = c.Request.Context()
			req  pb.GetBooksRequest
			resp *pb.GetBooksReply
			err  error
		)
		err = c.ShouldBindJSON(&req)
		if err != nil {
			g.log.WithContext(ctx).Errorf("GetMainAsnHandler: %v", err)
			AbortWithError(c, http.StatusBadRequest, err)
			return
		}
		ctx = ginx.NewDTOContext(ctx, &req)
		c.Request = c.Request.WithContext(ctx)

		resp, err = book.GetBook(ctx, &req)
		if err != nil {
			g.log.WithContext(ctx).Errorf("获取asn列表失败: %v", err)
			AbortWithError(c, http.StatusBadRequest, err)
			return
		}
		ginx.Success(c, resp)
	}

}
