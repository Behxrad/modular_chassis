package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
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

func (r *router) HandleRequest(ctx context.Context, serviceType, method string, request string) (string, error) {
	mInfo := GetRegistry().GetService(serviceType, method)
	if !mInfo.function.IsValid() {
		return "", errors.New("service not found")
	}

	req := reflect.New(mInfo.request).Interface()
	err := json.Unmarshal([]byte(request), req)
	if err != nil {
		return "", err
	}

	err = r.validateRequest(req)
	if err != nil {
		return "", err
	}

	reflectVals := mInfo.function.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req).Elem()})
	res, errs := reflectVals[0].Interface(), reflectVals[1].Interface()
	if errs != nil {
		return "", errs.(error)
	}

	marshal, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func (r *router) validateRequest(request any) error {
	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		fmt.Println(validationErrors)
		return err
	}
	return nil
}
