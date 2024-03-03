package config

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

// use a single instance , it caches struct info
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func InitTranslator(language string) error {
	// 修改gin框架中的validator引擎属性
	validate = binding.Validator.Engine().(*validator.Validate)

	// 注册一个获取json tag的自定义方法
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	zhtrans := zh.New()
	entrans := en.New()
	// 第一个参数是fallback locale 后面的参数是locales it should support
	uni = ut.New(entrans, zhtrans, entrans)
	var ok bool
	trans, ok = uni.GetTranslator(language)
	if !ok {
		return fmt.Errorf("uni.GetTranslator(%s)", language)
	}
	switch language {
	case "entrans":
		err := entranslations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			return err
		}
	case "zhtrans":
		err := zhtranslations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			return err
		}
	default:
		err := zhtranslations.RegisterDefaultTranslations(validate, trans)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisterDTO.name 删除错误信息前面的RegisterDTO.
func removeTopStruct(errMsg map[string]string) map[string]string {
	res := make(map[string]string)
	for field, msg := range errMsg {
		res[field[strings.Index(field, ".")+1:]] = msg
	}
	return res
}

func ValidateError(c *gin.Context, err error) {
	// 打印出 body
	data, _ := io.ReadAll(c.Request.Body)
	fmt.Printf("req.body=%s\n, content-type=%v\n", data, c.ContentType())
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	// 如果不是参数错误，比如是json格式错误
	if !ok {
		zap.L().Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//"error": {
	//	"RegisterDTO.Name": "Name长度必须至少为3个字符",
	//	"RegisterDTO.RePassword": "RePassword必须等于Password"
	//}
	zap.L().Error("errors", zap.Any("errors", removeTopStruct(errs.Translate(trans))))
	c.JSON(http.StatusBadRequest, gin.H{
		"data": removeTopStruct(errs.Translate(trans)),
	})
}
