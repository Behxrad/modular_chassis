package service

import (
	"context"
	"modular_chassis/echo/pkg/services"
)

var routerIns *router

type router struct{}

func getRouter() *router {
	once.Do(func() {
		if routerIns == nil {
			routerIns = &router{}
		}
	})
	return routerIns
}

type ImplFunc[P, R any] func(ctx context.Context, request services.ServiceRequest[P]) (services.ServiceResponse[R], error)

func HandleRequest(ctx context.Context, serviceType, method string, request services.ServiceRequest[string]) (services.ServiceResponse[string], error) {
	//call := v.Method(i).Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(services.ServiceRequest[carpay.ListVehiclesRequest]{})})
	////fmt.Println(call)
	//service, exists := GetRegistry().services[serviceType+"/"+method]
	//if !exists {
	//	return services.ServiceResponse[R]{}, errors.New("service not found")
	//}
	//return service.(ImplFunc[P, R])(ctx, request)
	return services.ServiceResponse[string]{
		"{\n  \"domain\": \"authentication\",\n  \"message\": \"Request received successfully\",\n  \"service\": \"fetch_user_id\"\n}",
	}, nil
}
