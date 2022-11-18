//go:build e2e_test

package test_http

import (
	"encoding/json"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"net/http"
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

type NotificationHttpTestSuite struct {
	suite.Suite
}

func (s *NotificationHttpTestSuite) TestNotificationLifecycleHTTP(t provider.T) {
	t.Title("HTTP - Test notification lifecycle (CRUD)")
	t.Description("Basic notification lifecycle - send, receive, acknowledge")
	client := NewHTTPClient(5, 1*time.Second)

	// Create a device
	deviceId := ""
	t.WithTestSetup(func(t provider.T) {
		platform := "ios"
		userId := "1"
		t.WithNewStep("Create a device", func(sCtx provider.StepCtx) {
			resp, err := client.CreateDevice(platform, userId)
			sCtx.Require().NoError(err)
			// Post returns 200, why not 201?
			sCtx.Assert().Equal(http.StatusOK, resp.StatusCode, "Status code should be correct. Want: %v. Got %v",
				http.StatusOK, resp.StatusCode)

			var createDevResp CreateDeviceResp
			err = json.NewDecoder(resp.Body).Decode(&createDevResp)
			sCtx.Require().NoError(err)
			deviceId = createDevResp.DeviceId
			sCtx.Assert().NotEmpty(deviceId, "Device ID should not be empty")
			sCtx.WithNewParameters("deviceId", deviceId)
		})
	})
	defer t.WithTestTeardown(func(t provider.T) {
		t.WithNewStep("Delete device", func(sCtx provider.StepCtx) {
			//sCtx.WithNewParameters("deviceId", deviceId)
			_, err := client.RemoveDevice(deviceId)
			sCtx.Require().NoError(err)
		})
	})

	// Send notification
	notifId := ""
	currentHour := time.Now().UTC().Hour()
	message := "Order is ready"
	lang := "LANG_ENGLISH"
	t.WithNewStep("Send notification", func(sCtx provider.StepCtx) {
		resp, err := client.SendNotification(deviceId, message)
		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(http.StatusOK, resp.StatusCode, "Status code should be correct. Want: %v. Got %v",
			http.StatusOK, resp.StatusCode)

		var SendNotifResp SendNotificationResp
		err = json.NewDecoder(resp.Body).Decode(&SendNotifResp)
		sCtx.Require().NoError(err)
		notifId = SendNotifResp.NotificationId
		sCtx.Assert().NotEmpty(notifId, "Notification ID should not be empty")
	})

	// Get notification
	t.WithNewStep("Get notification", func(sCtx provider.StepCtx) {
		resp, err := client.GetNotification(deviceId)
		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(http.StatusOK, resp.StatusCode, "Status code should be correct. Want: %v. Got %v",
			http.StatusOK, resp.StatusCode)

		var Notif GetNotificationResp
		err = json.NewDecoder(resp.Body).Decode(&Notif)
		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(notifId, Notif.Notification[0].NotificationId, "Notification IDs should match")
		sCtx.Assert().Equal(GetGreeting(currentHour, lang)+message, Notif.Notification[0].Message,
			"Messages should match. Got: %v. Want: %v", Notif.Notification[0].Message, GetGreeting(currentHour, lang)+message)
	})

	// Ack notification
	t.WithNewStep("Ack notification", func(sCtx provider.StepCtx) {
		resp, err := client.AckNotification(notifId)
		sCtx.Require().NoError(err)
		sCtx.Assert().Equal(http.StatusOK, resp.StatusCode, "Status code should be correct. Want: %v. Got %v",
			http.StatusOK, resp.StatusCode)

		AckNotifResp := new(AckNotificationResp)
		err = json.NewDecoder(resp.Body).Decode(&AckNotifResp)
		sCtx.Require().NoError(err)
		sCtx.Assert().True(AckNotifResp.Success, "Success should be true")
	})
}

func TestRunner(t *testing.T) {
	suite.RunSuite(t, new(NotificationHttpTestSuite))
}
