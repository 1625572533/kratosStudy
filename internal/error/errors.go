package errors

import (
	"errors"
	"fmt"

	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownMessage is unknown reason for error info.
	UnknownMessage = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

//go:generate protoc -I. --go_out=paths=source_relative:. errs.proto

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d message = %s detail = %s metadata = %v", e.Code, e.Msg, e.Detail, e.Metadata)
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(httpstatus.ToGRPCCode(int(e.Code)), e.Msg).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Msg,
			Metadata: e.Metadata,
		})
	return s
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Msg == e.Msg
	}
	return false
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := proto.Clone(e).(*Error)
	err.Metadata = md
	return err
}

// New returns an error object for the code, message.
func New(code int, message, detail string) *Error {
	return &Error{
		Code:   int32(code),
		Msg:    message,
		Detail: detail,
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, message, format string, a ...interface{}) *Error {
	return New(code, message, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code int, message, format string, a ...interface{}) error {
	return New(code, message, fmt.Sprintf(format, a...))
}

// NewWithErrCode returns an error object for the code, message.
func NewWithErrCode(code, eCode int, message, detail string) *Error {
	return &Error{
		Code:    int32(code),
		ErrCode: int32(eCode),
		Msg:     message,
		Detail:  detail,
	}
}

func NewfWithErrCode(code, eCode int, message, format string, a ...interface{}) *Error {
	return NewWithErrCode(code, eCode, message, fmt.Sprintf(format, a...))
}

func ErrorfWithErrCode(code, eCode int, message, format string, a ...interface{}) error {
	return NewWithErrCode(code, eCode, message, fmt.Sprintf(format, a...))
}

// Code returns the http code for a error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 200 //nolint:gomnd
	}
	if se := FromError(err); se != nil {
		return int(se.Code)
	}
	return UnknownCode
}

// Message returns the reason for a particular error.
// It supports wrapped errors.
func Message(err error) string {
	if se := FromError(err); se != nil {
		return se.Msg
	}
	return UnknownMessage
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		ret := New(
			httpstatus.FromGRPCCode(gs.Code()),
			UnknownMessage,
			gs.Message(),
		)
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				ret.Msg = d.Reason
				return ret.WithMetadata(d.Metadata)
			}
		}
		return ret
	}
	return New(UnknownCode, UnknownMessage, err.Error())
}
