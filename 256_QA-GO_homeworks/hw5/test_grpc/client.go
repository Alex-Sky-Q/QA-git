// Package test_grpc tests gRPC handles
package test_grpc

import (
	"context"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"google.golang.org/grpc"
	"log"
	"time"
)

const baseURL = "localhost:8082"

// NewGrpcClient creates a new gRPC device client
func NewGrpcClient() act_device_api.ActDeviceApiServiceClient {
	conn, err := grpc.Dial(baseURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := act_device_api.NewActDeviceApiServiceClient(conn)
	return client
}

// NewGrpcNotifClient creates a new gRPC notification client
func NewGrpcNotifClient() act_device_api.ActNotificationApiServiceClient {
	conn, err := grpc.Dial(baseURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := act_device_api.NewActNotificationApiServiceClient(conn)
	return client
}

// SetTimeoutCtx sets a timeout for a context
func SetTimeoutCtx(timeout int) context.Context {
	timeoutCtx, _ := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(timeout),
	)
	//defer cancel()
	return timeoutCtx
}
