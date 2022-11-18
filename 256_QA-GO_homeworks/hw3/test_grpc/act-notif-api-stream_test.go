//go:build e2e_test_noallure

package test_grpc

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// Test stream handle
func TestStreamNotification(t *testing.T) {
	// Test setup
	deviceId, err := CreateDevice("ios", 1)
	require.NoError(t, err)
	require.NotEmptyf(t, deviceId, "Device ID is empty")

	clientN := NewGrpcNotifClient()
	//ctx := context.Background()
	ctx := SetTimeoutCtx(5)

	t.Run("Test stream - send notification AFTER subscribe", func(t *testing.T) {
		// Subscribe
		reqSub := act_device_api.SubscribeNotificationRequest{DeviceId: deviceId}
		respSub, err := clientN.SubscribeNotification(ctx, &reqSub)
		require.NoError(t, err)

		// Send notification
		message := "Order is ready"
		expMessage := GetGreeting(time.Now().UTC().Hour(), "LANG_RUSSIAN") + message
		reqSend := act_device_api.SendNotificationV1Request{Notification: &act_device_api.Notification{
			//NotificationId:     0,
			DeviceId:           deviceId,
			Username:           "1",
			Message:            message,
			Lang:               act_device_api.Language_LANG_RUSSIAN,
			NotificationStatus: act_device_api.Status_STATUS_CREATED,
		}}
		respSend, err := clientN.SendNotificationV1(ctx, &reqSend)
		notifId := respSend.NotificationId
		require.NoError(t, err)
		require.NotEmptyf(t, notifId, "Notification ID is empty")

		// Receive from stream
		var msg act_device_api.UserNotification
		err = respSub.RecvMsg(&msg)
		//msg, err := respSub.Recv()
		require.NoError(t, err)
		assert.Equalf(t, msg.NotificationId, notifId, "Notification IDs not equal. Got: %v. Want: %v",
			msg.NotificationId, notifId)
		assert.Equalf(t, msg.Message, expMessage, "Message is not correct. Got: %s. Want: %s", msg.Message,
			expMessage)

		// Ack notification
		reqAck := act_device_api.AckNotificationV1Request{NotificationId: notifId}
		respAck, err := clientN.AckNotification(ctx, &reqAck)
		require.NoError(t, err)
		assert.Truef(t, respAck.Success, "Success is not True")

		// Check that received notification is not delivered again
		err = respSub.RecvMsg(&msg)
		require.Error(t, err)
		assert.Equalf(t, err.Error(), "rpc error: code = DeadlineExceeded desc = context deadline exceeded",
			"Error code is not as expected")
	})

	// Bug. Write to issues. Notifications sent before subscription are not delivered
	t.Run("Test stream - send notification BEFORE subscribe", func(t *testing.T) {
		// Send notification
		message := "Order is ready"
		expMessage := GetGreeting(time.Now().Hour()-3, "LANG_RUSSIAN") + message
		reqSend := act_device_api.SendNotificationV1Request{Notification: &act_device_api.Notification{
			//NotificationId:     0,
			DeviceId:           deviceId,
			Username:           "1",
			Message:            message,
			Lang:               act_device_api.Language_LANG_RUSSIAN,
			NotificationStatus: act_device_api.Status_STATUS_CREATED,
		}}
		respSend, err := clientN.SendNotificationV1(ctx, &reqSend)
		notifId := respSend.NotificationId
		require.NoError(t, err)
		require.NotEmptyf(t, notifId, "Notification ID is empty")

		// Subscribe
		reqSub := act_device_api.SubscribeNotificationRequest{DeviceId: deviceId}
		respSub, err := clientN.SubscribeNotification(ctx, &reqSub)
		require.NoError(t, err)

		// Receive from stream
		var msg act_device_api.UserNotification
		err = respSub.RecvMsg(&msg)
		//msg, err := respSub.Recv()
		require.NoError(t, err)
		assert.Equalf(t, msg.NotificationId, notifId, "Notification IDs not equal. Got: %v. Want: %v",
			msg.NotificationId, notifId)
		assert.Equalf(t, msg.Message, expMessage, "Message is not correct. Got: %s. Want: %s", msg.Message,
			expMessage)
	})

}
