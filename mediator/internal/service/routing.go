package service

import (
	"context"
	"encoding/json"
	"fmt"
	"modular_chassis/echo/pkg/services"
	"reflect"
)

var routerIns *router

type router struct{}

func GetRouter() *router {
	once.Do(func() {
		if routerIns == nil {
			routerIns = &router{}
		}
	})
	return routerIns
}

func (r *router) HandleRequest(ctx context.Context, serviceType, method string, request services.ServiceRequest[string]) (services.ServiceResponse[string], error) {
	mInfo := GetRegistry().GetService(serviceType, method)
	req := reflect.New(mInfo.request).Elem().Interface()
	err := json.Unmarshal([]byte(request.Payload), &req)
	if err != nil {
		return services.ServiceResponse[string]{}, err
	}

	//call := mInfo.function.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(services.ServiceRequest{Payload: req})})
	////fmt.Println(call)
	//service, exists := GetRegistry().services[serviceType+"/"+method]
	//if !exists {
	//	return services.ServiceResponse[R]{}, errors.New("service not found")
	//}
	//return service.(ImplFunc[P, R])(ctx, request)

	requestType := reflect.TypeOf(services.ServiceRequest[any]{})
	fmt.Println(requestType)
	//serviceRequestType := reflect.StructOf([]reflect.StructField{
	//	{
	//		Name: "Payload",
	//		Type: mInfo.request,
	//		Tag:  reflect.StructTag(`json:"payload"`),
	//	},
	//})

	// Create a new instance of ServiceRequest with the dynamic payload type
	serviceRequestValue := reflect.New(requestType).Elem()

	// Set the Payload field with the provided value
	serviceRequestValue.FieldByName("Payload").Set(reflect.ValueOf(req))

	call := mInfo.function.Call([]reflect.Value{reflect.ValueOf(context.Background()), serviceRequestValue})
	fmt.Println(call)

	return services.ServiceResponse[string]{
		"{\n  \"domain\": \"authentication\",\n  \"message\": \"Request received successfully\",\n  \"service\": \"fetch_user_id\"\n}",
	}, nil
}
