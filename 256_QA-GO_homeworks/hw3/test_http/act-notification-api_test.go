//go:build e2e_test_noallure

package test_http

import (
	"encoding/json"
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

// Basic notification lifecycle - send, receive, acknowledge
func TestNotificationLifecycle(t *testing.T) {
	client := NewHTTPClient(5, 1*time.Second)
	currentHour := time.Now().UTC().Hour()
	t.Run("Test notification lifecycle", func(t *testing.T) {
		platform := "ios"
		userId := "1"
		message := "Order is ready"
		lang := "LANG_ENGLISH"
		// Create a device
		resp, err := client.CreateDevice(platform, userId)
		if err != nil {
			t.Fatal(err)
		}

		deviceId := ""
		// Post returns 200, write to issues
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Status code is not %v. Got %v instead", http.StatusOK, resp.StatusCode)
		} else {
			var createDevResp CreateDeviceResp
			err = json.NewDecoder(resp.Body).Decode(&createDevResp)
			if err != nil {
				t.Fatal(err)
			}
			deviceId = createDevResp.DeviceId
			if deviceId == "" {
				t.Fatal("Device ID is empty")
			}
		}

		// Send notification
		resp, err = client.SendNotification(deviceId, message)
		if err != nil {
			t.Fatal(err)
		}
		notifId := ""
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Status code is not %v. Got %v instead", http.StatusOK, resp.StatusCode)
		} else {
			var SendNotifResp SendNotificationResp
			err = json.NewDecoder(resp.Body).Decode(&SendNotifResp)
			if err != nil {
				t.Fatal(err)
			}
			notifId = SendNotifResp.NotificationId
			if notifId == "" {
				t.Fatal("Notification ID is empty")
			}
		}

		// Get notification
		resp, err = client.GetNotification(deviceId)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Status code is not %v. Got %v instead", http.StatusOK, resp.StatusCode)
		} else {
			var Notif GetNotificationResp
			err = json.NewDecoder(resp.Body).Decode(&Notif)
			if err != nil {
				t.Fatal(err)
			}
			if Notif.Notification[0].NotificationId != notifId {
				t.Fatal("Notification ID mismatch")
			} else if Notif.Notification[0].Message != GetGreeting(currentHour, lang)+message {
				t.Errorf("Message is not correct. Got: %v. Want: %v", Notif.Notification[0].Message,
					GetGreeting(currentHour, lang)+message)
			}
		}

		// Ack notification
		resp, err = client.AckNotification(notifId)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Status code is not %v. Got %v instead", http.StatusOK, resp.StatusCode)
		}
		AckNotifResp := new(AckNotificationResp)
		err = json.NewDecoder(resp.Body).Decode(&AckNotifResp)
		if !AckNotifResp.Success {
			t.Fatal("Success is not true")
		}
	})
}
