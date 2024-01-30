package ginx

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"net/http"
	"strings"

	errors "student/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/valyala/bytebufferpool"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	baseContentType = "application"
)

func Error(c *gin.Context, err interface{}) {
	if err == nil {
		c.Status(http.StatusOK)
		return
	}

	var (
		ctx     = c.Request.Context()
		jsonErr error
		data    []byte
		code    int
	)

	if e, ok := err.(error); ok {
		se := errors.FromError(e)
		code = int(se.Code)
		if len(se.Msg) == 0 {
			se.Msg = se.String()
		}

		data, jsonErr = json.Marshal(se)
		if jsonErr != nil {
			panic(jsonErr)
		}
		//c.AbortWithStatusJSON(code, se)
	} else {
		data, jsonErr = json.Marshal(err)
		if jsonErr != nil {
			panic(jsonErr)
		}
		code = http.StatusInternalServerError
	}

	if trace.SpanContextFromContext(ctx).HasSpanID() {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Key("http.response.body").String(string(data)))
	}
	c.Abort()
	c.Data(code, "application/json; charset=utf-8", data)

	return
}

func Success(c *gin.Context, obj interface{}) {
	var (
		ctx = c.Request.Context()
	)

	if obj == nil {
		obj = struct{}{}
	}

	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	if trace.SpanContextFromContext(ctx).HasSpanID() {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Key("http.response.body").String(string(data)))
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
	return
}

func ShouldBind(c *gin.Context, req interface{}) error {
	rv := reflect.ValueOf(req)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(req)}
	}

	var (
		method = c.Request.Method
		ctx    = c.Request.Context()
	)

	if trace.SpanContextFromContext(ctx).HasSpanID() {
		buf := bytebufferpool.Get()
		if c.Request.Header.Clone().Write(buf) == nil {
			trace.SpanFromContext(ctx).SetAttributes([]attribute.KeyValue{
				attribute.Key("http.header").String(fmt.Sprintf("%+v", buf.String())),
			}...)
		}
		bytebufferpool.Put(buf)
	}

	switch method {
	case http.MethodGet:
		err := c.ShouldBindQuery(req)
		if err != nil {
			return err
		}

		if trace.SpanContextFromContext(ctx).HasSpanID() {
			trace.SpanFromContext(ctx).SetAttributes([]attribute.KeyValue{
				attribute.Key("query").String(fmt.Sprintf("%+v", req)),
			}...)
		}

	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		data, err := c.GetRawData()
		if err != nil {
			return err
		}

		if trace.SpanContextFromContext(ctx).HasSpanID() {
			trace.SpanFromContext(ctx).SetAttributes(attribute.Key("body").String(string(data)))
			trace.SpanFromContext(ctx).SetAttributes(semconv.HTTPServerAttributesFromHTTPRequest("", c.FullPath(), c.Request)...)
		}

		err = json.Unmarshal(data, req)
		if err != nil {
			return err
		}
		//TODO(如果出现既有Body又有Query, 怎么处理? 两种情况, 一种是Query也是req的字段, 另一种是Query不属于req字段)

	default:
		return fmt.Errorf("undefine request method: %v", method)
	}

	return nil
}

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}

// Middlewares return middlewares wrapper
func Middlewares(m ...middleware.Middleware) gin.HandlerFunc {
	chain := middleware.Chain(m...)
	return func(c *gin.Context) {
		next := func(ctx context.Context, req interface{}) (interface{}, error) {
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			var err error
			if c.Writer.Status() >= 400 {
				err = errors.Errorf(c.Writer.Status(), errors.UnknownMessage, errors.UnknownMessage)
			}
			return c.Writer, err
		}
		next = chain(next)
		ctx := NewGinContext(c.Request.Context(), c)
		if ginCtx, ok := FromGinContext(ctx); ok {
			khttp.SetOperation(ctx, ginCtx.FullPath())
		}
		_, err := next(ctx, c)
		if err != nil {
			if c.Writer.Size() > 0 {
				c.Writer.Flush()
				return
			}
			se := errors.FromError(err)
			c.AbortWithStatusJSON(int(se.Code), se)
			return
		}
	}
}

type dtoKey struct{}

// NewDTOContext .
func NewDTOContext(ctx context.Context, dto interface{}) context.Context {
	ctx = context.WithValue(ctx, dtoKey{}, dto)
	return ctx
}

// FromDTOContext .
func FromDTOContext(ctx context.Context) (dto interface{}, ok bool) {
	dto, ok = ctx.Value(dtoKey{}).(interface{})
	return
}

type ginKey struct{}

// NewGinContext returns a new Context that carries ginx.Context value.
func NewGinContext(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ginKey{}, c)
}

// FromGinContext returns the ginx.Context value stored in ctx, if any.
func FromGinContext(ctx context.Context) (c *gin.Context, ok bool) {
	c, ok = ctx.Value(ginKey{}).(*gin.Context)
	return
}
