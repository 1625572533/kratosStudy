package server

import (
	// v1 "student/api/helloworld/v1"
	"student/internal/conf"
	"student/internal/service"

	"student/internal/server/ginrouter"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, student *service.StudentService,
	book *service.BookService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	// v1.RegisterGreeterHTTPServer(srv, greeter)
	// v1.RegisterStudentHTTPServer(srv, student)

	router := gin.New()
	apis := router.Group("/v1/apis")
	ginrouter.RegisterBookHttpServer(apis, book, logger)
	srv.HandlePrefix("/", router)

	return srv
}
