//go:build e2e_test

package test_grpc

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"testing"
	"time"
)

func GetGreeting(hour int, lang string) string {
	greeting := ""
	if lang == "LANG_ENGLISH" {
		greeting = "Good night "
		if hour > 24 || hour < 0 {
			return "Hour is out of range, should be between 0 and 24"
		}
		if hour >= 6 && hour < 11 {
			greeting = "Good morning "
		} else if hour >= 11 && hour < 15 {
			greeting = "Good afternoon "
		} else if hour >= 15 && hour < 21 {
			greeting = "Good evening "
		}
	} else if lang == "LANG_ESPANOL" {
		greeting = "Buenas noches "
		if hour > 24 || hour < 0 {
			return "Hour is out of range, should be between 0 and 24"
		}
		if hour >= 6 && hour < 11 {
			greeting = "Buenos dias "
		} else if hour >= 11 && hour < 15 {
			greeting = "Buenas tardes "
		} else if hour >= 15 && hour < 21 {
			greeting = "Buenas noches "
		}
	} else if lang == "LANG_ITALIAN" {
		greeting = "Buona notte "
		if hour > 24 || hour < 0 {
			return "Hour is out of range, should be between 0 and 24"
		}
		if hour >= 6 && hour < 11 {
			greeting = "Buon giorno "
		} else if hour >= 11 && hour < 15 {
			greeting = "Buon pomeriggio "
		} else if hour >= 15 && hour < 21 {
			greeting = "Buona serata "
		}
	} else if lang == "LANG_RUSSIAN" {
		greeting = "Доброй ночи "
		if hour > 24 || hour < 0 {
			return "Hour is out of range, should be between 0 and 24"
		}
		if hour >= 6 && hour < 11 {
			greeting = "Доброе утро "
		} else if hour >= 11 && hour < 15 {
			greeting = "Добрый день "
		} else if hour >= 15 && hour < 21 {
			greeting = "Добрый вечер "
		}
	}
	return greeting
}

func CreateDevice(platform string, userId uint64) (uint64, error) {
	clientDev := NewGrpcClient()
	ctx := SetTimeoutCtx(5)
	// Create a device
	reqCreate := act_device_api.CreateDeviceV1Request{
		Platform: platform,
		UserId:   userId,
	}
	respCreate, err := clientDev.CreateDeviceV1(ctx, &reqCreate)
	if err != nil {
		return 0, err
	}
	deviceId := respCreate.DeviceId

	// Check that device was created
	reqDesc := act_device_api.DescribeDeviceV1Request{DeviceId: deviceId}
	_, err = clientDev.DescribeDeviceV1(ctx, &reqDesc)

	return deviceId, err
}

type NotifStreamGrpcTestSuite struct {
	suite.Suite
}

// Test stream handle
func (n *NotifStreamGrpcTestSuite) TestStreamNotificationAfter(t provider.T) {
	t.Title("gRPC - Test stream - send notification AFTER subscribe")
	t.Description("Subscribe device then send notification")

	// Test setup
	var deviceId uint64
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create device", func(sCtx provider.StepCtx) {
			var err error
			deviceId, err = CreateDevice("ios", 1)
			sCtx.Require().NoError(err)
			sCtx.Require().NotEmpty(deviceId, "Device ID should not be empty")
			sCtx.WithNewParameters("deviceId", deviceId)
		})
	})

	clientN := NewGrpcNotifClient()
	//ctx := context.Background()
	ctx := SetTimeoutCtx(5)

	// Subscribe
	var respSub act_device_api.ActNotificationApiService_SubscribeNotificationClient
	t.WithNewStep("Subscribe", func(sCtx provider.StepCtx) {
		reqSub := act_device_api.SubscribeNotificationRequest{DeviceId: deviceId}
		var err error
		respSub, err = clientN.SubscribeNotification(ctx, &reqSub)
		sCtx.Require().NoError(err)
	})

	// Send notification
	var notifId uint64
	message := "Order is ready"
	expMessage := GetGreeting(time.Now().UTC().Hour(), "LANG_RUSSIAN") + message
	t.WithNewStep("Send notification", func(sCtx provider.StepCtx) {
		reqSend := act_device_api.SendNotificationV1Request{Notification: &act_device_api.Notification{
			//NotificationId:     0,
			DeviceId:           deviceId,
			Username:           "1",
			Message:            message,
			Lang:               act_device_api.Language_LANG_RUSSIAN,
			NotificationStatus: act_device_api.Status_STATUS_CREATED,
		}}
		respSend, err := clientN.SendNotificationV1(ctx, &reqSend)
		notifId = respSend.NotificationId
		sCtx.Require().NoError(err)
		sCtx.Assert().NotEmpty(notifId, "Notification ID should not be empty")
	})

	// Receive from stream
	var msg act_device_api.UserNotification
	t.WithNewStep("Receive from stream", func(sCtx provider.StepCtx) {
		err := respSub.RecvMsg(&msg)
		//msg, err := respSub.Recv()
		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(notifId, msg.NotificationId, "Notification IDs should be equal. Got: %v. Want: %v",
			msg.NotificationId, notifId)
		sCtx.Assert().Equal(expMessage, msg.Message, "Messages should match. Got: %s. Want: %s", msg.Message,
			expMessage)
	})

	//Ack notification
	t.WithNewStep("Ack notification", func(sCtx provider.StepCtx) {
		reqAck := act_device_api.AckNotificationV1Request{NotificationId: notifId}
		respAck, err := clientN.AckNotification(ctx, &reqAck)
		sCtx.Require().NoError(err)
		sCtx.Assert().True(respAck.Success, "Success should be True")
	})

	// Check that received notification is not delivered again
	t.WithNewStep("Check that received notification is not delivered again", func(sCtx provider.StepCtx) {
		err := respSub.RecvMsg(&msg)
		sCtx.Require().Error(err)
		sCtx.Assert().Equal("rpc error: code = DeadlineExceeded desc = context deadline exceeded", err.Error(),
			"Error code should be as expected")
	})
}

// Bug. Notifications sent before subscription are not delivered
func (n *NotifStreamGrpcTestSuite) TestStreamNotificationBefore(t provider.T) {
	t.Skip()
	t.Title("gRPC - Test stream - send notification BEFORE subscribe")
	t.Description("Send notification then subscribe device")

	// Test setup
	var deviceId uint64
	t.WithTestSetup(func(t provider.T) {
		t.WithNewStep("Create device", func(sCtx provider.StepCtx) {
			var err error
			deviceId, err = CreateDevice("ios", 1)
			sCtx.Require().NoError(err)
			sCtx.Require().NotEmpty(deviceId, "Device ID should not be empty")
			sCtx.WithNewParameters("deviceId", deviceId)
		})
	})

	clientN := NewGrpcNotifClient()
	//ctx := context.Background()
	ctx := SetTimeoutCtx(5)

	// Send notification
	var notifId uint64
	message := "Order is ready"
	expMessage := GetGreeting(time.Now().UTC().Hour(), "LANG_RUSSIAN") + message
	t.WithNewStep("Send notification", func(sCtx provider.StepCtx) {
		reqSend := act_device_api.SendNotificationV1Request{Notification: &act_device_api.Notification{
			//NotificationId:     0,
			DeviceId:           deviceId,
			Username:           "1",
			Message:            message,
			Lang:               act_device_api.Language_LANG_RUSSIAN,
			NotificationStatus: act_device_api.Status_STATUS_CREATED,
		}}
		respSend, err := clientN.SendNotificationV1(ctx, &reqSend)
		notifId = respSend.NotificationId
		sCtx.Require().NoError(err)
		sCtx.Assert().NotEmpty(notifId, "Notification ID should not be empty")
	})

	// Subscribe
	var respSub act_device_api.ActNotificationApiService_SubscribeNotificationClient
	t.WithNewStep("Subscribe", func(sCtx provider.StepCtx) {
		reqSub := act_device_api.SubscribeNotificationRequest{DeviceId: deviceId}
		var err error
		respSub, err = clientN.SubscribeNotification(ctx, &reqSub)
		sCtx.Require().NoError(err)
	})

	// Receive from stream
	var msg act_device_api.UserNotification
	t.WithNewStep("Receive from stream", func(sCtx provider.StepCtx) {
		err := respSub.RecvMsg(&msg)
		//msg, err := respSub.Recv()
		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(notifId, msg.NotificationId, "Notification IDs should be equal. Got: %v. Want: %v",
			msg.NotificationId, notifId)
		sCtx.Assert().Equal(expMessage, msg.Message, "Messages should match. Got: %s. Want: %s", msg.Message,
			expMessage)
	})
}

func TestRunner(t *testing.T) {
	suite.RunSuite(t, new(NotifStreamGrpcTestSuite))
}
