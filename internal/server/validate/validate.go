package validate

import (
	"context"
	"strings"
	"sync"

	"student/internal/error"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	once             sync.Once
	Validator        *validator.Validate
	uni              *ut.UniversalTranslator
	enTrans, zhTrans ut.Translator
)

func setup() {
	once.Do(func() {
		var found bool

		Validator = validator.New()

		ent := en.New()
		uni = ut.New(ent, ent)
		enTrans, found = uni.GetTranslator("en")
		if !found {
			panic("en translation not found")
		}
		zht := zh.New()
		uni = ut.New(zht, zht)
		zhTrans, found = uni.GetTranslator("zh")
		if !found {
			panic("zh translation not found")
		}
		err := en_translations.RegisterDefaultTranslations(Validator, enTrans)
		if err != nil {
			panic(err)
		}

		err = zh_translations.RegisterDefaultTranslations(Validator, zhTrans)
		if err != nil {
			panic(err)
		}
	})
}

func init() {
	setup()
}

func GetValidator() *validator.Validate {
	if Validator == nil {
		setup()
	}
	return Validator
}

type protoMessage interface {
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)
}

type localeTranslationKey struct{}

func NewTranslationContext(ctx context.Context, locale string) context.Context {
	var translation = enTrans
	switch locale {
	case "zh", "zh-CN":
		translation = zhTrans
	}
	return context.WithValue(ctx, localeTranslationKey{}, translation)
}

func FromTranslationContext(ctx context.Context) (ut.Translator, bool) {
	trans, ok := ctx.Value(localeTranslationKey{}).(ut.Translator)
	return trans, ok
}

func Validating(ctx context.Context, field interface{}) error {
	err := Validator.Struct(field)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			if trans, ok := FromTranslationContext(ctx); ok && len(errs) > 0 {
				errArr := make([]string, 0, len(errs))
				var es string
				for _, fe := range errs {
					errStr := fe.Translate(trans)
					errArr = append(errArr, errStr)
				}
				es = strings.Join(errArr, ", ")
				return errors.New(400, es, es)
			}
		}
		return err
	}
	return nil
}

type Option func(*options)

type options struct {
	validate *validator.Validate
}

func WithValidator(v *validator.Validate) Option {
	return func(o *options) {
		o.validate = v
	}
}

// Validate is a validator middleware.
func Validate(opts ...Option) middleware.Middleware {
	op := options{
		validate: Validator,
	}
	for _, o := range opts {
		o(&op)
	}
	if op.validate == nil {
		panic("validate is nil")
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if _, ok := req.(protoMessage); ok {
				err = op.validate.StructCtx(ctx, req)
				if err != nil {
					if errs, ok := err.(validator.ValidationErrors); ok {
						if trans, ok := FromTranslationContext(ctx); ok && len(errs) > 0 {
							errArr := make([]string, 0, len(errs))
							var es string
							for _, fe := range errs {
								errStr := fe.Translate(trans)
								errArr = append(errArr, errStr)
							}
							es = strings.Join(errArr, ", ")
							return nil, errors.New(400, es, es)
						}
					}
					return nil, err
				}
			}
			return handler(ctx, req)
		}
	}
}
