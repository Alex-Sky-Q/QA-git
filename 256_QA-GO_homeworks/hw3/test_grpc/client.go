package test_grpc

import (
	"context"
	//"context"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"google.golang.org/grpc"
	"log"
	"time"
	//"time"
)

const baseURL = "localhost:8082"

func NewGrpcClient() act_device_api.ActDeviceApiServiceClient {
	conn, err := grpc.Dial(baseURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := act_device_api.NewActDeviceApiServiceClient(conn)

	return client
}

func NewGrpcNotifClient() act_device_api.ActNotificationApiServiceClient {
	conn, err := grpc.Dial(baseURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := act_device_api.NewActNotificationApiServiceClient(conn)

	return client
}

func SetTimeoutCtx(timeout int) context.Context {
	timeoutCtx, _ := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(timeout),
	)
	//defer cancel()
	return timeoutCtx
}
