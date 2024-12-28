package service

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"modular_chassis/echo/pkg/errs"
	"reflect"
)

var routerIns *router
var validate = validator.New(validator.WithRequiredStructEnabled())

type router struct{}

func GetRouter() *router {
	once.Do(func() {
		if routerIns == nil {
			routerIns = &router{}
		}
	})
	return routerIns
}

func (r *router) HandleRequest(ctx context.Context, serviceType, method string, request any) (any, error) {
	mInfo := GetRegistry().GetService(serviceType, method)
	if !mInfo.Function.IsValid() {
		return "", errors.New("service not found")
	}

	err := r.validateRequest(request)
	if err != nil {
		return "", errs.NewServiceErrorCodeAndCause(errs.BadRequest, err)
	}

	reflectVals := mInfo.Function.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(request).Elem()})
	res, errs := reflectVals[0].Interface(), reflectVals[1].Interface()
	if errs != nil {
		return "", errs.(error)
	}

	return res, nil
}

func (r *router) validateRequest(request any) error {
	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Println(validationErrors)
		return err
	}
	return nil
}
