// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.2
// - protoc             v3.19.4
// source: api/helloworld/v1/book.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationBookGetBook = "/api.helloworld.v1.Book/GetBook"

type BookHTTPServer interface {
	GetBook(context.Context, *GetBooksRequest) (*GetBooksReply, error)
}

func RegisterBookHTTPServer(s *http.Server, srv BookHTTPServer) {
	r := s.Route("/")
	r.GET("/book/{id}", _Book_GetBook0_HTTP_Handler(srv))
}

func _Book_GetBook0_HTTP_Handler(srv BookHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetBooksRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBookGetBook)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetBook(ctx, req.(*GetBooksRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetBooksReply)
		return ctx.Result(200, reply)
	}
}

type BookHTTPClient interface {
	GetBook(ctx context.Context, req *GetBooksRequest, opts ...http.CallOption) (rsp *GetBooksReply, err error)
}

type BookHTTPClientImpl struct {
	cc *http.Client
}

func NewBookHTTPClient(client *http.Client) BookHTTPClient {
	return &BookHTTPClientImpl{client}
}

func (c *BookHTTPClientImpl) GetBook(ctx context.Context, in *GetBooksRequest, opts ...http.CallOption) (*GetBooksReply, error) {
	var out GetBooksReply
	pattern := "/book/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationBookGetBook))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
